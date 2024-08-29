package api

import (
	"net/http"
	"stark8/internal/token"

	"github.com/gin-gonic/gin"
)

const authorizationPayloadKey = "authorization_payload"

func authMiddleware(TokenMaker token.Maker, host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: fix the scheme(http)
		// authorizationHeader := c.GetHeader(authorizationHeaderKey)
		cookie, err := c.Request.Cookie("stark8.token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "http://"+host+"/users/login")
			c.Abort()
			return
		}
		if len(cookie.Value) == 0 {
			c.Redirect(http.StatusTemporaryRedirect, "http://"+host+"/users/login")
			c.Abort()
			return
		}

		payload, err := TokenMaker.VerifyToken(cookie.Value)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "http://"+host+"/users/login")
			c.Abort()
			return
		}
		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}
