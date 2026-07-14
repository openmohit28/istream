package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"istream/backend/internal/jobsearch"
)

// handleJobSearchURL builds a pre-filtered LinkedIn jobs-search deep link
// from the user's filters.
func (s *Server) handleJobSearchURL(c *gin.Context) {
	params := jobsearch.SearchParams{
		Keywords:     c.Query("keywords"),
		Location:     c.Query("location"),
		Workplace:    c.Query("workplace"),
		Experience:   c.Query("experience"),
		JobType:      c.Query("jobType"),
		PostedWithin: c.Query("postedWithin"),
	}
	url, err := jobsearch.BuildLinkedInURL(params)
	if err != nil {
		errorJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"provider": "linkedin", "url": url})
}
