// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.

package identity

import (
	"net/http"

	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/pep"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/provision"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/roles"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-config"
)

var (
	toggle = true
	config = configuration.Get()
)

// RoleCheck ...
func RoleCheck(handler http.Handler) http.Handler {
	if config.GetBool("feature_toggles.use_pep_service") == false {
		return roles.RoleCheck(handler)
	}

	return pep.RoleCheck(handler)
}

// VerifyIdentity ...
func VerifyIdentity(handler http.Handler) http.Handler {
	if config.GetBool("feature_toggles.use_pep_service") == false {
		return provision.VerifyIdentity(handler)
	}

	return pep.VerifyIdentity(handler)
}
