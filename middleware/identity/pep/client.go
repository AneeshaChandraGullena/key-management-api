// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.

package pep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	pepproto "github.ibm.com/Alchemy-Key-Protect/key-management-api/middleware/identity/protocols"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-config"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-consts"
	"github.ibm.com/Alchemy-Key-Protect/kp-go-models/collections"
	keyerror "github.ibm.com/Alchemy-Key-Protect/kp-go-models/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/url"
	"regexp"
	"strings"
)

var (
	config = configuration.Get()
)

// client defines PEP client connection envelope
type client struct {
	client pepproto.PEPAPIClient
	conn   *grpc.ClientConn
}

// newPEPClient creates connection to gRPC and returns
// PEP client along with gRPC connection
func newPEPClient() (*client, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", config.GetString("pepService.ipv4_address"), config.GetInt("pepService.port")),
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second*time.Duration(config.GetInt("timeouts.grpcTimeout"))),
	)

	if err != nil {
		return &client{}, err
	}

	return &client{
		pepproto.NewPEPAPIClient(conn),
		conn,
	}, nil
}

// RoleCheck ...
func RoleCheck(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c, err := newPEPClient()
		if err != nil {
			errorCollection := collections.NewErrorCollection().Append(err)
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}
		defer c.conn.Close()
		client := c.client

		resource, action, err := ExtractIAMData(request)
		if err != nil {
			httperror := keyerror.ConvertError(err)
			errorCollection := collections.NewErrorCollection().Append(httperror)
			writer.WriteHeader(int(httperror.StatusCode))
			json.NewEncoder(writer).Encode(errorCollection)
			return
		} else if resource == "" || action == "" {
			notImplemented := errors.New(http.StatusText(http.StatusNotImplemented))
			errorCollection := collections.NewErrorCollection().Append(notImplemented)
			writer.WriteHeader(http.StatusNotImplemented)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		authRequest := &pepproto.AuthzRequest{

			Header: map[string]string{
				constants.AuthorizationHeader: request.Header.Get(constants.AuthorizationHeader),
				constants.BluemixSpaceHeader:  request.Header.Get(constants.BluemixSpaceHeader),
				constants.CorrelationIDHeader: request.Header.Get(constants.CorrelationIDHeader),
			},

			Resource: resource,
			Action:   action,
		}

		reply, err := client.CheckAuthorization(context.Background(), authRequest)
		if err != nil {
			httperror := keyerror.ConvertError(err)
			errorCollection := collections.NewErrorCollection().Append(httperror)
			writer.WriteHeader(int(httperror.StatusCode))
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		if reply.Error != "" {
			httperror := keyerror.ConvertError(reply.Error)
			errorCollection := collections.NewErrorCollection().Append(httperror)
			writer.WriteHeader(int(httperror.StatusCode))
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		/* Have to fill in the headers based on UAA or IAM.
		   If the token provided was UAA:
		      the Role (for action validation in k-m-core) & UserId will be in the reply
		   If the token provided was IAM:
		      all role validation for the action has taken place and the Allowed
		      response determines whether to forge on. */
		if reply.Role != "" {
			request.Header.Set(constants.BluemixUserRole, reply.Role)
			request.Header.Set(constants.UserIDHeader, reply.UserId)
		} else if reply.Allowed {
			request.Header.Set(constants.BluemixUserRole, constants.RoleManager)
			request.Header.Set(constants.UserIDHeader, reply.UserId)
		} else {
			unauthorized := errors.New(http.StatusText(http.StatusForbidden) + ": user does not have access to provided space")
			errorCollection := collections.NewErrorCollection().Append(unauthorized)
			writer.WriteHeader(http.StatusForbidden)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		handler.ServeHTTP(writer, request)
	})
}

// VerifyIdentity ...
func VerifyIdentity(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c, err := newPEPClient()
		if err != nil {
			errorCollection := collections.NewErrorCollection().Append(err)
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}
		defer c.conn.Close()
		client := c.client

		header := &pepproto.UAAHeaders{
			Token:         request.Header.Get(constants.AuthorizationHeader),
			CorrelationId: request.Header.Get(constants.CorrelationIDHeader),
		}

		reply, err := client.UAAVerifyIdentity(context.Background(), header)
		if err != nil {
			errorCollection := collections.NewErrorCollection().Append(err)
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		if reply.Error != "" {
			errorCollection := collections.NewErrorCollection().Append(errors.New(reply.Error))
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		if reply.Authenticated == false {
			unauthorized := errors.New(http.StatusText(http.StatusUnauthorized) + ": request requires valid authorization")
			errorCollection := collections.NewErrorCollection().Append(unauthorized)
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(errorCollection)
			return
		}

		handler.ServeHTTP(writer, request)
	})
}

//ExtractIAMData is a helper method to retrieve the resource (key or space), and the IAM action
//Takes in the request
func ExtractIAMData(request *http.Request) (resource string, action string, err error) {

	spaceID := request.Header.Get(constants.BluemixSpaceHeader)
	method := request.Method
	path := strings.ToLower(request.URL.Path)

	parsedQuery, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		return "", "", err
	}

	r, err := regexp.Compile("[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}")
	if err != nil {
		return "", "", err
	}

	if textguid := r.FindString(path); textguid != "" {
		resource = textguid
		if wrapunwrap := parsedQuery.Get("action"); wrapunwrap != "" {
			if wrapunwrap == "wrap" {
				action = constants.ActionWrapSecret
			} else if wrapunwrap == "unwrap" {
				action = constants.ActionUnwrapSecret
			} else {
				//Only wrap/unwrap allowed atm.
				return "", "", errors.New(http.StatusText(http.StatusNotImplemented))
			}
		} else if method == "GET" {
			action = constants.ActionReadSecret
		} else if method == "DELETE" {
			action = constants.ActionDeleteSecret
		} else {
			//Only POST (wrap/unwrap), GET & DELETE have an <id> in the path
			return "", "", errors.New(http.StatusText(http.StatusNotImplemented))
		}
	} else if path == "/api/v2/secrets" || path == "/api/v2/keys"{
		resource = spaceID
		if method == "POST" {
			action = constants.ActionCreateSecret
		} else {
			action = constants.ActionListSecrets
		}
	} else {
		// means someone is putting in a bad guid or random stuff
		return "", "", errors.New(http.StatusText(http.StatusBadRequest))
	}

	return resource, action, nil
}
