package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pmoorani/tasksAPI/config"
	"github.com/pmoorani/tasksAPI/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetCurrentURLPath(c *gin.Context) string {
	return c.Request.URL.Path
}

func TokenAuthMiddleware() gin.HandlerFunc {
	SecretKey := config.New().SecretKey
	return func(c *gin.Context) {
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "DELETE" {
			c.Next()
			return
		}
		tokenOK := true

		//List of endpoints that doesn't require Authentication
		notAuth := []string{"/api/login", "/api/register"}
		requestPath := GetCurrentURLPath(c)

		// Check if request does not need Authentication
		for _, value := range notAuth {
			if value == requestPath {
				c.Next()
				return
			}
		}

		// Get Authorization Header from Request Headers
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg":     "Token missing",
				"success": 0,
			})
			return
		}

		// Check if token came in format `JWT {token-body}`,
		splitted := strings.Split(authHeader, " ")
		if len(splitted) != 2 {
			tokenOK = false
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"msg":     "Malformed Auth Token!",
				"success": 0,
			})
			return
		}

		if len(splitted) == 2 {
			prefixJWT := splitted[0]
			if prefixJWT != "JWT" {
				tokenOK = false
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"msg":     "Malformed Auth Token",
					"success": 0,
				})
				return
			}
		}

		if tokenOK {
			tokenPart := splitted[1]

			claims := &models.Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(SecretKey), nil
			})
			fmt.Println(token.Valid, "claims=", claims, "err=", err)
			fmt.Println("token = ", token)
			if !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"msg":     "Token invalid",
					"success": 0,
				})
				return
			}

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					c.AbortWithStatus(http.StatusUnauthorized)
				}
				c.AbortWithStatus(http.StatusBadRequest)
			}

			claims.IsAuthenticated = true
			c.Set("claims", claims)
			c.Next()
		}
	}
}
