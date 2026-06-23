package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type GatewayProxy struct {
	authServiceURL *url.URL
	authProxy      *httputil.ReverseProxy
}

func NewGatewayProxy(authUrl string) (*GatewayProxy, error) {
	parsedAuthURL, err := url.Parse(authUrl)

	if err != nil {
		return nil, err
	}

	return &GatewayProxy{
		authServiceURL: parsedAuthURL,
		authProxy:      httputil.NewSingleHostReverseProxy(parsedAuthURL),
	}, nil

}

func (gp *GatewayProxy) RouteRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path, "r.URL.Path")

	if strings.HasPrefix(r.URL.Path, "/auth") {
		fmt.Println(r.Host, "2.r.host")

		r.Header.Set("X-Forwarded-Host", r.Host) // req header host to downstream service accepts it.
		fmt.Println(gp.authServiceURL.Host, "gp.authServiceURL.Host")
		r.Host = gp.authServiceURL.Host

		gp.authProxy.ServeHTTP(w, r)
		return
	}

	// Default fallback to route not found
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Route not found in API Gateway"}`))

}
