package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{}
}

func (r *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ae := r.Header.Get("Api-Endpoint")
		if len(ae) < 1 {
			next.ServeHTTP(w, r)
			return
		}

		ap := r.Header.Get("Api-Port")

		scheme := r.Header.Get("Api-Scheme")
		if scheme == "" {
			scheme = "http"
		}
		if ap == "" && scheme == "https" {
			ap = "443"
		}

		if ap == "" {
			ap = "8955"
		}

		link := fmt.Sprintf("%s://%s:%s", scheme, ae, ap)
		uri, _ := url.Parse(link)

		if uri.Host == r.Host {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Set("Reverse-Proxy", "true")

		proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
			r.URL.Scheme = uri.Scheme
			r.URL.Host = uri.Host
			r.URL.Path = uri.Path + r.URL.Path
			r.Host = uri.Host
		}}

		proxy.ServeHTTP(w, r)
	})
}
