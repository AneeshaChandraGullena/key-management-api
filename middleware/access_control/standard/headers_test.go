// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package standard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

// testHandler only returns so that we can get the response from the middleware
func testHandler(writer http.ResponseWriter, request *http.Request) {
	return
}

func TestReturnedHandlerAuthEmpty(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testRequest.Header.Set(constants.AuthorizationHeader, "")

	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	//error 401 == unauthorized
	if responseRecorder.Code != 401 {
		t.Errorf("Expected %d, received %d", 400, responseRecorder.Code)
	}
}

func TestAccessControlHeadersNewAuth(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testAuth := "test"
	testRequest.Header.Set(constants.AuthorizationHeader, testAuth)

	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	auth := responseRecorder.Header().Get(constants.AuthorizationHeader)

	if auth != testAuth {
		t.Errorf("Expected %s, received %s", testAuth, auth)
	}
}

func TestAccessControlHeadersGoodBluemix(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	checkHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testAuth := "test"
	orgUUID := uuid.NewV4().String()
	spaceUUID2 := uuid.NewV4().String()

	testRequest.Header.Set(constants.AuthorizationHeader, testAuth)
	testRequest.Header.Set(constants.BluemixOrgHeader, orgUUID)
	testRequest.Header.Set(constants.BluemixSpaceHeader, spaceUUID2)

	checkHandler.ServeHTTP(responseRecorder, testRequest)

	returnedUUIDv4 := responseRecorder.Header().Get(constants.BluemixOrgHeader)
	if returnedUUIDv4 != orgUUID {
		t.Errorf("Expected %s, received %s", orgUUID, returnedUUIDv4)
	}
}

func TestAccessControlHeadersBadBluemixSpace(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testAuth := "test"
	orgUUID := uuid.NewV4().String()

	testRequest.Header.Set(constants.AuthorizationHeader, testAuth)
	testRequest.Header.Set(constants.BluemixOrgHeader, orgUUID)
	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	//400 == badRequest
	if responseRecorder.Code != 400 {
		t.Errorf("Expected %d, received %d", 400, responseRecorder.Code)
	}
}

func TestAccessControlHeadersUnexpectedMethod(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	headerCheckHandler := AccessControlHeaders(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	testRequest.Method = http.MethodOptions
	headerCheckHandler.ServeHTTP(responseRecorder, testRequest)

	//200 == ok, invalid method returns without error
	if responseRecorder.Code != 200 {
		t.Errorf("Expected %d, received %d", 200, responseRecorder.Code)
	}
}
