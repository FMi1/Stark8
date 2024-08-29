package api

import (
	"fmt"
	"net/http"
	"stark8/internal/store"
	"stark8/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        string `json:"user"`
}

type userLoginRequest struct {
	Username string `form:"username" binding:"required,alphanum,min=3,max=30"`
	Password string `form:"password" binding:"required,customPassword"`
}

type userSignupRequest struct {
	Username string `form:"username" binding:"required,alphanum,min=3,max=30"`
	Password string `form:"password" binding:"required,customPassword"`
}

func (s *Server) loginUserRequest(c *gin.Context) {
	// TODO: add service handler

	var req userLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req)

	if req.Username == "admin" && req.Password == "admin" {
		accessToken, err := s.TokenMaker.CreateToken(req.Username, s.config.TokenDuration)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//TODO: remove port
		subdomainHost := strings.Split(s.config.Hostname, ":")[0]
		c.SetCookie("stark8.token", accessToken, 3600, "/", "."+subdomainHost, false, true)
		//// HERE
		c.Header("HX-Redirect", "/") // Redirect to the home page
		c.JSON(http.StatusOK, loginUserResponse{
			AccessToken: accessToken,
			User:        req.Username,
		})
		return
	}

	user, err := s.db.GetUser(req.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := s.TokenMaker.CreateToken(req.Username, s.config.TokenDuration)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//TODO: remove port
	subdomainHost := strings.Split(s.config.Hostname, ":")[0]

	c.SetCookie("stark8.token", accessToken, 3600, "/", "."+subdomainHost, false, true)
	//// HERE
	c.Header("HX-Redirect", "/") // Redirect to the home page
	c.JSON(http.StatusOK, loginUserResponse{
		AccessToken: accessToken,
		User:        req.Username,
	})
}

func (s *Server) logoutUserRequest(c *gin.Context) {
	c.SetCookie("stark8.token", "", -1, "/", s.config.Hostname, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (s *Server) createUserRequest(c *gin.Context) {
	var req userSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Password = hashedPassword

	arg := store.UserParams{
		Username: req.Username,
		Password: hashedPassword,
	}
	err = s.db.CreateUser(arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
