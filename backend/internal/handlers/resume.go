package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/resume"
	"istream/backend/internal/store"
)

// decodeResumeDoc parses and validates the resume document from the body.
func decodeResumeDoc(c *gin.Context) (resume.Document, []byte, bool) {
	var doc resume.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return doc, nil, false
	}
	if err := doc.Validate(); err != nil {
		errorJSON(c, http.StatusBadRequest, err.Error())
		return doc, nil, false
	}
	data, err := json.Marshal(doc)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not encode resume")
		return doc, nil, false
	}
	return doc, data, true
}

func (s *Server) handleResumeCreate(c *gin.Context) {
	doc, data, ok := decodeResumeDoc(c)
	if !ok {
		return
	}
	r, err := s.resumes.Create(c.Request.Context(), c.GetString(userIDKey), doc.TargetTitle, data)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not save resume")
		return
	}
	c.JSON(http.StatusCreated, r)
}

func (s *Server) handleResumeUpdate(c *gin.Context) {
	doc, data, ok := decodeResumeDoc(c)
	if !ok {
		return
	}
	r, err := s.resumes.Update(c.Request.Context(), c.Param("id"), c.GetString(userIDKey), doc.TargetTitle, data)
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "resume not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not update resume")
		return
	}
	c.JSON(http.StatusOK, r)
}

func (s *Server) handleResumeList(c *gin.Context) {
	resumes, err := s.resumes.ListByUser(c.Request.Context(), c.GetString(userIDKey))
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load resumes")
		return
	}
	// Keep list responses light: id, title, timestamps only.
	summaries := make([]gin.H, 0, len(resumes))
	for _, r := range resumes {
		summaries = append(summaries, gin.H{
			"id":        r.ID,
			"title":     r.Title,
			"createdAt": r.CreatedAt,
			"updatedAt": r.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"resumes": summaries})
}

func (s *Server) handleResumeGet(c *gin.Context) {
	r, err := s.resumes.ByIDForUser(c.Request.Context(), c.Param("id"), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "resume not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load resume")
		return
	}
	c.JSON(http.StatusOK, r)
}

func (s *Server) handleResumeDelete(c *gin.Context) {
	err := s.resumes.Delete(c.Request.Context(), c.Param("id"), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "resume not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not delete resume")
		return
	}
	c.Status(http.StatusNoContent)
}

type keywordCheckRequest struct {
	JobDescription string `json:"jobDescription"`
}

// handleResumeKeywordCheck scores a stored resume against a pasted job
// description the way an ATS keyword filter would.
func (s *Server) handleResumeKeywordCheck(c *gin.Context) {
	var req keywordCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.JobDescription == "" {
		errorJSON(c, http.StatusBadRequest, "jobDescription is required")
		return
	}

	r, err := s.resumes.ByIDForUser(c.Request.Context(), c.Param("id"), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "resume not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load resume")
		return
	}

	var doc resume.Document
	if err := json.Unmarshal(r.Data, &doc); err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not decode resume")
		return
	}
	c.JSON(http.StatusOK, resume.CheckKeywords(doc, req.JobDescription))
}
