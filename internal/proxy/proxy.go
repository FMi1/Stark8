package proxy

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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
	r.SetXForwarded()
	r.SetURL(backendURL)
	r.Out.Header = r.In.Header
}

// modifyResponse handles response modification for redirects
func modifyResponse(r *http.Response, name string, backendURL *url.URL, hostname string) error {
	if isRedirect(r.StatusCode) {
		location := r.Header.Get("Location")
		locationURL, err := url.Parse(location)
		log.Println("Redirect Location:", locationURL.String(), "Backend:", backendURL.String())
		if err != nil {
			log.Printf("Error parsing Location header: %v\n", err)
			return err
		}
		if locationURL.Scheme != backendURL.Scheme {
			backendURL.Scheme = locationURL.Scheme
		}
		if locationURL.Host == backendURL.Host {
			locationURL.Host = name + "." + hostname
			locationURL.Scheme = "https" // Cambia schema in base a TLS
			r.Header.Set("Location", locationURL.String())
		}
	}
	//// TODO: rimuovi l'attributo Secure perche siamo in http per ora
	// Modifica i cookie `Set-Cookie` per rimuovere l'attributo `Secure`
	setCookies := r.Header.Values("Set-Cookie")
	for i, cookie := range setCookies {
		// Se il cookie contiene l'attributo Secure, lo rimuove
		if strings.Contains(cookie, "Secure") {
			modifiedCookie := strings.ReplaceAll(cookie, "; Secure", "")
			setCookies[i] = modifiedCookie
		}
	}

	// Se è stata fatta una modifica, aggiorna l'header `Set-Cookie`
	if len(setCookies) > 0 {
		r.Header["Set-Cookie"] = setCookies
	}

	return nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
