package quiz

import (
	"strings"
	"testing"
)

// allAnswers builds a complete answer set with the given default value,
// then applies overrides.
func allAnswers(defaultValue int, overrides map[string]int) map[string]int {
	answers := make(map[string]int, len(Questions))
	for _, q := range Questions {
		answers[q.ID] = defaultValue
	}
	for id, v := range overrides {
		answers[id] = v
	}
	return answers
}

func TestScoreAllMin(t *testing.T) {
	scores, err := Score(allAnswers(1, nil))
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range Dimensions {
		if scores[d] != 0 {
			t.Errorf("dimension %s: want 0, got %d", d, scores[d])
		}
	}
}

func TestScoreAllMax(t *testing.T) {
	scores, err := Score(allAnswers(5, nil))
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range Dimensions {
		if scores[d] != 100 {
			t.Errorf("dimension %s: want 100, got %d", d, scores[d])
		}
	}
}

func TestScoreSkewedProfile(t *testing.T) {
	// Max out Investigative, minimum everywhere else.
	scores, err := Score(allAnswers(1, map[string]int{"i1": 5, "i2": 5, "i3": 5, "i4": 5}))
	if err != nil {
		t.Fatal(err)
	}
	if scores[Investigative] != 100 {
		t.Errorf("Investigative: want 100, got %d", scores[Investigative])
	}
	if scores[Realistic] != 0 {
		t.Errorf("Realistic: want 0, got %d", scores[Realistic])
	}
}

func TestScoreRejectsIncomplete(t *testing.T) {
	answers := allAnswers(3, nil)
	delete(answers, "r1")
	if _, err := Score(answers); err == nil {
		t.Fatal("expected error for incomplete answers")
	}
}

func TestScoreRejectsUnknownQuestion(t *testing.T) {
	answers := allAnswers(3, nil)
	delete(answers, "r1")
	answers["bogus"] = 3
	if _, err := Score(answers); err == nil {
		t.Fatal("expected error for unknown question id")
	}
}

func TestScoreRejectsOutOfRange(t *testing.T) {
	for _, bad := range []int{0, 6, -1} {
		if _, err := Score(allAnswers(3, map[string]int{"r1": bad})); err == nil {
			t.Fatalf("expected error for answer value %d", bad)
		}
	}
}

func TestMatchRanksMatchingProfileFirst(t *testing.T) {
	// A strongly Investigative+Realistic profile should surface technical
	// jobs (IR* Holland codes) at the top.
	scores := Scores{Realistic: 80, Investigative: 100, Artistic: 10, Social: 10, Enterprising: 20, Conventional: 30}
	matches := Match(scores, 5)
	if len(matches) != 5 {
		t.Fatalf("want 5 matches, got %d", len(matches))
	}
	top := matches[0]
	if !strings.HasPrefix(top.HollandCode, "IR") && !strings.HasPrefix(top.HollandCode, "RI") {
		t.Errorf("top match %q (%s) does not lead with I/R dimensions", top.Title, top.HollandCode)
	}
	for _, m := range matches {
		if m.Fit < 0 || m.Fit > 100 {
			t.Errorf("fit out of range for %q: %d", m.Title, m.Fit)
		}
	}
}

func TestMatchSocialProfileSurfacesCareRoles(t *testing.T) {
	scores := Scores{Realistic: 10, Investigative: 20, Artistic: 20, Social: 100, Enterprising: 30, Conventional: 20}
	matches := Match(scores, 5)
	top := matches[0]
	if top.HollandCode[0] != 'S' {
		t.Errorf("top match %q (%s) is not Social-led", top.Title, top.HollandCode)
	}
}

func TestMatchZeroLimitReturnsAll(t *testing.T) {
	scores := Scores{Realistic: 50, Investigative: 50, Artistic: 50, Social: 50, Enterprising: 50, Conventional: 50}
	if got := len(Match(scores, 0)); got != len(Jobs) {
		t.Errorf("want %d matches, got %d", len(Jobs), got)
	}
}

func TestCatalogIntegrity(t *testing.T) {
	valid := map[byte]bool{'R': true, 'I': true, 'A': true, 'S': true, 'E': true, 'C': true}
	seen := map[string]bool{}
	for _, job := range Jobs {
		if seen[job.Title] {
			t.Errorf("duplicate job title %q", job.Title)
		}
		seen[job.Title] = true
		if len(job.HollandCode) != 3 {
			t.Errorf("%q: holland code %q must be 3 letters", job.Title, job.HollandCode)
		}
		chars := map[byte]bool{}
		for i := 0; i < len(job.HollandCode); i++ {
			ch := job.HollandCode[i]
			if !valid[ch] {
				t.Errorf("%q: invalid holland dimension %q", job.Title, string(ch))
			}
			if chars[ch] {
				t.Errorf("%q: repeated dimension in %q", job.Title, job.HollandCode)
			}
			chars[ch] = true
		}
		if job.Blurb == "" || job.Category == "" {
			t.Errorf("%q: missing blurb or category", job.Title)
		}
	}
}

func TestQuestionBankBalanced(t *testing.T) {
	counts := map[Dimension]int{}
	for _, q := range Questions {
		counts[q.Dimension]++
	}
	for _, d := range Dimensions {
		if counts[d] != 4 {
			t.Errorf("dimension %s has %d questions, want 4", d, counts[d])
		}
	}
}
