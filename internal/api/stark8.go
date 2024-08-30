package api

import (
	"fmt"
	"net/http"
	ginrender "stark8/internal/render"
	"stark8/internal/templates"

	"github.com/gin-gonic/gin"
)

type createStark8Request struct {
	Name     string `form:"name" binding:"required"`
	Protocol string `form:"protocol" binding:"required"`
	Port     int    `form:"port" binding:"required"`
	Logo     string `form:"selectedLogoName" binding:"required"`
}

func (s *Server) createStark8Request(c *gin.Context) {
	// Define a struct to hold the form data
	var f createStark8Request
	// Bind the incoming form data to the struct
	if err := c.ShouldBind(&f); err != nil {
		// If there's an error in binding, return a bad request response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if a proxy with the same name already exists

	if _, err := s.proxyHub.GetProxy(f.Name); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("proxy with name %s already exists", f.Name)})
		return
	}
	url := fmt.Sprintf("%s://%s.%s:%d", f.Protocol, c.Param("service"), c.Param("namespace"), f.Port)
	if err := s.proxyHub.NewProxy(f.Name, url, f.Logo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &getStark8sRequest{
		PageNumber: 1,
	}
	proxies, err := s.proxyHub.GetListProxy(10, req.PageNumber)
	if err != nil {
		c.JSON(http.StatusCreated, f)
		return
	}
	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.Stark8sComponent(proxies))
	c.Render(http.StatusCreated, r)
}

type getStark8sRequest struct {
	PageNumber int `form:"page" binding:"min=1"`
}

func (s *Server) getStark8sRequest(c *gin.Context) {
	req := &getStark8sRequest{
		PageNumber: 1,
	}
	proxies, err := s.proxyHub.GetListProxy(10, req.PageNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.Stark8sComponent(proxies))
	c.Render(http.StatusOK, r)
}
