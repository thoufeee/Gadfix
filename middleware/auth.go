package middleware

import (
	"gadfix/utlis"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// auth check
func Auth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		tokenstr := strings.TrimPrefix(header, "Bearer ")
		jwtkey := []byte(os.Getenv("jwtkey"))

		token, err := jwt.ParseWithClaims(tokenstr, &utlis.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or token expired"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utlis.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token claims"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				c.Set("userid", claims.UserId)
				c.Set("email", claims.Email)
				c.Set("role", claims.Role)
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"err": "access denied"})
		c.Abort()
	}
}
