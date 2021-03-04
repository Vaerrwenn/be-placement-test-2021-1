package middleware

import (
	"b-pay/config/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthJWT is a middleware for protected APIs. Checks whether the User who's
// trying to use an API is authenticated or not.
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the 'authorization' from the Header.
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No token in header",
			})
			c.Abort()
			return
		}

		// Check whether the Token is valid.
		jwtWrapper := auth.JwtWrapper{
			SecretKey: os.Getenv("JWT_SECRET"),
			Issuer:    "AuthService",
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
