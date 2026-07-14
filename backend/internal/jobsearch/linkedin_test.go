package jobsearch

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildLinkedInURLFull(t *testing.T) {
	got, err := BuildLinkedInURL(SearchParams{
		Keywords:     "machine learning engineer",
		Location:     "Bengaluru, India",
		Workplace:    "remote",
		Experience:   "mid-senior",
		JobType:      "fulltime",
		PostedWithin: "week",
	})
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "www.linkedin.com" || u.Path != "/jobs/search/" {
		t.Errorf("unexpected base: %s", got)
	}
	q := u.Query()
	checks := map[string]string{
		"keywords": "machine learning engineer",
		"location": "Bengaluru, India",
		"f_WT":     "2",
		"f_E":      "4",
		"f_JT":     "F",
		"f_TPR":    "r604800",
	}
	for key, want := range checks {
		if q.Get(key) != want {
			t.Errorf("%s: want %q, got %q", key, want, q.Get(key))
		}
	}
}

func TestBuildLinkedInURLMinimal(t *testing.T) {
	got, err := BuildLinkedInURL(SearchParams{Keywords: "nurse"})
	if err != nil {
		t.Fatal(err)
	}
	u, _ := url.Parse(got)
	q := u.Query()
	if q.Get("keywords") != "nurse" {
		t.Errorf("keywords: got %q", q.Get("keywords"))
	}
	for _, param := range []string{"location", "f_WT", "f_E", "f_JT", "f_TPR"} {
		if q.Has(param) {
			t.Errorf("param %s should be absent, got %q", param, q.Get(param))
		}
	}
}

func TestBuildLinkedInURLRequiresKeywords(t *testing.T) {
	if _, err := BuildLinkedInURL(SearchParams{Keywords: "   "}); err == nil {
		t.Fatal("expected error for empty keywords")
	}
}

func TestBuildLinkedInURLRejectsInvalidEnums(t *testing.T) {
	cases := []SearchParams{
		{Keywords: "x", Workplace: "moon"},
		{Keywords: "x", Experience: "wizard"},
		{Keywords: "x", JobType: "gig"},
		{Keywords: "x", PostedWithin: "decade"},
	}
	for _, p := range cases {
		if _, err := BuildLinkedInURL(p); err == nil {
			t.Errorf("expected error for %+v", p)
		}
	}
}

func TestBuildLinkedInURLNormalizesCase(t *testing.T) {
	got, err := BuildLinkedInURL(SearchParams{Keywords: "x", Workplace: "Remote", PostedWithin: "MONTH"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "f_WT=2") || !strings.Contains(got, "f_TPR=r2592000") {
		t.Errorf("case normalization failed: %s", got)
	}
}
