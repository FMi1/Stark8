package api

import (
	"net/http"
	ginrender "stark8/internal/render"
	"stark8/internal/templates"

	"github.com/gin-gonic/gin"
)

func (s *Server) getModalStruct(c *gin.Context) {
	// TODO: add service handler
	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.ModalComponent())
	c.Render(http.StatusOK, r)
}

func (s *Server) getModalNamespacesRequest(c *gin.Context) {
	// Add the get of namespaces
	namespaces, err := s.k8sClientset.GetNamespaces(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.ModalBodyNamespacesComponent(namespaces))
	c.Render(http.StatusOK, r)
}

func (s *Server) getModalServicesRequest(c *gin.Context) {
	// TODO: add service handler
	namespace := c.Param("namespace")

	services, err := s.k8sClientset.GetServices(c.Request.Context(), namespace)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.ModalBodyServicesComponent(namespace, services))
	c.Render(http.StatusOK, r)
}

func (s *Server) getModalSettingsRequest(c *gin.Context) {
	// TODO: add service handler
	namespace := c.Param("namespace")
	service := c.Param("service")
	errors := make(map[string]bool)
	values := make(map[string]string)
	ports, err := s.k8sClientset.GetServicePorts(c.Request.Context(), namespace, service)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.ModalBodySettingsComponent(namespace, service, ports, errors, values))
	c.Render(http.StatusOK, r)
}

func (s *Server) getLoginUserRequest(c *gin.Context) {
	// TODO: add service handler
	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.Login())
	c.Render(http.StatusOK, r)

}
