//Package proxy .
// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Prox holds business logic for reverse proxy
type Prox struct {
	target   *url.URL
	revProxy *httputil.ReverseProxy
}

// New is a factory pattern to create new Proxy
func New(target string) *Prox {
	url, err := url.Parse(target)
	if err != nil {
		panic("bad string for proxy")
	}
	return &Prox{target: url, revProxy: httputil.NewSingleHostReverseProxy(url)}
}

// Handle implements http Handler for reverse proxy
func (p *Prox) Handle(w http.ResponseWriter, r *http.Request) {
	p.revProxy.ServeHTTP(w, r)
}
