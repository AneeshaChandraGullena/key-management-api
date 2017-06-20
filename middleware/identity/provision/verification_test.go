package provision

import "testing"

func TestGetAuthURLFromTokenNoBearer(t *testing.T) {
	_, actualerror := getAuthURLFromToken("bear sdfpawoifjaslfd")
	expectederror := "token should contain 'bearer '"
	if actualerror.Error() != expectederror {
		t.Errorf("Expected %s, Received %s", expectederror, actualerror)
	}
}

func TestGetAuthURLFromTokenInvalidToken(t *testing.T) {
	_, actualerror := getAuthURLFromToken("bearer Testing.TestRing")
	expectederror := "invalid character 'M' looking for beginning of value"
	if actualerror.Error() != expectederror {
		t.Errorf("Expected %s, Received %s", expectederror, actualerror)
	}
}
