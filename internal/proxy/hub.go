package proxy

import (
	"fmt"
	"sort"
)

// ProxyHub manages a collection of proxies
type ProxyHub struct {
	hub      map[string]Proxy
	hostname string
}

// NewProxyHub initializes and returns a new ProxyHub
func NewProxyHub(hostname string) *ProxyHub {
	return &ProxyHub{
		hub:      make(map[string]Proxy),
		hostname: hostname,
	}
}

func (p *ProxyHub) GetProxy(name string) (*Proxy, error) {
	if proxy, ok := p.hub[name]; ok {
		return &proxy, nil
	}
	return nil, fmt.Errorf("proxy not found")
}

// GetListProxy retrieves a list of proxies within a given range
// The method is idempotent.
func (p *ProxyHub) GetListProxy(elementsPerPage int, pageNumber int) (map[string]map[string]string, error) {
	pageSize := elementsPerPage
	offset := (pageNumber - 1) * pageSize
	keys := p.sortedKeys(offset, offset+pageSize)
	proxies := make(map[string]map[string]string)
	for _, key := range keys {
		proxy := p.hub[key]
		if proxy.logoUrl == nil {
			proxies[key] = map[string]string{
				"internalURL": proxy.url.String(),
				"externalURL": fmt.Sprintf("//%s.%s", key, p.hostname),
				"color":       proxy.color,
			}
		} else {
			proxies[key] = map[string]string{
				"internalURL": proxy.url.String(),
				"externalURL": fmt.Sprintf("//%s.%s", key, p.hostname),
				"color":       proxy.color,
				"logoURL":     proxy.logoUrl.String(),
			}
		}

	}
	return proxies, nil
}

// sortedKeys returns a slice of keys sorted in lexicographic order.
// The method is idempotent.
func (p *ProxyHub) sortedKeys(offset, endIndex int) []string {
	keys := make([]string, 0, len(p.hub))
	for key := range p.hub {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if endIndex > len(keys) {
		endIndex = len(keys)
	}
	return keys[offset:endIndex]
}
