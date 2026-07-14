package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func sampleResume() map[string]any {
	return map[string]any{
		"targetTitle": "Backend Engineer",
		"contact": map[string]any{
			"fullName": "Mohit Rawat",
			"email":    "mohit@example.com",
			"location": "Bengaluru",
		},
		"summary": "Backend engineer building APIs in Go and PostgreSQL.",
		"experience": []map[string]any{{
			"company": "Acme",
			"title":   "Software Engineer",
			"bullets": []string{"Built REST APIs with Go and Gin"},
		}},
		"skills": []string{"Go", "PostgreSQL", "Docker"},
	}
}

func TestJobSearchURLRequiresAuth(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "GET", "/api/jobs/search-url?keywords=nurse", nil, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}

func TestJobSearchURLBuildsLink(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "jobs@example.com")

	rec := doJSON(t, srv, "GET",
		"/api/jobs/search-url?keywords=machine+learning+engineer&location=Remote&workplace=remote&experience=entry&postedWithin=week",
		nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Provider string `json:"provider"`
		URL      string `json:"url"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Provider != "linkedin" {
		t.Errorf("provider: got %q", resp.Provider)
	}
	for _, want := range []string{"linkedin.com/jobs/search", "f_WT=2", "f_E=2", "f_TPR=r604800"} {
		if !containsStr(resp.URL, want) {
			t.Errorf("url missing %q: %s", want, resp.URL)
		}
	}
}

func containsStr(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func TestJobSearchURLRejectsBadFilter(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "jobs@example.com")

	rec := doJSON(t, srv, "GET", "/api/jobs/search-url?keywords=x&workplace=moon", nil, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestResumeCRUD(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "resume@example.com")

	// Create
	created := doJSON(t, srv, "POST", "/api/resumes", sampleResume(), token)
	if created.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d: %s", created.Code, created.Body.String())
	}
	var resumeResp struct {
		ID    string          `json:"id"`
		Title string          `json:"title"`
		Data  json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &resumeResp); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	if resumeResp.Title != "Backend Engineer" {
		t.Errorf("title derived from targetTitle: got %q", resumeResp.Title)
	}

	// List
	list := doJSON(t, srv, "GET", "/api/resumes", nil, token)
	var listResp struct {
		Resumes []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"resumes"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResp.Resumes) != 1 || listResp.Resumes[0].ID != resumeResp.ID {
		t.Fatalf("list should contain the created resume, got %+v", listResp.Resumes)
	}

	// Update
	updated := sampleResume()
	updated["targetTitle"] = "Platform Engineer"
	upd := doJSON(t, srv, "PUT", "/api/resumes/"+resumeResp.ID, updated, token)
	if upd.Code != http.StatusOK {
		t.Fatalf("update: want 200, got %d: %s", upd.Code, upd.Body.String())
	}
	var updResp struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(upd.Body.Bytes(), &updResp); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	if updResp.Title != "Platform Engineer" {
		t.Errorf("update title: got %q", updResp.Title)
	}

	// Get
	get := doJSON(t, srv, "GET", "/api/resumes/"+resumeResp.ID, nil, token)
	if get.Code != http.StatusOK {
		t.Fatalf("get: want 200, got %d", get.Code)
	}

	// Delete
	del := doJSON(t, srv, "DELETE", "/api/resumes/"+resumeResp.ID, nil, token)
	if del.Code != http.StatusNoContent {
		t.Fatalf("delete: want 204, got %d", del.Code)
	}
	gone := doJSON(t, srv, "GET", "/api/resumes/"+resumeResp.ID, nil, token)
	if gone.Code != http.StatusNotFound {
		t.Fatalf("after delete: want 404, got %d", gone.Code)
	}
}

func TestResumeValidation(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "resume@example.com")

	noTitle := sampleResume()
	noTitle["targetTitle"] = ""
	rec := doJSON(t, srv, "POST", "/api/resumes", noTitle, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for missing title, got %d", rec.Code)
	}

	noEmail := sampleResume()
	noEmail["contact"] = map[string]any{"fullName": "Mohit"}
	rec = doJSON(t, srv, "POST", "/api/resumes", noEmail, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for missing email, got %d", rec.Code)
	}
}

func TestResumeUserScoping(t *testing.T) {
	srv := newTestServer(t)
	ownerToken := registerUser(t, srv, "owner-r@example.com")
	otherToken := registerUser(t, srv, "other-r@example.com")

	created := doJSON(t, srv, "POST", "/api/resumes", sampleResume(), ownerToken)
	var resumeResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &resumeResp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	for _, tc := range []struct {
		method, path string
		body         any
	}{
		{"GET", "/api/resumes/" + resumeResp.ID, nil},
		{"PUT", "/api/resumes/" + resumeResp.ID, sampleResume()},
		{"DELETE", "/api/resumes/" + resumeResp.ID, nil},
		{"POST", "/api/resumes/" + resumeResp.ID + "/keyword-check", map[string]string{"jobDescription": "Go"}},
	} {
		rec := doJSON(t, srv, tc.method, tc.path, tc.body, otherToken)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s %s as other user: want 404, got %d", tc.method, tc.path, rec.Code)
		}
	}
}

func TestResumeKeywordCheck(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "resume@example.com")

	created := doJSON(t, srv, "POST", "/api/resumes", sampleResume(), token)
	var resumeResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &resumeResp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	jd := "Backend Engineer with Go, PostgreSQL, Kubernetes. Kubernetes required."
	rec := doJSON(t, srv, "POST", fmt.Sprintf("/api/resumes/%s/keyword-check", resumeResp.ID),
		map[string]string{"jobDescription": jd}, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var report struct {
		Score   int      `json:"score"`
		Matched []string `json:"matched"`
		Missing []string `json:"missing"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &report); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if report.Score <= 0 || report.Score >= 100 {
		t.Errorf("expected partial score, got %d", report.Score)
	}
	if len(report.Missing) == 0 {
		t.Error("kubernetes should be missing")
	}

	// Empty job description rejected.
	rec = doJSON(t, srv, "POST", fmt.Sprintf("/api/resumes/%s/keyword-check", resumeResp.ID),
		map[string]string{"jobDescription": ""}, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for empty JD, got %d", rec.Code)
	}
}
