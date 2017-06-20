package pep

import (
	"fmt"
	"net/http"
	"testing"

	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
)

var validGUID = "2561b2bc-45aa-4a39-914b-dc25d2c38521"
var validSpace = "766cd95e-4ddc-403f-a669-9c85d59ac7b1"
var baseURL = "/api/v2/secrets"

func TestExtractIAMDataListGET(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL, nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionListSecrets || resource != validSpace {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataListHEAD(t *testing.T) {
	req, err := http.NewRequest("HEAD", baseURL, nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionListSecrets || resource != validSpace {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataCreatePOST(t *testing.T) {
	req, err := http.NewRequest("POST", baseURL, nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionCreateSecret || resource != validSpace {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataDeleteDELETE(t *testing.T) {
	req, err := http.NewRequest("DELETE", baseURL+"/"+validGUID, nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionDeleteSecret || resource != validGUID {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataReadGET(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL+"/"+validGUID, nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionReadSecret || resource != validGUID {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataWrap(t *testing.T) {
	req, err := http.NewRequest("POST", baseURL+"/"+validGUID+"?action=wrap", nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionWrapSecret || resource != validGUID {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataUnwrap(t *testing.T) {
	req, err := http.NewRequest("POST", baseURL+"/"+validGUID+"?action=unwrap", nil)
	req.Header.Set(constants.BluemixSpaceHeader, validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != constants.ActionUnwrapSecret || resource != validGUID {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataBlankQuery(t *testing.T) {
	req, err := http.NewRequest("", "", nil)
	resource, action, err := ExtractIAMData(req)
	if action != "" || resource != "" {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err != nil {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataInvalidAction(t *testing.T) {
	req, err := http.NewRequest("POST", baseURL+"/"+validGUID+"?action=invalid", nil)
	req.Header.Set("BluemixSpaceHeader", validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != "" || resource != "" {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err.Error() != http.StatusText(http.StatusNotImplemented) {
		fmt.Printf("Error found: %s", err)
	}
}

func TestExtractIAMDataNoAction(t *testing.T) {
	req, err := http.NewRequest("POST", baseURL+"/"+validGUID, nil)
	req.Header.Set("BluemixSpaceHeader", validSpace)
	resource, action, err := ExtractIAMData(req)
	if action != "" || resource != "" {
		t.Errorf("Resource: %s and Action: %s", resource, action)
	}
	if err.Error() != http.StatusText(http.StatusNotImplemented) {
		fmt.Printf("Error found: %s", err)
	}
}
