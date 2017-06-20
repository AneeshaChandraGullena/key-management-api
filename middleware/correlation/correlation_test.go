// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.
package correlation

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

func TestIDValidCorrHeader(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := ID(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	UUIDv4 := uuid.NewV4().String()

	testRequest.Header.Set(constants.CorrelationIDHeader, UUIDv4)
	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	if corrHeader := responseRecorder.Header().Get(constants.CorrelationIDHeader); corrHeader != UUIDv4 {
		t.Errorf("Expected %s, received %s", UUIDv4, corrHeader)
	}
}

func TestIDBadCorrHeadert(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := ID(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testRequest.Header.Set(constants.CorrelationIDHeader, "test")
	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	//400 bad response, CorrelationIDHeader is not a valid UUID
	if httpCode := responseRecorder.Code; httpCode != 400 {
		t.Errorf("Expected %d, received %d", 200, httpCode)
	}
}

func TestIDBlankCorrHeader(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := ID(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testRequest.Header.Set(constants.CorrelationIDHeader, "")
	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	//corrHeader will get a new valid UUID
	if corrHeader := responseRecorder.Header().Get(constants.CorrelationIDHeader); corrHeader == "" {
		t.Errorf("Expected %s, received %s", "", corrHeader)
	}
}
