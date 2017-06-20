// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package provision

import (
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

func testHandler(writer http.ResponseWriter, request *http.Request) {
	return
}

//test for bad request http response
func TestAccessControlHeadersNoCorrelationID(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	checkHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	checkHandler.ServeHTTP(responseRecorder, testRequest)

	expectedResponse := ""
	if actualResponse := responseRecorder.Header().Get(constants.CorrelationIDHeader); expectedResponse != actualResponse {
		t.Errorf("Expected %s, received %s", expectedResponse, actualResponse)
	}
}

//test for valid http response
func TestAccessControlHeadersValidHeaders(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	checkHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	UUIDv4 := uuid.NewV4().String()
	testRequest.Header.Set(constants.AuthorizationHeader, "test")
	testRequest.Header.Set(constants.CorrelationIDHeader, UUIDv4)

	checkHandler.ServeHTTP(responseRecorder, testRequest)

	if actualAuth := responseRecorder.Header().Get(constants.AuthorizationHeader); actualAuth != "test" {
		t.Errorf("Expected %s, received %s", "test", actualAuth)
	}

	if actualCorr := responseRecorder.Header().Get(constants.CorrelationIDHeader); actualCorr != UUIDv4 {
		t.Errorf("Expected %s, received %s", UUIDv4, actualCorr)
	}
}

func TestAccessControlHeadersBadAuth(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	checkHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	UUIDv4 := uuid.NewV4().String()
	testRequest.Header.Set(constants.CorrelationIDHeader, UUIDv4)
	testRequest.Header.Set(constants.AuthorizationHeader, "")

	checkHandler.ServeHTTP(responseRecorder, testRequest)

	if actualAuth := responseRecorder.Header().Get(constants.AuthorizationHeader); actualAuth != "" {
		t.Errorf("Expected %s, received %s", "", actualAuth)
	}
}
