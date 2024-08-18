package api

import (
	"net/http"
	"stark8/internal/token"

	"github.com/gin-gonic/gin"
)

const authorizationPayloadKey = "authorization_payload"

func authMiddleware(TokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authorizationHeader := c.GetHeader(authorizationHeaderKey)
		cookie, err := c.Request.Cookie("stark8_token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/users/login")
			c.Abort()
			return
		}
		if len(cookie.Value) == 0 {
			c.Redirect(http.StatusSeeOther, "/users/login")
			c.Abort()
			return
		}

		payload, err := TokenMaker.VerifyToken(cookie.Value)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/users/login")
			c.Abort()
			return
		}
		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
	// TODO: add auth handler
}
