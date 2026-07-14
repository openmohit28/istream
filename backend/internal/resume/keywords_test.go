package resume

import (
	"slices"
	"strings"
	"testing"
)

func sampleDoc() Document {
	return Document{
		TargetTitle: "Backend Engineer",
		Contact:     Contact{FullName: "Mohit", Email: "m@example.com"},
		Summary:     "Backend engineer building APIs in Go and PostgreSQL.",
		Experience: []Experience{{
			Company: "Acme",
			Title:   "Software Engineer",
			Bullets: []string{"Built REST APIs with Go and Gin", "Cut query latency 40% on PostgreSQL"},
		}},
		Skills: []string{"Go", "PostgreSQL", "Docker"},
	}
}

func TestValidate(t *testing.T) {
	doc := sampleDoc()
	if err := doc.Validate(); err != nil {
		t.Fatalf("valid doc rejected: %v", err)
	}

	missingTitle := sampleDoc()
	missingTitle.TargetTitle = " "
	if err := missingTitle.Validate(); err == nil {
		t.Error("expected error for missing target title")
	}

	missingName := sampleDoc()
	missingName.Contact.FullName = ""
	if err := missingName.Validate(); err == nil {
		t.Error("expected error for missing name")
	}

	missingEmail := sampleDoc()
	missingEmail.Contact.Email = ""
	if err := missingEmail.Validate(); err == nil {
		t.Error("expected error for missing email")
	}
}

func TestCheckKeywordsMatchesAndMisses(t *testing.T) {
	jd := `We need a Backend Engineer with Go, PostgreSQL and Kubernetes experience.
	Kubernetes is essential. Kafka is a plus. Strong communication skills required.`

	report := CheckKeywords(sampleDoc(), jd)

	for _, want := range []string{"go", "postgresql", "backend", "engineer"} {
		if !slices.Contains(report.Matched, want) {
			t.Errorf("expected %q in matched, got %v", want, report.Matched)
		}
	}
	for _, want := range []string{"kubernetes", "kafka"} {
		if !slices.Contains(report.Missing, want) {
			t.Errorf("expected %q in missing, got %v", want, report.Missing)
		}
	}
	if report.Score <= 0 || report.Score >= 100 {
		t.Errorf("expected partial score, got %d", report.Score)
	}
}

func TestCheckKeywordsFiltersBoilerplate(t *testing.T) {
	jd := "The candidate must have strong experience and excellent skills for this role."
	report := CheckKeywords(sampleDoc(), jd)
	all := append(append([]string{}, report.Matched...), report.Missing...)
	for _, boiler := range []string{"candidate", "experience", "skills", "role", "strong", "excellent"} {
		if slices.Contains(all, boiler) {
			t.Errorf("boilerplate %q should be filtered, got %v", boiler, all)
		}
	}
}

func TestCheckKeywordsRanksByFrequency(t *testing.T) {
	jd := "Python Python Python Terraform Terraform Ansible"
	report := CheckKeywords(Document{}, jd)
	if len(report.Missing) != 3 {
		t.Fatalf("want 3 missing keywords, got %v", report.Missing)
	}
	if report.Missing[0] != "python" || report.Missing[1] != "terraform" || report.Missing[2] != "ansible" {
		t.Errorf("wrong frequency order: %v", report.Missing)
	}
	if report.Score != 0 {
		t.Errorf("empty resume should score 0, got %d", report.Score)
	}
}

func TestCheckKeywordsEmptyJobDescription(t *testing.T) {
	report := CheckKeywords(sampleDoc(), "")
	if report.Score != 0 || len(report.Matched) != 0 || len(report.Missing) != 0 {
		t.Errorf("empty JD should produce empty report, got %+v", report)
	}
}

func TestCheckKeywordsFullCoverage(t *testing.T) {
	report := CheckKeywords(sampleDoc(), "Go PostgreSQL Docker")
	if report.Score != 100 {
		t.Errorf("want 100, got %d (missing: %v)", report.Score, report.Missing)
	}
}

func TestCheckKeywordsStripsSentencePunctuation(t *testing.T) {
	// Words at sentence end ("Docker.") must match their clean form, while
	// interior dots survive ("node.js").
	report := CheckKeywords(sampleDoc(), "We use Docker. Also node.js daily.")
	if !slices.Contains(report.Matched, "docker") {
		t.Errorf("docker should match despite trailing period, got matched=%v missing=%v",
			report.Matched, report.Missing)
	}
	if !slices.Contains(report.Missing, "node.js") {
		t.Errorf("node.js should keep its interior dot, got %v", report.Missing)
	}
}

func TestPlainTextIncludesAllSections(t *testing.T) {
	doc := sampleDoc()
	doc.Education = []Education{{School: "IIT", Degree: "B.Tech", Field: "Computer Science"}}
	doc.Certifications = []string{"AWS Solutions Architect"}
	text := doc.PlainText()
	for _, want := range []string{"Backend Engineer", "Acme", "REST APIs", "IIT", "Computer Science", "AWS Solutions Architect", "Docker"} {
		if !strings.Contains(text, want) {
			t.Errorf("plain text missing %q", want)
		}
	}
}
