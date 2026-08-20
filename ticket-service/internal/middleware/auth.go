package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc { //Bu fonsiyon loginde oluşturduğumjuz JWT yi isteklerde kontrol edicek
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization") //Request Headerından aldık tkeni
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort() //işlemi durdur
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") { //bearer kontrolü
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized ",
			})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("role", claims["role"].(string))
		c.Set("user_id", claims["user_id"].(float64)) //json sayıları float64 müş
		c.Set("entity_id", claims["entity_id"].(float64))
		c.Next() //devam et
	}
}
func LimitBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			1<<20, //max 1mb
		)

		c.Next()
	}
}
