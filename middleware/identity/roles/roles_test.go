// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.

package roles

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

// testHandler only returns so that we can get the response from the middleware
func testHandler(writer http.ResponseWriter, request *http.Request) {
	return
}

func TestRoleCheckGoodPath(t *testing.T) {
	t.SkipNow()

	// throw away handler, it isn't used, but is required to call `roleCheck` function
	tempHandler := http.HandlerFunc(testHandler)

	// Need to dereference the tempHandler
	roleCheckHandler := RoleCheck(tempHandler)

	// create test request to pass to handler
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)

	// Token for keymaster@bg.vnet.ibm.com logged into AlchemyStaging Org in bluemix staging.
	testRequest.Header.Set(constants.AuthorizationHeader, "bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI5N2IzZTg4Ny05YWM1LTRjODYtYjBmOS04MDI2YmRiNTU3ODUiLCJzdWIiOiJjNTQ3ZmYwMS1hOGI4LTQzODctYmU1Ni1iNGIyNjJiYmRjMTQiLCJzY29wZSI6WyJvcGVuaWQiLCJ1YWEudXNlciIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiYzU0N2ZmMDEtYThiOC00Mzg3LWJlNTYtYjRiMjYyYmJkYzE0Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImVtYWlsIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImF1dGhfdGltZSI6MTQ4MTgyMjc4NiwicmV2X3NpZyI6IjI4ZWMzYjc3IiwiaWF0IjoxNDgxODIyNzg2LCJleHAiOjE0ODMwMzIzODYsImlzcyI6Imh0dHBzOi8vdWFhLnN0YWdlMS5uZy5ibHVlbWl4Lm5ldC9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjZiIsIm9wZW5pZCIsInVhYSIsImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCJdfQ.xoLIWF3kpjbHs-yese3e_f82-cN84_YxUQOe7Fgv7BY")
	// Space GUID for is for the Key Management space.
	testRequest.Header.Set(constants.BluemixSpaceHeader, "5187f101-39f8-4b7f-b961-5ef2cb65a7e9")

	responseRecorder := httptest.NewRecorder()

	roleCheckHandler.ServeHTTP(responseRecorder, testRequest)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Fail()
	}

	// keymaster@bg.vnet.ibm.com should have developer role in Key Management Space
	if role := testRequest.Header.Get(constants.BluemixUserRole); role != constants.RoleDeveloper {
		t.Fail()
	}
}

func TestRoleCheckBadPathNoSpaceHeader(t *testing.T) {
	t.SkipNow()

	// throw away handler, it isn't used, but is required to call `roleCheck` function
	tempHandler := http.HandlerFunc(testHandler)

	// Need to dereference the tempHandler
	roleCheckHandler := RoleCheck(tempHandler)

	// create test request to pass to handler
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)

	// Token for keymaster@bg.vnet.ibm.com logged into AlchemyStaging Org in bluemix staging.
	testRequest.Header.Set(constants.AuthorizationHeader, "bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI5N2IzZTg4Ny05YWM1LTRjODYtYjBmOS04MDI2YmRiNTU3ODUiLCJzdWIiOiJjNTQ3ZmYwMS1hOGI4LTQzODctYmU1Ni1iNGIyNjJiYmRjMTQiLCJzY29wZSI6WyJvcGVuaWQiLCJ1YWEudXNlciIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiYzU0N2ZmMDEtYThiOC00Mzg3LWJlNTYtYjRiMjYyYmJkYzE0Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImVtYWlsIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImF1dGhfdGltZSI6MTQ4MTgyMjc4NiwicmV2X3NpZyI6IjI4ZWMzYjc3IiwiaWF0IjoxNDgxODIyNzg2LCJleHAiOjE0ODMwMzIzODYsImlzcyI6Imh0dHBzOi8vdWFhLnN0YWdlMS5uZy5ibHVlbWl4Lm5ldC9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjZiIsIm9wZW5pZCIsInVhYSIsImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCJdfQ.xoLIWF3kpjbHs-yese3e_f82-cN84_YxUQOe7Fgv7BY")

	responseRecorder := httptest.NewRecorder()

	roleCheckHandler.ServeHTTP(responseRecorder, testRequest)

	// should return a 200 as it doesn't check if the space if valid.
	if status := responseRecorder.Code; status != http.StatusForbidden {
		t.Fail()
	}

	// should not have a role since there was no space
	if role := testRequest.Header.Get(constants.BluemixUserRole); role != "" {
		t.Fail()
	}

	expectedResponse := `{"metadata":{"collectionType":"application/vnd.ibm.kms.error+json","collectionTotal":1},"resources":[{"errorMsg":"Forbidden: User does not have access to provided space"}]}`
	if responseBody := responseRecorder.Body.String(); strings.Compare(strings.TrimSpace(responseBody), expectedResponse) != 0 {
		fmt.Println(strings.Compare(responseBody, expectedResponse))
		t.Errorf("Expected %s, recieved %s", expectedResponse, responseBody)
	}
}

