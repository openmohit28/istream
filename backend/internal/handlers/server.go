package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/config"
	"istream/backend/internal/store"
)

type Server struct {
	cfg   config.Config
	users *store.Users
}

func NewServer(db *sql.DB, cfg config.Config) http.Handler {
	s := &Server{
		cfg:   cfg,
		users: &store.Users{DB: db},
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), s.cors())

	api := r.Group("/api")
	api.GET("/health", s.handleHealth)

	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", s.handleRegister)
	authRoutes.POST("/login", s.handleLogin)
	authRoutes.GET("/me", s.requireAuth(), s.handleMe)

	return r
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func errorJSON(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}