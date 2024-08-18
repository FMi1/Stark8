package k8sclient

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client is a structure that holds a Kubernetes client
type Client struct {
	client *kubernetes.Clientset
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// GetNamespaces returns all namespaces in the cluster
func (c *Client) GetNamespaces(ctx context.Context) ([]string, error) {
	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, namespace := range namespaces.Items {
		out = append(out, namespace.Name)
	}
	return out, nil
}

// GetServices returns all services in the given namespace
func (c *Client) GetServices(ctx context.Context, namespace string) ([]string, error) {
	services, err := c.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, service := range services.Items {
		out = append(out, service.Name)
	}
	return out, nil
}

// GetService returns a service struct given the service name and namespace
func (c *Client) GetServicePorts(ctx context.Context, namespace, name string) ([]v1.ServicePort, error) {
	service, err := c.client.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return service.Spec.Ports, nil
}
