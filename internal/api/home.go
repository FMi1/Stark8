package api

import (
	"net/http"
	"stark8/internal/templates"

	ginrender "stark8/internal/render"

	"github.com/gin-gonic/gin"
)

type getHomeRequest struct {
	PageNumber int `form:"page" binding:"min=1"`
}

func (s *Server) getHomeRequest(c *gin.Context) {
	req := &getHomeRequest{
		PageNumber: 1,
	}
	if err := c.ShouldBindQuery(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proxies, err := s.proxyHub.GetListProxy(10, req.PageNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r := ginrender.New(c.Request.Context(), http.StatusOK, templates.Home(proxies))

	c.Render(http.StatusOK, r)
	// c.JSON(http.StatusOK, services)
}
