package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createStark8Request struct {
	Name     string `form:"name" binding:"required"`
	Protocol string `form:"protocol" binding:"required"`
	Port     int    `form:"port" binding:"required"`
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

	url := fmt.Sprintf("%s://%s.%s:%d", f.Protocol, c.Param("service"), c.Param("namespace"), f.Port)
	if err := s.proxyHub.NewProxy(f.Name, url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Print the form data

	// Return a JSON response with the form data
	c.JSON(http.StatusOK, f)
	// r := ginrender.New(c.Request.Context(), http.StatusOK, templates.CreateErrorBadge())
	// c.Render(http.StatusOK, r)

	// Uncomment this part when implementing the service handler logic
	// if err := s.proxyHub.NewProxy(f.Name, f.URL); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
}
