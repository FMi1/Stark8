package proxy

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy represents a reverse proxy with its associated URL
type Proxy struct {
	url   *url.URL
	proxy *httputil.ReverseProxy
	logo  string
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
func (p *ProxyHub) NewProxy(name string, urlStr string, logo string) error {
	backendURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Error parsing URL: %v\n", err)
		return err
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	proxy := &httputil.ReverseProxy{
		Transport: transport,
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
		logo:  logo,
	}

	return nil
}

// GetProxy retrieves a proxy by its name

// rewriteRequest sets up the request headers and URL for the reverse proxy
func rewriteRequest(r *httputil.ProxyRequest, backendURL *url.URL) {
	r.SetURL(backendURL)
	r.Out.Header = r.In.Header
	r.SetXForwarded()
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
		if locationURL.Host != "" {
			if locationURL.Host == backendURL.Host {
				locationURL.Host = name + "." + hostname
			}
			if locationURL.Scheme != "" {
				if locationURL.Scheme != backendURL.Scheme {
					backendURL.Scheme = locationURL.Scheme
				}
				locationURL.Scheme = "https"
			}
			r.Header.Set("Location", locationURL.String())
		}
	}

	return nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
