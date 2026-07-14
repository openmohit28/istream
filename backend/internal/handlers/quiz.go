package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/quiz"
	"istream/backend/internal/store"
)

const matchLimit = 8

type submitRequest struct {
	Answers map[string]int `json:"answers"`
}

func (s *Server) handleQuizQuestions(c *gin.Context) {
	questions := make([]gin.H, 0, len(quiz.Questions))
	for _, q := range quiz.Questions {
		questions = append(questions, gin.H{"id": q.ID, "text": q.Text})
	}
	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
		"scale": gin.H{
			"min":      quiz.AnswerMin,
			"max":      quiz.AnswerMax,
			"minLabel": "Strongly disagree",
			"maxLabel": "Strongly agree",
		},
	})
}

func (s *Server) handleQuizSubmit(c *gin.Context) {
	var req submitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	scores, err := quiz.Score(req.Answers)
	if err != nil {
		errorJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	matches := quiz.Match(scores, matchLimit)

	answersJSON, err := json.Marshal(req.Answers)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not encode answers")
		return
	}
	scoresJSON, err := json.Marshal(scores)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not encode scores")
		return
	}
	matchesJSON, err := json.Marshal(matches)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not encode matches")
		return
	}

	result, err := s.results.Create(c.Request.Context(), c.GetString(userIDKey), answersJSON, scoresJSON, matchesJSON)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not save result")
		return
	}

	c.JSON(http.StatusCreated, resultDetail(result))
}

func (s *Server) handleQuizResults(c *gin.Context) {
	results, err := s.results.ListByUser(c.Request.Context(), c.GetString(userIDKey))
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load results")
		return
	}

	summaries := make([]gin.H, 0, len(results))
	for _, r := range results {
		summary := gin.H{"id": r.ID, "createdAt": r.CreatedAt, "scores": r.Scores}
		var matches []quiz.JobMatch
		if err := json.Unmarshal(r.Matches, &matches); err == nil && len(matches) > 0 {
			summary["topMatch"] = gin.H{"title": matches[0].Title, "fit": matches[0].Fit}
		}
		summaries = append(summaries, summary)
	}
	c.JSON(http.StatusOK, gin.H{"results": summaries})
}

func (s *Server) handleQuizResult(c *gin.Context) {
	result, err := s.results.ByIDForUser(c.Request.Context(), c.Param("id"), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "result not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load result")
		return
	}
	c.JSON(http.StatusOK, resultDetail(result))
}

func resultDetail(r store.TestResult) gin.H {
	return gin.H{
		"id":        r.ID,
		"createdAt": r.CreatedAt,
		"scores":    r.Scores,
		"matches":   r.Matches,
	}
}
