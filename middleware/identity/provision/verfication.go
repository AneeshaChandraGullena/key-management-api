// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package provision

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/helpers"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/utils/logging"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-config"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-models/collections"
)

var (
	basicAuth        string
	regions          = []string{constants.DallasRegionSubDomain, constants.LondonRegionSubDomain}
	once             sync.Once
	middlewareLogger log.Logger
	config           configuration.Configuration
)

func init() {
	middlewareLogger = logging.GlobalLogger()
	config = configuration.Get()
}

func getAuthURLFromToken(authToken string) (issuerURL string, err error) {
	baseSplit := strings.Split(authToken, " ")
	if len(baseSplit) != 2 || strings.ToLower(baseSplit[0]) != "bearer" {
		err = fmt.Errorf("token should contain 'bearer '")
		middlewareLogger.Log("err", err)
		return
	}

	authToken = baseSplit[1]

	var environment string
	var region string

	encodedSplitAuth := strings.Split(authToken, ".")
	encodedUserInfo := encodedSplitAuth[1]

	decodedUserInfo, err := helpers.DecodeSegment(encodedUserInfo)
	if err != nil {
		middlewareLogger.Log("err", err)
		return
	}

	var data map[string]interface{}
	if err = json.Unmarshal(decodedUserInfo, &data); err != nil {
		middlewareLogger.Log("err", err)
		return
	}

	// iss -> issuer
	issuerBaseURL := data["iss"].(string)

	for _, r := range regions {
		if strings.Contains(issuerBaseURL, r) {
			region = r
			break
		}
	}

	if strings.Contains(issuerBaseURL, constants.StagingSubDomain) {
		environment = "." + constants.StagingSubDomain
	}
	middlewareLogger.Log("issuerBaseURL", issuerBaseURL)

	issuerURL = fmt.Sprintf("https://login%s.%s.bluemix.net/UAALoginServerWAR/check_token?token=%s", environment, region, authToken)
	return
}

func validClientToken(issuerURL string) (response *http.Response, err error) {
	request, err := http.NewRequest("POST", issuerURL, nil)
	if err != nil {
		middlewareLogger.Log("err", err)
		return
	}

	// setup headers
	request.Header.Set(constants.AuthorizationHeader, basicAuth)

	client := &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		middlewareLogger.Log("err", err)
		return
	}
	defer response.Body.Close()
	return
}

// VerifyIdentity is used to verify the identity of a client token
func VerifyIdentity(handler http.Handler) http.Handler {
	// Initalize the basicAuth URL the first time it is called
	once.Do(func() {
		stringToBase64Encode := config.GetString("authorization_info.client_token_credentials.client_id") + ":" + config.GetString("authorization_info.client_token_credentials.client_secret")
		basicAuth = "Basic " + base64.URLEncoding.EncodeToString([]byte(stringToBase64Encode))
	})

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get(constants.AuthorizationHeader)
		issuerURL, issuerErr := getAuthURLFromToken(authToken)
		if issuerErr != nil {
			unauthorized := http.StatusText(http.StatusUnauthorized) + ": Request requires valid Authorization"
			middlewareLogger.Log("err", unauthorized)
			errorCollection := collections.NewErrorCollection().Append(errors.New(unauthorized))
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		hystrix.ConfigureCommand(constants.HystrixBluemixUserCommand, hystrix.CommandConfig{
			Timeout:               config.GetInt(constants.HystrixBluemixLoginTimeoutInMs),
			MaxConcurrentRequests: config.GetInt(constants.HystrixBluemixLoginMaxConcurrentRequests),
			ErrorPercentThreshold: config.GetInt(constants.HystrixBluemixLoginMaxConcurrentRequests),
		})

		//response, err := validClientToken(authToken)
		var response *http.Response
		var err error
		hystrix.Do(constants.HystrixBluemixUserCommand, func() error {
			response, err = validClientToken(issuerURL)
			if err != nil {
				return err
			}
			return nil
		}, func(err error) error {
			response, err = validClientToken(helpers.FindAlternativeIssuer(issuerURL))
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			middlewareLogger.Log("err", err, "correlation-id", request.Header.Get(constants.CorrelationIDHeader))
			errorCollection := collections.NewErrorCollection().Append(err)
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		middlewareLogger.Log("msg", response.StatusCode)
		if err != nil || response.StatusCode != http.StatusOK {
			unauthorized := http.StatusText(http.StatusUnauthorized) + ": Request requires valid Authorization"
			middlewareLogger.Log("err", unauthorized, "correlation-id", request.Header.Get(constants.CorrelationIDHeader))
			errorCollection := collections.NewErrorCollection().Append(errors.New(unauthorized))
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		handler.ServeHTTP(writer, request)
	})
}
