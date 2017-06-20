// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package correlation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/satori/go.uuid"

	"github.ibm.com/Alchemy-Key-Protect/key-management-api/utils/logging"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-models/collections"
)

var middlewareLogger log.Logger

func init() {
	middlewareLogger = logging.GlobalLogger()
}

// ID will transport or create a new Correlation id for a request depending on if one exists
func ID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if userCorrelationID := request.Header.Get(constants.CorrelationIDHeader); userCorrelationID != "" {
			if _, err := uuid.FromString(userCorrelationID); err != nil {
				badRequest := http.StatusText(http.StatusBadRequest) + ": Request requires valid UUID v4 Correlation-ID"
				middlewareLogger.Log("err", badRequest)
				errorCollection := collections.NewErrorCollection().Append(errors.New(badRequest))
				writer.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(writer).Encode(errorCollection)
				return
			}
			writer.Header().Set(constants.CorrelationIDHeader, userCorrelationID)
		} else {
			uuidV4 := uuid.NewV4()
			request.Header.Set(constants.CorrelationIDHeader, uuidV4.String())
			writer.Header().Set(constants.CorrelationIDHeader, uuidV4.String())
		}

		handler.ServeHTTP(writer, request)
	})
}
