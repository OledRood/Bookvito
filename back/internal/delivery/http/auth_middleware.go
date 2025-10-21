package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("AuthMiddleware: Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("AuthMiddleware: Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что метод подписи HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("AuthMiddleware: unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			// Логируем конкретную ошибку парсинга токена
			log.Printf("AuthMiddleware: Invalid token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			log.Println("AuthMiddleware: Invalid token claims")
			return
		}

		userID, ok := claims["userId"].(string)
		if !ok || userID == "" {
			log.Println("AuthMiddleware: userId not found or is not a string in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userId not found in token"})
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			log.Println("AuthMiddleware: role not found or is not a string in token")
		}

		c.Set("userId", userID)
		c.Set("role", userRole)
		c.Next()
	}
}
