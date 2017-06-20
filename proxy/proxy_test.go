//Package proxy .
// © Copyright 2016 IBM Corp. Licensed Materials – Property of IBM.
package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"
)

var urlTests = []struct {
	in     string
	scheme string
	host   string
	path   string
}{
	{"http://localhost:8989", "http", "localhost:8989", ""},
	{"http://localhost:8989/admin/v1/healthcheck", "http", "localhost:8989", "/admin/v1/healthcheck"},
	{"http://localhost:8989/admin/v1", "http", "localhost:8989", "/admin/v1"},
}

func TestNew(t *testing.T) {
	for _, scenario := range urlTests {
		result := New(scenario.in)
		if result.target.Host != scenario.host {
			t.Errorf("New(%s) => %s want %s", scenario.in, result.target.Host, scenario.host)
		}
		if result.target.Scheme != scenario.scheme {
			t.Errorf("New(%s) => %s want %s", scenario.in, result.target.Scheme, scenario.scheme)
		}
		if result.target.Path != scenario.path {
			t.Errorf("New(%s) => %s want %s", scenario.in, result.target.Path, scenario.path)
		}

	}
}

// TestHandle validates proxy.
// derived from  https://github.com/golang/go/blob/master/src/net/http/httputil/reverseproxy_test.go
func TestHandle(t *testing.T) {
	const backendResponse = "I am the backend"
	const backendStatus = 404
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.FormValue("mode") == "hangup" {
			c, _, _ := w.(http.Hijacker).Hijack()
			c.Close()
			return
		}
		if len(r.TransferEncoding) > 0 {
			t.Errorf("backend got unexpected TransferEncoding: %v", r.TransferEncoding)
		}
		if r.Header.Get("X-Forwarded-For") == "" {
			t.Errorf("didn't get X-Forwarded-For header")
		}
		// if r.Header.Get("X-MyReverseProxy") == "" {
		// 	t.Errorf("didn't get X-MyReverseProxy header")
		// }
		if !strings.Contains(runtime.Version(), "1.5") {
			if c := r.Header.Get("Proxy-Connection"); c != "" {
				t.Errorf("handler got Proxy-Connection header value %q", c)
			}
		}
		if g, e := r.Host, "some-name"; g != e {
			t.Errorf("backend got Host header %q, want %q", g, e)
		}
		w.Header().Set("X-Foo", "bar")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	// setup Frontend using admin proxy
	mux := http.NewServeMux()
	reverseProxy := &Prox{}
	reverseProxy = New(backendURL.String())
	mux.HandleFunc("/", reverseProxy.Handle)

	frontend := httptest.NewServer(mux)
	defer frontend.Close()

	getReq, _ := http.NewRequest("GET", frontend.URL, nil)
	getReq.Host = "some-name"
	getReq.Header.Set("Proxy-Connection", "should be deleted")
	getReq.Close = true
	res, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if g, e := res.StatusCode, backendStatus; g != e {
		t.Errorf("got res.StatusCode %d; expected %d", g, e)
	}
	if g, e := res.Header.Get("X-Foo"), "bar"; g != e {
		t.Errorf("got X-Foo %q; expected %q", g, e)
	}
}
