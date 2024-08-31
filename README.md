# Stark8 - Kubernetes Native Web Application Secure Access

Stark8 is a project born for learning purpose, it's a GOTTH (Go, Templ, Tailwind CSS, HTMX) stack example that uses the official Kubernetes Go client to interact with the Kubernetes API.

## Problem Overview

In typical Kubernetes deployments, services are often exposed through an Ingress resource. This approach can lead to challenges, such as:

- **Unrestricted Visibility**: All services become visible to the public, potentially exposing sensitive internal services.
- **Difficult Management**: Managing access to each service individually can become complex and error-prone, especially when administrators must remember the details of every service and its corresponding URL.

Stark8 solves these problems by providing a centralized, secure access point that controls which services are exposed, and only to authenticated users.

---
**Stark8** is a Kubernetes-native reverse proxy that enables secure access to services running within a Kubernetes cluster via a single point of entry. By integrating authentication mechanisms, Stark8 ensures that only authorized users can reach internal services, making it an ideal solution for managing access to applications in a distributed environment.

---

## Key Features

- **Unified Access Point**: Provides a single entry point for accessing multiple services within a Kubernetes cluster.
- **Authentication and Authorization**: Enforces user authentication and authorization to control access to internal services.
- **Kubernetes-Native**: Seamlessly integrates with Kubernetes, utilizing native DNS, service discovery, and routing capabilities.
- **Zero Trust Access**: Ensures that users must authenticate and be authorized to access any internal services, regardless of network trust.
- **Dynamic Routing**: Automatically routes incoming traffic to the appropriate service based on pre-defined rules.
- **UI for Access Management**: Provides a user interface (UI) where administrators can easily create and manage access to Kubernetes services.
- **Service-Specific URLs**: Each exposed service can be accessed via a dedicated URL in the format:  
  `<service_name>.<stark8_host>`

---

## How It Works

Stark8 operates as a reverse proxy within your Kubernetes cluster. It authenticates users via a configured OAuth2. This allows external users to securely access internal services without exposing them directly to the internet.

1. **User Authentication**: When a user attempts to access a service, Stark8 prompts for authentication through the configured provider.
2. **Service Dashboard**: After successful authentication, the user is presented with an interface displaying all the services they have access to. From this dashboard, users can easily select and access the services available to them within the Kubernetes cluster.
3. **Service Access**: Upon selecting a service, Stark8 routes the request to the corresponding service within the cluster, ensuring secure access.
4. **Dynamic Routing**: Routing is handled based on predefined paths or through service-specific URLs in the format:  
   `<service_name>.<stark8_host>`


---

## Technologies Used

It uses [Gin](https://github.com/gin-gonic/gin) for the web framework, [Viper](https://github.com/spf13/viper) for configuration, [Paseto](https://github.com/o1egl/paseto) for token authentication, [Tailwind CSS](https://tailwindcss.com/) for styles and [Htmx](https://htmx.org/) for ajax requests.


---

# Roadmap

This is early development version. I am currently considering:

- [x] HTTPS with self sign certs
- [ ] Edit and Delete service
- [ ] Modify manifests adding PVC, ingress
- [ ] RBAC
- [ ] User Management
- [ ] Port Agnostic
- [ ] Helm Chart
- [ ] OIDC


---
## Installation
