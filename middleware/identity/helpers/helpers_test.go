// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package helpers

import (
	"encoding/base64"
	"testing"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

func TestDecodeSegment(t *testing.T) {
	testByte := "0000000"
	receivedByte, _ := DecodeSegment(testByte)
	expectedByte, _ := base64.URLEncoding.DecodeString("0000000=")
	receivedByteString := string(receivedByte[:2])
	expectedByteString := string(expectedByte[:2])

	if expectedByteString != receivedByteString {
		t.Errorf("Expected %s, received %s", expectedByte, receivedByte)
	}
}

func TestFindAlternativeUser(t *testing.T) {
	testIssuer := "test"
	receivedIssuer := FindAlternativeIssuer(testIssuer)
	if testIssuer != receivedIssuer {
		t.Errorf("Expected %s, received %s", testIssuer, receivedIssuer)
	}
}

func TestFindAlternativeUserForDallas(t *testing.T) {
	testIssuer := constants.DallasRegionSubDomain
	receivedIssuer := FindAlternativeIssuer(testIssuer)
	if testIssuer == receivedIssuer {
		t.Errorf("Expected, %s, received %s", constants.LondonRegionSubDomain, receivedIssuer)
	}
}

func TestFindAlternativeUserForLondon(t *testing.T) {
	testIssuer := constants.LondonRegionSubDomain
	receivedIssuer := FindAlternativeIssuer(testIssuer)
	if testIssuer == receivedIssuer {
		t.Errorf("Expected, %s, received %s", constants.DallasRegionSubDomain, receivedIssuer)
	}
}
