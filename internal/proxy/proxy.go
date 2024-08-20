package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"stark8/internal/utils"
	"strings"
	"time"

	"golang.org/x/exp/rand"
	"golang.org/x/net/html"
)

// Proxy represents a reverse proxy with its associated URL
type Proxy struct {
	url     *url.URL
	proxy   *httputil.ReverseProxy
	logoUrl *url.URL
	color   string
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

	logoURL, err := dummyCall(backendURL)
	if err != nil {
		log.Printf("Error during dummy call: %v\n", err)
	}
	randomColor := "from-" + utils.TailwindColors[rand.Intn(len(utils.TailwindColors))] + " to-" + utils.TailwindColors[rand.Intn(len(utils.TailwindColors))]
	p.hub[name] = Proxy{
		url:     backendURL,
		proxy:   proxy,
		logoUrl: logoURL,
		color:   randomColor,
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

// dummyCall performs a dummy HTTP call to find an element with "logo" in its class, id or link.
// It follows redirects and handles them appropriately.
func dummyCall(backendURL *url.URL) (*url.URL, error) {
	// Create a custom HTTP client to handle redirects manually
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("stopped after 3 redirects")
			}
			// If redirected, update backendURL to the new location
			newURL, err := url.Parse(req.URL.String())
			if err != nil {
				return err
			}
			*backendURL = *newURL
			return nil
		},
	}

	// Perform the request
	resp, err := client.Get(backendURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the HTML document
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	var logoURL *url.URL
	var firstPngURL *url.URL

	// Function to recursively find elements with "logo" or first PNG image
	var findLogo func(*html.Node)
	findLogo = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Check if the element has "logo" in its attributes
			for _, attr := range n.Attr {
				if strings.Contains(attr.Val, "logo") {
					for _, attr := range n.Attr {
						if strings.HasSuffix(attr.Val, ".png") && isInternalURL(attr.Val, backendURL) {
							logoURL, _ = backendURL.Parse(attr.Val)
							return
						}
					}
				}
			}
			// Check if it's an image and store the first PNG found
			if n.Data == "img" {
				for _, attr := range n.Attr {
					if strings.HasSuffix(attr.Val, ".png") && isInternalURL(attr.Val, backendURL) {
						if firstPngURL == nil {
							firstPngURL, _ = backendURL.Parse(attr.Val)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLogo(c)
		}
	}

	findLogo(doc)

	// Return the found logo URL or the first PNG URL
	if logoURL != nil {
		return logoURL, nil
	}
	if firstPngURL != nil {
		return firstPngURL, nil
	}
	return nil, nil
}

// isInternalURL checks if a given URL is internal based on the hostname of the backend URL.
func isInternalURL(link string, backendURL *url.URL) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	// If the parsed URL doesn't have a hostname, it is a relative URL and thus internal
	if parsedURL.Hostname() == "" {
		return true
	}

	// Otherwise, compare the hostname with the backend URL's hostname
	return parsedURL.Hostname() == backendURL.Hostname()
}
