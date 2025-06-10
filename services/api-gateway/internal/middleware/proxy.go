package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(target string) http.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic("Invalid target URL: " + target)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		if strings.HasPrefix(req.URL.Path, "/users") {
			req.URL.Path = singleJoiningSlash(targetURL.Path, strings.TrimPrefix(req.URL.Path, "/users"))
		} else {
			req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
