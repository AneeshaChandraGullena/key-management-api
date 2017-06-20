// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package ratelimiting

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/youtube/vitess/go/ratelimiter"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/utils/logging"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-config"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-models/collections"
)

var (
	limiter *ratelimiter.RateLimiter

	middlewareLogger log.Logger
	config           configuration.Configuration
)

func init() {
	middlewareLogger = logging.GlobalLogger()
	config = configuration.Get()

	limiter = ratelimiter.NewRateLimiter(config.GetInt(constants.RateLimitingAmount), time.Duration(config.GetInt(constants.RateLimitingSeconds))*time.Second)
}

// Limit will set a limit for the rate of a request.
func Limit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if limiter.Allow() == false {
			middlewareLogger.Log("msg", "Rate Limit exceeded")
			tooMany := http.StatusText(http.StatusTooManyRequests) + ": Please try again."
			errorCollection := collections.NewErrorCollection().Append(errors.New(tooMany))
			writer.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}
		handler.ServeHTTP(writer, request)
	})
}
