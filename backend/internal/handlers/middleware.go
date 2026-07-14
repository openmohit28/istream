package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/auth"
)

const userIDKey = "userID"

func (s *Server) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", s.cfg.FrontendOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (s *Server) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || tokenString == "" {
			errorJSON(c, http.StatusUnauthorized, "missing or malformed authorization header")
			return
		}
		userID, err := auth.ParseToken(s.cfg.JWTSecret, tokenString)
		if err != nil {
			errorJSON(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		c.Set(userIDKey, userID)
		c.Next()
	}
}
