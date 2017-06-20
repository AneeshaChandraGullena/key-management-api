// © Copyright 2017 IBM Corp. Licensed Materials – Property of IBM.

package helpers

import (
	"encoding/base64"
	"strings"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

// DecodeSegment decodes JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

// FindAlternativeIssuer will take a currentIssuer and return one in a different region from the original
func FindAlternativeIssuer(currentIssuer string) string {
	switch {
	// if dallas, try London
	case strings.Contains(currentIssuer, constants.DallasRegionSubDomain):
		return strings.Replace(currentIssuer, constants.DallasRegionSubDomain, constants.LondonRegionSubDomain, 1)
		// if london, try dallas
	case strings.Contains(currentIssuer, constants.LondonRegionSubDomain):
		return strings.Replace(currentIssuer, constants.LondonRegionSubDomain, constants.DallasRegionSubDomain, 1)
	}
	// default don't replace anything
	return currentIssuer
}
