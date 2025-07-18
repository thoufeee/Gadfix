package controllers

import (
	"gadfix/utlis"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func StaffLogout(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid token"})
		return
	}

	refresh := strings.TrimPrefix(header, "Bearer ")
	token, err := jwt.ParseWithClaims(refresh, &utlis.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("refreshkey")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
		return
	}

	cliams, ok := token.Claims.(*utlis.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unathorized"})
		return
	}

	if err := utlis.DeleteRefresh(cliams.TokenID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly loged out"})
}
