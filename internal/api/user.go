package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type loginUserRequest struct {
	Username string `form:"username" binding:"required,alphanum,min=3,max=30"`
	Password string `form:"password" binding:"required,min=3,max=30"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        string `json:"user"`
}

func (s *Server) loginUserRequest(c *gin.Context) {
	// TODO: add service handler
	var req loginUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/////// test
	if req.Username == "admin" && req.Password == "admin" {
		accessToken, err := s.TokenMaker.CreateToken(req.Username, s.config.TokenDuration)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//TODO: remove port
		subdomainHost := strings.Split(s.config.Hostname, ":")[0]
		c.SetCookie("stark8_token", accessToken, 3600, "/", "."+subdomainHost, false, true)
		//// HERE
		c.Header("HX-Redirect", "/") // Redirect to the home page
		c.JSON(http.StatusOK, loginUserResponse{
			AccessToken: accessToken,
			User:        req.Username,
		})
		return
	}
	///////

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}
