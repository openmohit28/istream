package resume

import (
	"regexp"
	"sort"
	"strings"
)

// KeywordReport compares a resume against a job description the way an ATS
// keyword filter would: which significant terms from the posting appear in
// the resume, which are missing, and a 0-100 coverage score.
type KeywordReport struct {
	Score   int      `json:"score"`
	Matched []string `json:"matched"`
	Missing []string `json:"missing"`
}

// maxKeywords caps how many job-description terms we score against, ranked
// by frequency in the posting.
const maxKeywords = 30

var (
	tokenPattern  = regexp.MustCompile(`[a-z0-9][a-z0-9+#./-]*`)
	numberPattern = regexp.MustCompile(`^[0-9.]+$`)
)

// stopwords covers common English plus boilerplate that appears in nearly
// every job posting and carries no matching signal.
var stopwords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "been": true, "but": true, "by": true, "can": true, "do": true,
	"for": true, "from": true, "has": true, "have": true, "in": true, "is": true,
	"it": true, "its": true, "of": true, "on": true, "or": true, "our": true,
	"that": true, "the": true, "their": true, "them": true, "they": true,
	"this": true, "to": true, "was": true, "we": true, "were": true, "will": true,
	"with": true, "you": true, "your": true, "not": true, "all": true, "if": true,
	"more": true, "other": true, "into": true, "than": true, "each": true,
	"about": true, "across": true, "also": true, "who": true, "what": true,
	"when": true, "where": true, "how": true, "per": true, "etc": true,
	// job-posting boilerplate
	"job": true, "role": true, "work": true, "working": true, "team": true,
	"teams": true, "years": true, "year": true, "experience": true,
	"experienced": true, "skills": true, "ability": true, "strong": true,
	"excellent": true, "required": true, "requirements": true, "preferred": true,
	"qualifications": true, "responsibilities": true, "candidate": true,
	"candidates": true, "including": true, "include": true, "includes": true,
	"looking": true, "join": true, "opportunity": true, "position": true,
	"company": true, "us": true, "plus": true, "bonus": true, "must": true,
	"should": true, "would": true, "well": true, "based": true, "using": true,
	"knowledge": true, "understanding": true, "related": true, "relevant": true,
	"day": true, "help": true, "new": true, "within": true, "environment": true,
}

func tokenize(text string) []string {
	raw := tokenPattern.FindAllString(strings.ToLower(text), -1)
	tokens := make([]string, 0, len(raw))
	for _, t := range raw {
		// Interior dots/slashes stay ("node.js", "ci/cd"); sentence-final
		// punctuation must not ("Docker." should match "docker").
		t = strings.TrimRight(t, "./-")
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func significant(token string) bool {
	if stopwords[token] {
		return false
	}
	// Keep short tech terms that carry signal (go, r, ai, ml, c#, qa...)
	// but drop bare numbers and single generic letters.
	if len(token) < 2 {
		return false
	}
	if numberPattern.MatchString(token) {
		return false
	}
	return true
}

// CheckKeywords extracts the most frequent significant terms from the job
// description and reports which ones the resume covers.
func CheckKeywords(doc Document, jobDescription string) KeywordReport {
	freq := map[string]int{}
	order := []string{}
	for _, tok := range tokenize(jobDescription) {
		if !significant(tok) {
			continue
		}
		if freq[tok] == 0 {
			order = append(order, tok)
		}
		freq[tok]++
	}

	// Rank by frequency, then by first appearance for stability.
	position := map[string]int{}
	for i, tok := range order {
		position[tok] = i
	}
	sort.SliceStable(order, func(i, j int) bool {
		if freq[order[i]] != freq[order[j]] {
			return freq[order[i]] > freq[order[j]]
		}
		return position[order[i]] < position[order[j]]
	})
	if len(order) > maxKeywords {
		order = order[:maxKeywords]
	}

	resumeTokens := map[string]bool{}
	for _, tok := range tokenize(doc.PlainText()) {
		resumeTokens[tok] = true
	}

	report := KeywordReport{Matched: []string{}, Missing: []string{}}
	for _, tok := range order {
		if resumeTokens[tok] {
			report.Matched = append(report.Matched, tok)
		} else {
			report.Missing = append(report.Missing, tok)
		}
	}
	if len(order) > 0 {
		report.Score = len(report.Matched) * 100 / len(order)
	}
	return report
}
