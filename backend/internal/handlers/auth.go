package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/auth"
	"istream/backend/internal/models"
	"istream/backend/internal/store"
)

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (s *Server) handleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)

	if _, err := mail.ParseAddress(email); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid email address")
		return
	}
	if name == "" {
		errorJSON(c, http.StatusBadRequest, "name is required")
		return
	}
	if len(req.Password) < 8 {
		errorJSON(c, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not process password")
		return
	}

	user, err := s.users.Create(c.Request.Context(), email, name, hash)
	if errors.Is(err, store.ErrDuplicateEmail) {
		errorJSON(c, http.StatusConflict, "email already registered")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not create user")
		return
	}

	s.respondWithToken(c, http.StatusCreated, user)
}

func (s *Server) handleLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.users.ByEmail(c.Request.Context(), email)
	if errors.Is(err, store.ErrNotFound) || (err == nil && !auth.CheckPassword(user.PasswordHash, req.Password)) {
		errorJSON(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not log in")
		return
	}

	s.respondWithToken(c, http.StatusOK, user)
}

func (s *Server) handleMe(c *gin.Context) {
	user, err := s.users.ByID(c.Request.Context(), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusUnauthorized, "user no longer exists")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load user")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *Server) respondWithToken(c *gin.Context, status int, user models.User) {
	token, err := auth.SignToken(s.cfg.JWTSecret, user.ID)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not sign token")
		return
	}
	c.JSON(status, authResponse{Token: token, User: user})
}