func TestRoleCheckBadRequest(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	roleCheckHandler := RoleCheck(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	roleCheckHandler.ServeHTTP(responseRecorder, testRequest)

	//error 400 == badRequest
	if responseRecorder.Code != 400 {
		t.Errorf("Expected %d, received %d", 400, responseRecorder.Code)
	}
}

func TestPerformCheckBadHeaders(t *testing.T) {
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	authToken := "test"
	corrID := "test"
	spaceGUID := "test"

	_, err := performCheck(authToken, corrID, spaceGUID, responseRecorder, testRequest)

	expectedErr := "Bad Request: Please check authorization token and space guid"
	if expectedErr != err.Error() {
		t.Errorf("Expected %s, received %s", expectedErr, err.Error())
	}
}

func TestAcquireIDNoToken(t *testing.T) {
	authToken := ""

	_, err := acquireID(authToken)

	expectedErr := "Bad Request : Token should be prefixed with bearer"
	if expectedErr != err.Error() {
		t.Errorf("Expected %s, received %s", expectedErr, err.Error())
	}
}

func TestAcquireIDMisformedToken(t *testing.T) {
	authToken := "bearer test"

	_, err := acquireID(authToken)

	expectedErr := "Bad Request : missing or malformed token"
	if expectedErr != err.Error() {
		t.Errorf("Expected %s, received %s", expectedErr, err.Error())
	}
}

func TestAcquireIDGoodToken(t *testing.T) {
	authToken := "bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI5N2IzZTg4Ny05YWM1LTRjODYtYjBmOS04MDI2YmRiNTU3ODUiLCJzdWIiOiJjNTQ3ZmYwMS1hOGI4LTQzODctYmU1Ni1iNGIyNjJiYmRjMTQiLCJzY29wZSI6WyJvcGVuaWQiLCJ1YWEudXNlciIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiYzU0N2ZmMDEtYThiOC00Mzg3LWJlNTYtYjRiMjYyYmJkYzE0Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImVtYWlsIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImF1dGhfdGltZSI6MTQ4MTgyMjc4NiwicmV2X3NpZyI6IjI4ZWMzYjc3IiwiaWF0IjoxNDgxODIyNzg2LCJleHAiOjE0ODMwMzIzODYsImlzcyI6Imh0dHBzOi8vdWFhLnN0YWdlMS5uZy5ibHVlbWl4Lm5ldC9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjZiIsIm9wZW5pZCIsInVhYSIsImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCJdfQ.xoLIWF3kpjbHs-yese3e_f82-cN84_YxUQOe7Fgv7BY"

	receivedUserID, _ := acquireID(authToken)

	expectedUserID := "c547ff01-a8b8-4387-be56-b4b262bbdc14"
	if expectedUserID != receivedUserID {
		t.Errorf("Expected %s, received %s", expectedUserID, receivedUserID)
	}
}

func TestAcquireIDBadToken(t *testing.T) {
	authToken := "bearer eyJhbGciOiJIUzI1NiJ9.test.xoLIWF3kpjbHs-yese3e_f82-cN84_YxUQOe7Fgv7BY"

	_, err := acquireID(authToken)

	expectedErr := "invalid character 'µ' looking for beginning of value"
	if expectedErr != err.Error() {
		t.Errorf("Expected %s, received %s", expectedErr, err.Error())
	}
}

func TestMakeCacheKey(t *testing.T) {
	authToken := "auth"
	space := "test"

	expectedReturn := "authtest"
	if actualReturn := makeCacheKey(authToken, space); actualReturn != expectedReturn {
		t.Errorf("Expected %s, received %s", expectedReturn, actualReturn)
	}
}
