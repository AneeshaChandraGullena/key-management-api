//Package middleware .
// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
package middleware

import (
	provAC "github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/access_control/provision"
	stdAC "github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/access_control/standard"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/correlation"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/cache"
	"github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/ratelimiting"

	"github.com/justinas/alice"
)

// StdMiddlewareChain will register the middleware required for the API server
func StdMiddlewareChain(chain alice.Chain) alice.Chain {
	cache.InitCacheCompute()
	newChain := alice.New(stdAC.AccessControlHeaders, identity.RoleCheck)
	return chain.Extend(newChain)
}

// ProvisionMiddlewareChain will register the middleware required for the API server
func ProvisionMiddlewareChain(chain alice.Chain) alice.Chain {
	newChain := alice.New(provAC.AccessControlHeaders, identity.VerifyIdentity)
	return chain.Extend(newChain)
}

// BaseMiddlewareChain is the middleware installed at '/' and contains things we want for all routes to have
func BaseMiddlewareChain() alice.Chain {
	return alice.New(ratelimiting.Limit, correlation.ID)
}
