// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.

package roles

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	uuid "github.com/satori/go.uuid"

	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/cache"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/helpers"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/utils/logging"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-models/collections"
)

const (
	// Auditor role can access the HEAD API to determine the number of key that are in a space.
	Auditor = "auditors"

	// Developer role can do all the things that an Auditor can do, as well as create, retrieve and update a secret
	Developer = "developers"

	// Manager role can do all the things that a Developer can do, as well as delete a secret.
	Manager = "managers"

	// login in url for dallas.
	loginURL = "https://login%s.ng.bluemix.net/UAALoginServerWAR/rolecheck?role=%s&space_guid=%s"

	// user ID key, found in the decoded token.
	user_id_key = "user_id"
)

var (
	// Array of roles in order from greatest to least privilege
	roles = []string{Manager, Developer, Auditor}

	// rolesMap maps a role string to a header role.
	rolesMap = map[string]string{
		Auditor:   constants.RoleAuditor,
		Developer: constants.RoleDeveloper,
		Manager:   constants.RoleManager,
	}

	stage1           string
	middlewareLogger log.Logger
)

func init() {
	middlewareLogger = logging.GlobalLogger()
	host, err := os.Hostname()
	if err != nil {
		panic("Unable to find hostname")
	}

	// Ansible declares hostanems by <env>-<datacenter>-keyprotect-<machinerole>-<instance>-<domain>
	env := strings.Split(host, "-")[0]
	if env != "prod" {
		// the '.' in the assignment is so that this will fit into the formatted string loginURL.
		// for example, if this is a prod env, then the var stage1 will be "" and loginURL will have a host of login.ng.bluemix.net.
		// if it is a staging env such as dev, prestage or stage then the var stage1 will be ".stage1" and loginURL will be login.stage1.ng.bluemix.net
		stage1 = ".stage1"
	}
}

// RoleCheck is a handler for checking the users role
func RoleCheck(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get(constants.AuthorizationHeader)
		userCorrelationID := request.Header.Get(constants.CorrelationIDHeader)
		spaceGUID := request.Header.Get(constants.BluemixSpaceHeader)

		// Check that bearer is in token
		baseSplit := strings.Split(authToken, " ")
		if len(baseSplit) != 2 || strings.ToLower(baseSplit[0]) != "bearer" {
			err := errors.New("token should contain 'bearer '")
			middlewareLogger.Log("err", err, "correlation_id", userCorrelationID)
			errorCollection := collections.NewErrorCollection().Append(err)
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		userID, err := acquireID(authToken)
		if err != nil {
			errReponse := errors.New("Unable to get user ID from token")
			middlewareLogger.Log("err", errReponse, "correlation_id", userCorrelationID)
			errorCollection := collections.NewErrorCollection().Append(err)
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		var role string
		cacheKey := makeCacheKey(authToken, spaceGUID)
		hit, role := cache.Get(cacheKey)
		if !hit {
			role, err = performCheck(authToken, userCorrelationID, spaceGUID, writer, request)
			if err != nil {
				return
			}
		}

		// protects us against someone passing in a good role with a bad space.
		// it would allow a user to pass in a different user's space, a valid token and a valid role (i.e Developer or Manager) on the header and be able access get routes
		if role == "" {
			errReponse := errors.New(http.StatusText(http.StatusForbidden) + ": User does not have access to provided space")
			middlewareLogger.Log("err", errReponse, "correlation_id", userCorrelationID)
			errorCollection := collections.NewErrorCollection().Append(errReponse)
			writer.WriteHeader(http.StatusForbidden)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		} else {
			request.Header.Set(constants.BluemixUserRole, rolesMap[role])
			request.Header.Set(constants.UserIDHeader, userID)
		}

		if !hit {
			cache.Insert(cacheKey, 0, true, role)
		}

		handler.ServeHTTP(writer, request)
	})
}

