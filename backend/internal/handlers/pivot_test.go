package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"istream/backend/internal/pivot"
)

type threadResponse struct {
	ID         string         `json:"id"`
	Steps      []pivot.Step   `json:"steps"`
	ForkedFrom string         `json:"forkedFrom"`
	Current    *pivot.Node    `json:"current"`
	Outcome    *pivot.Outcome `json:"outcome"`
}

func decodeThread(t *testing.T, body []byte) threadResponse {
	t.Helper()
	var resp threadResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode thread: %v (%s)", err, body)
	}
	return resp
}

// burnoutPath walks driver -> hours-fix -> employer-flex -> reduce-hours.
var burnoutPath = []pivot.Step{
	{NodeID: "driver", Option: "I'm burnt out - I need more time and energy for life"},
	{NodeID: "hours-fix", Option: "Yes - the job is fine, it's the hours"},
	{NodeID: "employer-flex", Option: "Open - there are precedents for part-time or 4-day weeks"},
}

func TestPivotThreadCreateStartsAtRoot(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "pivot@example.com")

	rec := doJSON(t, srv, "POST", "/api/pivot/threads", nil, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	thread := decodeThread(t, rec.Body.Bytes())
	if thread.Current == nil || thread.Current.ID != pivot.RootID {
		t.Fatalf("new thread should sit at root, got %+v", thread.Current)
	}
	if thread.Outcome != nil {
		t.Error("new thread should have no outcome")
	}
	if len(thread.Current.Options) != 4 {
		t.Errorf("root should offer 4 drivers, got %d", len(thread.Current.Options))
	}
}

func TestPivotThreadAnswerToOutcome(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "pivot@example.com")

	created := decodeThread(t, doJSON(t, srv, "POST", "/api/pivot/threads", nil, token).Body.Bytes())

	// Answer step by step, as the UI does.
	for i := 1; i <= len(burnoutPath); i++ {
		rec := doJSON(t, srv, "PUT", "/api/pivot/threads/"+created.ID,
			map[string]any{"steps": burnoutPath[:i]}, token)
		if rec.Code != http.StatusOK {
			t.Fatalf("step %d: want 200, got %d: %s", i, rec.Code, rec.Body.String())
		}
		thread := decodeThread(t, rec.Body.Bytes())
		if i < len(burnoutPath) && thread.Current == nil {
			t.Fatalf("step %d: expected a next question", i)
		}
		if i == len(burnoutPath) {
			if thread.Outcome == nil || thread.Outcome.ID != "reduce-hours" {
				t.Fatalf("want reduce-hours outcome, got %+v", thread.Outcome)
			}
			if len(thread.Outcome.Plan) == 0 || len(thread.Outcome.Resources) == 0 {
				t.Error("outcome must carry plan and resources")
			}
		}
	}
}

func TestPivotThreadRejectsInvalidPath(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "pivot@example.com")

	created := decodeThread(t, doJSON(t, srv, "POST", "/api/pivot/threads", nil, token).Body.Bytes())

	rec := doJSON(t, srv, "PUT", "/api/pivot/threads/"+created.ID,
		map[string]any{"steps": []pivot.Step{{NodeID: "driver", Option: "not a real option"}}}, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPivotThreadFork(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "pivot@example.com")

	created := decodeThread(t, doJSON(t, srv, "POST", "/api/pivot/threads", nil, token).Body.Bytes())
	doJSON(t, srv, "PUT", "/api/pivot/threads/"+created.ID, map[string]any{"steps": burnoutPath}, token)

	// Fork keeping only the first answer.
	rec := doJSON(t, srv, "POST", "/api/pivot/threads/"+created.ID+"/fork",
		map[string]any{"atStep": 1}, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("fork: want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	fork := decodeThread(t, rec.Body.Bytes())
	if fork.ID == created.ID {
		t.Fatal("fork must be a new thread")
	}
	if fork.ForkedFrom != created.ID {
		t.Errorf("fork should record its origin, got %q", fork.ForkedFrom)
	}
	if len(fork.Steps) != 1 || fork.Steps[0].NodeID != "driver" {
		t.Fatalf("fork should keep 1 step, got %+v", fork.Steps)
	}
	if fork.Current == nil || fork.Current.ID != "hours-fix" {
		t.Errorf("fork should sit at hours-fix, got %+v", fork.Current)
	}

	// The original thread is untouched, still at its outcome.
	orig := decodeThread(t, doJSON(t, srv, "GET", "/api/pivot/threads/"+created.ID, nil, token).Body.Bytes())
	if orig.Outcome == nil || orig.Outcome.ID != "reduce-hours" {
		t.Errorf("original thread must keep its outcome, got %+v", orig.Outcome)
	}
	if len(orig.Steps) != 3 {
		t.Errorf("original thread must keep all steps, got %d", len(orig.Steps))
	}

	// Both threads appear in the list.
	var list struct {
		Threads []threadResponse `json:"threads"`
	}
	if err := json.Unmarshal(doJSON(t, srv, "GET", "/api/pivot/threads", nil, token).Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list.Threads) != 2 {
		t.Fatalf("want 2 threads after fork, got %d", len(list.Threads))
	}
}

func TestPivotThreadForkAtStepBounds(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "pivot@example.com")

	created := decodeThread(t, doJSON(t, srv, "POST", "/api/pivot/threads", nil, token).Body.Bytes())
	doJSON(t, srv, "PUT", "/api/pivot/threads/"+created.ID, map[string]any{"steps": burnoutPath}, token)

	for _, bad := range []int{-1, 4} {
		rec := doJSON(t, srv, "POST", "/api/pivot/threads/"+created.ID+"/fork",
			map[string]any{"atStep": bad}, token)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("atStep=%d: want 400, got %d", bad, rec.Code)
		}
	}
}

func TestPivotThreadDelete(t *testing.T) {
	srv := newTestServer(t)
	token := registerUser(t, srv, "pivot@example.com")

	created := decodeThread(t, doJSON(t, srv, "POST", "/api/pivot/threads", nil, token).Body.Bytes())
	rec := doJSON(t, srv, "DELETE", "/api/pivot/threads/"+created.ID, nil, token)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rec.Code)
	}
	gone := doJSON(t, srv, "GET", "/api/pivot/threads/"+created.ID, nil, token)
	if gone.Code != http.StatusNotFound {
		t.Fatalf("after delete: want 404, got %d", gone.Code)
	}
}

func TestPivotThreadsAreUserScoped(t *testing.T) {
	srv := newTestServer(t)
	ownerToken := registerUser(t, srv, "owner-p@example.com")
	otherToken := registerUser(t, srv, "other-p@example.com")

	created := decodeThread(t, doJSON(t, srv, "POST", "/api/pivot/threads", nil, ownerToken).Body.Bytes())

	for _, tc := range []struct {
		method, path string
		body         any
	}{
		{"GET", "/api/pivot/threads/" + created.ID, nil},
		{"PUT", "/api/pivot/threads/" + created.ID, map[string]any{"steps": burnoutPath[:1]}},
		{"POST", "/api/pivot/threads/" + created.ID + "/fork", map[string]any{"atStep": 0}},
		{"DELETE", "/api/pivot/threads/" + created.ID, nil},
	} {
		rec := doJSON(t, srv, tc.method, tc.path, tc.body, otherToken)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s %s as other user: want 404, got %d", tc.method, tc.path, rec.Code)
		}
	}
}
