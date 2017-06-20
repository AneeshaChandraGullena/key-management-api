package ratelimiting

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler(writer http.ResponseWriter, request *http.Request) {
	return
}

//test that ratelimiter works correctly
func TestLimitCorrectRate(t *testing.T) {
	tempHandler := http.HandlerFunc(testHandler)
	checkHandler := Limit(tempHandler)
	testRequest, _ := http.NewRequest(http.MethodHead, "/", nil)
	responseRecorder := httptest.NewRecorder()

	checkHandler.ServeHTTP(responseRecorder, testRequest)

	//200 == ok, Limit() handles ratelimit correctly
	if httpCode := responseRecorder.Code; httpCode != 200 {
		t.Errorf("Expected %d, received %d", 200, httpCode)
	}
}
