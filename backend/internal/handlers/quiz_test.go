package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"istream/backend/internal/quiz"
)

// registerUser creates a fresh user and returns their auth token.
func registerUser(t *testing.T, srv http.Handler, email string) string {
	t.Helper()
	rec := doJSON(t, srv, "POST", "/api/auth/register",
		map[string]string{"email": email, "name": "Quiz Tester", "password": "supersecret1"}, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	token, _ := decodeAuth(t, rec)
	return token
}

// investigativeAnswers builds a complete answer set skewed toward
// Investigative and Realistic dimensions.
func investigativeAnswers() map[string]int {
	answers := map[string]int{}
	for _, q := range quiz.Questions {
		switch q.Dimension {
		case quiz.Investigative:
			answers[q.ID] = 5
		case quiz.Realistic:
			answers[q.ID] = 4
		default:
			answers[q.ID] = 2
		}
	}
	return answers
}

type resultPayload struct {
	ID        string          `json:"id"`
	CreatedAt string          `json:"createdAt"`
	Scores    map[string]int  `json:"scores"`
	Matches   []quiz.JobMatch `json:"matches"`
}

func TestQuizQuestionsRequiresAuth(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "GET", "/api/quiz/questions", nil, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}

func TestQuizQuestionsReturnsBank(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "quiz@example.com")

	rec := doJSON(t, srv, "GET", "/api/quiz/questions", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Questions []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"questions"`
		Scale struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"scale"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Questions) != len(quiz.Questions) {
		t.Errorf("want %d questions, got %d", len(quiz.Questions), len(resp.Questions))
	}
	if resp.Scale.Min != 1 || resp.Scale.Max != 5 {
		t.Errorf("unexpected scale: %+v", resp.Scale)
	}
	// The dimension must not leak to clients.
	if body := rec.Body.String(); json.Valid([]byte(body)) && containsDimensionKey(body) {
		t.Error("questions response must not expose dimensions")
	}
}

func containsDimensionKey(body string) bool {
	var raw struct {
		Questions []map[string]any `json:"questions"`
	}
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		return false
	}
	for _, q := range raw.Questions {
		if _, ok := q["dimension"]; ok {
			return true
		}
	}
	return false
}

func TestQuizSubmitComputesScoresAndMatches(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "quiz@example.com")

	rec := doJSON(t, srv, "POST", "/api/quiz/submit",
		map[string]any{"answers": investigativeAnswers()}, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var result resultPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.ID == "" {
		t.Error("expected a result id")
	}
	if result.Scores["I"] != 100 {
		t.Errorf("Investigative score: want 100, got %d", result.Scores["I"])
	}
	if len(result.Matches) == 0 {
		t.Fatal("expected job matches")
	}
	top := result.Matches[0]
	if top.HollandCode[0] != 'I' && top.HollandCode[0] != 'R' {
		t.Errorf("top match %q (%s) should lead with I or R", top.Title, top.HollandCode)
	}
}

func TestQuizSubmitRejectsIncompleteAnswers(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "quiz@example.com")

	answers := investigativeAnswers()
	delete(answers, "r1")
	rec := doJSON(t, srv, "POST", "/api/quiz/submit", map[string]any{"answers": answers}, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestQuizSubmitRejectsOutOfRangeValue(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "quiz@example.com")

	answers := investigativeAnswers()
	answers["r1"] = 9
	rec := doJSON(t, srv, "POST", "/api/quiz/submit", map[string]any{"answers": answers}, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestQuizResultsListAndDetail(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "quiz@example.com")

	sub := doJSON(t, srv, "POST", "/api/quiz/submit", map[string]any{"answers": investigativeAnswers()}, token)
	var created resultPayload
	if err := json.Unmarshal(sub.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode submit: %v", err)
	}

	list := doJSON(t, srv, "GET", "/api/quiz/results", nil, token)
	if list.Code != http.StatusOK {
		t.Fatalf("list: want 200, got %d", list.Code)
	}
	var listResp struct {
		Results []struct {
			ID       string `json:"id"`
			TopMatch struct {
				Title string `json:"title"`
				Fit   int    `json:"fit"`
			} `json:"topMatch"`
		} `json:"results"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResp.Results) != 1 || listResp.Results[0].ID != created.ID {
		t.Fatalf("expected the submitted result in list, got %+v", listResp.Results)
	}
	if listResp.Results[0].TopMatch.Title == "" {
		t.Error("list summary should include topMatch")
	}

	detail := doJSON(t, srv, "GET", fmt.Sprintf("/api/quiz/results/%s", created.ID), nil, token)
	if detail.Code != http.StatusOK {
		t.Fatalf("detail: want 200, got %d: %s", detail.Code, detail.Body.String())
	}
}

func TestQuizResultsAreUserScoped(t *testing.T) {
	srv := newTestServer(t)
	ownerToken := registerUser(t, srv, "owner@example.com")
	otherToken := registerUser(t, srv, "other@example.com")

	sub := doJSON(t, srv, "POST", "/api/quiz/submit", map[string]any{"answers": investigativeAnswers()}, ownerToken)
	var created resultPayload
	if err := json.Unmarshal(sub.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode submit: %v", err)
	}

	// The other user cannot fetch the owner's result...
	detail := doJSON(t, srv, "GET", fmt.Sprintf("/api/quiz/results/%s", created.ID), nil, otherToken)
	if detail.Code != http.StatusNotFound {
		t.Fatalf("cross-user detail: want 404, got %d", detail.Code)
	}

	// ...and their own list is empty.
	list := doJSON(t, srv, "GET", "/api/quiz/results", nil, otherToken)
	var listResp struct {
		Results []any `json:"results"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResp.Results) != 0 {
		t.Fatalf("other user's list should be empty, got %d items", len(listResp.Results))
	}
}
