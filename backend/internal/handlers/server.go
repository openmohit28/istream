package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/config"
	"istream/backend/internal/store"
)

type Server struct {
	cfg     config.Config
	users   *store.Users
	results *store.Results
	resumes *store.Resumes
}

func NewServer(db *sql.DB, cfg config.Config) http.Handler {
	s := &Server{
		cfg:     cfg,
		users:   &store.Users{DB: db},
		results: &store.Results{DB: db},
		resumes: &store.Resumes{DB: db},
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

	quizRoutes := api.Group("/quiz", s.requireAuth())
	quizRoutes.GET("/questions", s.handleQuizQuestions)
	quizRoutes.POST("/submit", s.handleQuizSubmit)
	quizRoutes.GET("/results", s.handleQuizResults)
	quizRoutes.GET("/results/:id", s.handleQuizResult)

	jobRoutes := api.Group("/jobs", s.requireAuth())
	jobRoutes.GET("/search-url", s.handleJobSearchURL)

	resumeRoutes := api.Group("/resumes", s.requireAuth())
	resumeRoutes.POST("", s.handleResumeCreate)
	resumeRoutes.GET("", s.handleResumeList)
	resumeRoutes.GET("/:id", s.handleResumeGet)
	resumeRoutes.PUT("/:id", s.handleResumeUpdate)
	resumeRoutes.DELETE("/:id", s.handleResumeDelete)
	resumeRoutes.POST("/:id/keyword-check", s.handleResumeKeywordCheck)

	return r
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func errorJSON(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
