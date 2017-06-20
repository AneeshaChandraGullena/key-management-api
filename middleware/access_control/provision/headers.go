// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package provision

import (
	"encoding/json"
	"errors"
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
		// Correlation-ID is a required header for looking up information within the Logs,
		// the header must also be a valid UUID v4.
		// This condition will throw an error if Correlation-ID is not included as a header and is not a valid UUID.
		if userCorrelationID, ok := helpers.ValidateUUIDHeader(constants.CorrelationIDHeader, request); ok {
			writer.Header().Set(constants.CorrelationIDHeader, userCorrelationID)
		} else {
			badRequest := http.StatusText(http.StatusBadRequest) + ": Request requires valid Correlation-ID"
			middlewareLogger.Log("err", badRequest)
			errorCollection := collections.NewErrorCollection().Append(errors.New(badRequest))
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		// Authorization is a required header for looking up information within the Database,
		// This condition will throw an error if Authorization is not included as a header.
		// It does not check if it is valid Authorization. That is left to another middleware
		if auth := request.Header.Get(constants.AuthorizationHeader); auth != "" {
			writer.Header().Set(constants.AuthorizationHeader, auth)
		} else {
			badRequest := http.StatusText(http.StatusBadRequest) + ": Request requires Authorization"
			middlewareLogger.Log("err", badRequest, "correlation-id", request.Header.Get(constants.CorrelationIDHeader))
			errorCollection := collections.NewErrorCollection().Append(errors.New(badRequest))
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}
		h.ServeHTTP(writer, request)
	})
}
