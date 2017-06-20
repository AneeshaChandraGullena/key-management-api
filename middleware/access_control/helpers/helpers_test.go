// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package helpers

import (
	"net/http"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

func TestValidateUUIDHeaderInvalidRequest(t *testing.T) {

	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	testRequest.Header.Set(constants.AuthorizationHeader, "")

	headerResponse, isValidResponse := ValidateUUIDHeader(constants.AuthorizationHeader, testRequest)

	if headerResponse != "" {
		t.Errorf("Expected %s, received %s", "", headerResponse)
	}

	if isValidResponse != false {
		t.Errorf("Expected %t, received %t", true, isValidResponse)
	}
}

func TestValidateUUIDHeaderInvalidRequest2(t *testing.T) {

	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	testRequest.Header.Set(constants.AuthorizationHeader, "test")

	headerResponse, isValidResponse := ValidateUUIDHeader(constants.AuthorizationHeader, testRequest)

	if headerResponse != "" {
		t.Errorf("Expected %s, received %s", "", strings.TrimSpace(headerResponse))
	}

	if isValidResponse != false {
		t.Errorf("Expected %t, received %t", false, isValidResponse)
	}
}

func TestValidateUUIDHeaderValidRequest(t *testing.T) {

	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	validAuthHeader := uuid.NewV4().String()

	testRequest.Header.Set(constants.AuthorizationHeader, validAuthHeader) //set valid UUID

	headerResponse, isValidResponse := ValidateUUIDHeader(constants.AuthorizationHeader, testRequest)

	if headerResponse != validAuthHeader {
		t.Errorf("Expected %s, received %s", validAuthHeader, headerResponse)
	}

	if isValidResponse != true {
		t.Errorf("Expected %t, received %t", true, isValidResponse)
	}
}