func performCheck(authToken, userCorrelationID, spaceGUID string, writer http.ResponseWriter, request *http.Request) (string, error) {
	//Loop through all possible roles.
	//The roles array MUST be in highest -> lowest, or else higher roles will return auditor.
	for _, roleToCheck := range roles {
		// TODO: for now, only supporting dallas. GHE issue https://github.ibm.com/Alchemy-Key-Protect/key-protect-backlog/issues/820
		url := fmt.Sprintf(loginURL, stage1, roleToCheck, spaceGUID)
		requestRoleCheck, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			middlewareLogger.Log("err", err, "correlation_id", userCorrelationID)
			errRequest := errors.New(http.StatusText(http.StatusServiceUnavailable) + ": Unable to make request at this time")
			middlewareLogger.Log("err", errRequest, "correlation_id", userCorrelationID)
			errorCollection := collections.NewErrorCollection().Append(errRequest)
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errorCollection)
			return "", errRequest
		}

		// setup headers
		requestRoleCheck.Header.Set(constants.AuthorizationHeader, authToken)

		client := &http.Client{}
		response, err := client.Do(requestRoleCheck)
		if err != nil {
			middlewareLogger.Log("err", err, "correlation_id", userCorrelationID)
			clientError := errors.New(http.StatusText(http.StatusBadGateway) + ": Unable to communicate with the verification service")
			middlewareLogger.Log("err", clientError, "correlation_id", userCorrelationID)
			errorCollection := collections.NewErrorCollection().Append(clientError)
			writer.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(writer).Encode(errorCollection)
			return "", clientError
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			responseError := errors.New(http.StatusText(response.StatusCode) + ": Please check authorization token and space guid")
			middlewareLogger.Log("err", responseError, "correlation_id", userCorrelationID)
			errorCollection := collections.NewErrorCollection().Append(responseError)
			writer.WriteHeader(response.StatusCode)
			json.NewEncoder(writer).Encode(errorCollection)
			return "", responseError
		}

		decodedResponse, err := marshallResponse(userCorrelationID, response, writer)
		if err != nil {
			return "", err
		}
		if hasAccessRaw, ok := decodedResponse["hasaccess"]; ok {
			if hasAccess, ok := hasAccessRaw.(bool); ok {
				if hasAccess {
					return roleToCheck, nil
				}
			} else {
				// if this happens something has gone wrong with the login server as hasaccess should be of type bool.
				errReponse := errors.New(http.StatusText(http.StatusInternalServerError) + ": Error occured while attemping to determine user's access")
				middlewareLogger.Log("err", errReponse, "correlation_id", userCorrelationID)
				errorCollection := collections.NewErrorCollection().Append(errReponse)
				writer.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(writer).Encode(errorCollection)
				return "", errReponse
			}
		}
	}
	return "", nil
}

func marshallResponse(userCorrelationID string, response *http.Response, writer http.ResponseWriter) (map[string]interface{}, error) {
	// convert body to bytes for json Unmarshalling
	body, _ := ioutil.ReadAll(response.Body)
	// interface to decode the response body
	var decodeResponseInterface interface{}
	errDecodeResponseInterface := json.Unmarshal(body, &decodeResponseInterface)
	if errDecodeResponseInterface != nil {
		middlewareLogger.Log("err", errDecodeResponseInterface, "correlation_id", userCorrelationID)
		errReponse := errors.New(http.StatusText(http.StatusInternalServerError) + ": Error occured while attemping to decode user's role")
		middlewareLogger.Log("err", errReponse, "correlation_id", userCorrelationID)
		errorCollection := collections.NewErrorCollection().Append(errReponse)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorCollection)
		return nil, errReponse
	}
	return decodeResponseInterface.(map[string]interface{}), nil
}

func acquireID(authToken string) (string, error) {
	baseSplit := strings.Split(authToken, " ")
	var err error
	if len(baseSplit) != 2 || strings.ToLower(baseSplit[0]) != "bearer" {
		err = errors.New(http.StatusText(http.StatusBadRequest) + " : Token should be prefixed with bearer")
		return "", err
	}

	authToken = baseSplit[1]

	splitAuth := strings.Split(authToken, ".")

	// JWT tokens have 3 segments
	if len(splitAuth) != 3 {
		err = errors.New(http.StatusText(http.StatusBadRequest) + " : missing or malformed token")
		return "", err
	}

	encodedData := splitAuth[1]
	var decodedData []byte
	if decodedData, err = helpers.DecodeSegment(encodedData); err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err = json.Unmarshal(decodedData, &data); err != nil {
		return "", err
	}

	userID := data[user_id_key].(string)
	if _, err = uuid.FromString(userID); err != nil {
		return "", err
	}

	return userID, nil
}

func makeCacheKey(authToken, space string) string {
	return authToken + space
}
