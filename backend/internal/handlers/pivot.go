package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/pivot"
	"istream/backend/internal/store"
)

// threadView is a thread plus its computed position in the tree.
func threadView(t store.PivotThread) (gin.H, error) {
	var steps []pivot.Step
	if err := json.Unmarshal(t.Steps, &steps); err != nil {
		return nil, err
	}
	state, err := pivot.WalkPath(steps)
	if err != nil {
		return nil, err
	}
	view := gin.H{
		"id":        t.ID,
		"steps":     steps,
		"createdAt": t.CreatedAt,
		"updatedAt": t.UpdatedAt,
	}
	if t.ForkedFrom != nil {
		view["forkedFrom"] = *t.ForkedFrom
	}
	if state.Current != nil {
		view["current"] = state.Current
	}
	if state.Outcome != nil {
		view["outcome"] = state.Outcome
	}
	return view, nil
}

func (s *Server) respondThread(c *gin.Context, status int, t store.PivotThread) {
	view, err := threadView(t)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not read thread state")
		return
	}
	c.JSON(status, view)
}

type threadStepsRequest struct {
	Steps []pivot.Step `json:"steps"`
}

// validateSteps walks the submitted path and returns its JSON encoding.
func validateSteps(c *gin.Context, steps []pivot.Step) ([]byte, bool) {
	if steps == nil {
		steps = []pivot.Step{}
	}
	if _, err := pivot.WalkPath(steps); err != nil {
		errorJSON(c, http.StatusBadRequest, err.Error())
		return nil, false
	}
	data, err := json.Marshal(steps)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not encode steps")
		return nil, false
	}
	return data, true
}

func (s *Server) handlePivotThreadCreate(c *gin.Context) {
	var req threadStepsRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			errorJSON(c, http.StatusBadRequest, "invalid JSON body")
			return
		}
	}
	data, ok := validateSteps(c, req.Steps)
	if !ok {
		return
	}
	t, err := s.pivots.Create(c.Request.Context(), c.GetString(userIDKey), data, nil)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not create thread")
		return
	}
	s.respondThread(c, http.StatusCreated, t)
}

func (s *Server) handlePivotThreadUpdate(c *gin.Context) {
	var req threadStepsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	data, ok := validateSteps(c, req.Steps)
	if !ok {
		return
	}
	t, err := s.pivots.UpdateSteps(c.Request.Context(), c.Param("id"), c.GetString(userIDKey), data)
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "thread not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not update thread")
		return
	}
	s.respondThread(c, http.StatusOK, t)
}

type forkRequest struct {
	// AtStep is how many leading steps to keep (0 = restart from the root).
	AtStep int `json:"atStep"`
}

// handlePivotThreadFork copies the first AtStep answers of an existing
// thread into a NEW thread, so the original exploration stays intact.
func (s *Server) handlePivotThreadFork(c *gin.Context) {
	var req forkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorJSON(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	userID := c.GetString(userIDKey)
	src, err := s.pivots.ByIDForUser(c.Request.Context(), c.Param("id"), userID)
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "thread not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load thread")
		return
	}

	var steps []pivot.Step
	if err := json.Unmarshal(src.Steps, &steps); err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not read thread state")
		return
	}
	if req.AtStep < 0 || req.AtStep > len(steps) {
		errorJSON(c, http.StatusBadRequest, "atStep out of range")
		return
	}

	data, ok := validateSteps(c, steps[:req.AtStep])
	if !ok {
		return
	}
	fork, err := s.pivots.Create(c.Request.Context(), userID, data, &src.ID)
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not fork thread")
		return
	}
	s.respondThread(c, http.StatusCreated, fork)
}

func (s *Server) handlePivotThreadList(c *gin.Context) {
	threads, err := s.pivots.ListByUser(c.Request.Context(), c.GetString(userIDKey))
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load threads")
		return
	}
	views := make([]gin.H, 0, len(threads))
	for _, t := range threads {
		view, err := threadView(t)
		if err != nil {
			errorJSON(c, http.StatusInternalServerError, "could not read thread state")
			return
		}
		views = append(views, view)
	}
	c.JSON(http.StatusOK, gin.H{"threads": views})
}

func (s *Server) handlePivotThreadGet(c *gin.Context) {
	t, err := s.pivots.ByIDForUser(c.Request.Context(), c.Param("id"), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "thread not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not load thread")
		return
	}
	s.respondThread(c, http.StatusOK, t)
}

func (s *Server) handlePivotThreadDelete(c *gin.Context) {
	err := s.pivots.Delete(c.Request.Context(), c.Param("id"), c.GetString(userIDKey))
	if errors.Is(err, store.ErrNotFound) {
		errorJSON(c, http.StatusNotFound, "thread not found")
		return
	}
	if err != nil {
		errorJSON(c, http.StatusInternalServerError, "could not delete thread")
		return
	}
	c.Status(http.StatusNoContent)
}
