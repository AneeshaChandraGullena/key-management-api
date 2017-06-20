// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package standard

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"

	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/access_control/helpers"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/utils/logging"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-models/collections"
)

var middlewareLogger log.Logger

func init() {
	middlewareLogger = logging.GlobalLogger()
}

// AccessControlHeaders check the headers that come in on the request
func AccessControlHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(constants.ContentTypeHeader, constants.AppJSONMime+"; charset=utf-8")
		writer.Header().Set(constants.AllowOriginHeader, "*")
		writer.Header().Set(constants.StrictTransportSecurity, "max-age=31536000; includeSubDomains; preload")
		methods := fmt.Sprintf("%s, %s, %s, %s, %s", http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodHead)
		writer.Header().Set(constants.AllowMethodsHeader, methods)
		constantsAllowed := fmt.Sprintf("Origin, %s, %s, %s, %s, %s", constants.ContentTypeHeader, constants.BluemixSpaceHeader, constants.BluemixOrgHeader, constants.AuthorizationHeader, constants.Prefer)
		writer.Header().Set(constants.AllowHeadersHeader, constantsAllowed)
		exposedHeaders := fmt.Sprintf("%s, %s", constants.CorrelationIDHeader, constants.KeyTotalForSpaceHeader)
		writer.Header().Set(constants.ExposeHeadersHeader, exposedHeaders)

		if request.Method == http.MethodOptions {
			return
		}
		// Authorization is a required header for looking up information within the Database,
		// This condition will throw an error is Authorization is not including as a header.
		// It does not check if it is valid Authorization. That is left to another middleware
		if auth := request.Header.Get(constants.AuthorizationHeader); auth != "" {
			writer.Header().Set(constants.AuthorizationHeader, auth)
		} else {
			unauthorized := http.StatusText(http.StatusUnauthorized) + ": Request requires Authorization"
			middlewareLogger.Log("err", unauthorized, "correlation_id", request.Header.Get(constants.CorrelationIDHeader))
			errorCollection := collections.NewErrorCollection().Append(errors.New(unauthorized))
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		// Bluemix-Org is a required field for looking up information within the Database,
		// This condition will throw an error if Bluemix-org is not included as a header containing a valid UUID.
		// It does not check if it is a valid Org.
		if org, ok := helpers.ValidateUUIDHeader(constants.BluemixOrgHeader, request); ok {
			writer.Header().Set(constants.BluemixOrgHeader, org)
		} else {
			badRequest := http.StatusText(http.StatusBadRequest) + ": Request requires Bluemix-Org Header containing a valid UUID"
			middlewareLogger.Log("err", badRequest, "correlation_id", request.Header.Get(constants.CorrelationIDHeader))
			errorCollection := collections.NewErrorCollection().Append(errors.New(badRequest))
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		// Bluemix-Space is a required field for looking up information within the Database,
		// This condition will throw an error if Bluemix-Space is not included as a header containing a valid UUID.
		// It does not check if it is a valid space. That is left to another middleware
		if space, ok := helpers.ValidateUUIDHeader(constants.BluemixSpaceHeader, request); ok {
			writer.Header().Set(constants.BluemixSpaceHeader, space)
		} else {
			badRequest := http.StatusText(http.StatusBadRequest) + ": Request requires Bluemix-Space Header containing a valid UUID"
			middlewareLogger.Log("err", badRequest, "correlation_id", request.Header.Get(constants.CorrelationIDHeader))
			errorCollection := collections.NewErrorCollection().Append(errors.New(badRequest))
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		h.ServeHTTP(writer, request)
	})
}
