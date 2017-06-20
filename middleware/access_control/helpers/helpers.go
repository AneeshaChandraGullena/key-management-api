// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package helpers

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// ValidateUUIDHeader ensures that a a header that is suppose to be a uuid is a uuid
func ValidateUUIDHeader(header string, request *http.Request) (string, bool) {
	var isValid bool
	var headerValue string

	if gotHeader := request.Header.Get(header); gotHeader != "" {
		if _, err := uuid.FromString(gotHeader); err == nil {
			isValid = true
			headerValue = gotHeader
		}
	}
	return headerValue, isValid
}
