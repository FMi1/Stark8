package proxy

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Proxy represents a reverse proxy with its associated URL
type Proxy struct {
	url   *url.URL
	proxy *httputil.ReverseProxy
	// logoUrl *url.URL
}

// redirectStatusCodes is a package-level variable containing HTTP redirection status codes
var redirectStatusCodes = map[int]struct{}{
	http.StatusMultipleChoices:   {},
	http.StatusMovedPermanently:  {},
	http.StatusFound:             {},
	http.StatusSeeOther:          {},
	http.StatusNotModified:       {},
	http.StatusUseProxy:          {},
	http.StatusTemporaryRedirect: {},
	http.StatusPermanentRedirect: {},
}

// isRedirect checks if a status code is in the redirection set
func isRedirect(statusCode int) bool {
	_, exists := redirectStatusCodes[statusCode]
	return exists
}

// NewProxy creates a new proxy and adds it to the ProxyHub
func (p *ProxyHub) NewProxy(name string, urlStr string) error {
	backendURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Error parsing URL: %v\n", err)
		return err
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()      // 1
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // 2

	proxy := &httputil.ReverseProxy{
		Transport: transport, // 3
		Rewrite: func(r *httputil.ProxyRequest) {
			rewriteRequest(r, backendURL)
		},
		ModifyResponse: func(r *http.Response) error {
			return modifyResponse(r, name, backendURL, p.hostname)
		},
	}

	p.hub[name] = Proxy{
		url:   backendURL,
		proxy: proxy,
	}

	return nil
}

// GetProxy retrieves a proxy by its name

// rewriteRequest sets up the request headers and URL for the reverse proxy
func rewriteRequest(r *httputil.ProxyRequest, backendURL *url.URL) {
	r.SetXForwarded()
	r.SetURL(backendURL)
	fmt.Printf("%s - - [%s] \"%s %s %s\" %dus\n", strings.Split(r.In.RemoteAddr, ":")[0], time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.In.Method, r.Out.URL.Path, r.In.Proto, time.Since(time.Now()))
	fmt.Printf("%s - - [%s] Proxy Request to: %s %s %dus\n", strings.Split(r.In.RemoteAddr, ":")[0], time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.In.Method, r.Out.URL, time.Since(time.Now()))

	fmt.Println("Rewriting request to:", r.Out.URL)
	r.Out.Header = r.In.Header
}

// modifyResponse handles response modification for redirects
func modifyResponse(r *http.Response, name string, backendURL *url.URL, hostname string) error {
	if isRedirect(r.StatusCode) {
		location := r.Header.Get("Location")
		locationURL, err := url.Parse(location)
		if err != nil {
			log.Printf("Error parsing Location header: %v\n", err)
			return err
		}

		if locationURL.Host == backendURL.Host {
			locationURL.Host = name + "." + hostname
			locationURL.Scheme = "http"
			r.Header.Set("Location", locationURL.String())
		}
	}

	return nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
