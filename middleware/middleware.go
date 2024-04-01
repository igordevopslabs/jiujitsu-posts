package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/models"
)

func RequireAuth(c *gin.Context) {

	//verificar se existe o header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Divide o header Authorization em "Bearer" e o valor do token.
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is Bearer token"})
		return
	}

	//Pegar valor do header Authorization da requisição
	tokenValue := strings.Split(c.Request.Header["Authorization"][0], " ")[1]

	if len(tokenValue) == 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Empty Token! Provide a valid token in your request.",
		})

		return
	}

	//Validar o Token
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["Authorization"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Checar a expiration Date
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//Encontrar o user (sub) no banco para fazer a autenticação
		var user models.User
		config.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}
