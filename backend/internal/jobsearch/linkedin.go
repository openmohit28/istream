// Package jobsearch builds pre-filtered job-board search URLs. LinkedIn has
// no public job-search API and scraping violates its terms, so we deep-link
// users into LinkedIn's own search with their filters applied.
package jobsearch

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type SearchParams struct {
	Keywords     string // required
	Location     string
	Workplace    string // onsite | remote | hybrid
	Experience   string // internship | entry | associate | mid-senior | director | executive
	JobType      string // fulltime | parttime | contract | temporary | internship
	PostedWithin string // day | week | month
}

var workplaceCodes = map[string]string{
	"onsite": "1",
	"remote": "2",
	"hybrid": "3",
}

var experienceCodes = map[string]string{
	"internship": "1",
	"entry":      "2",
	"associate":  "3",
	"mid-senior": "4",
	"director":   "5",
	"executive":  "6",
}

var jobTypeCodes = map[string]string{
	"fulltime":   "F",
	"parttime":   "P",
	"contract":   "C",
	"temporary":  "T",
	"internship": "I",
}

var postedWithinCodes = map[string]string{
	"day":   "r86400",
	"week":  "r604800",
	"month": "r2592000",
}

// BuildLinkedInURL maps the filters onto LinkedIn's jobs-search query
// parameters (f_WT workplace, f_E experience, f_JT job type, f_TPR recency).
func BuildLinkedInURL(p SearchParams) (string, error) {
	keywords := strings.TrimSpace(p.Keywords)
	if keywords == "" {
		return "", errors.New("keywords are required")
	}

	q := url.Values{}
	q.Set("keywords", keywords)

	if loc := strings.TrimSpace(p.Location); loc != "" {
		q.Set("location", loc)
	}
	if err := setCoded(q, "f_WT", "workplace", p.Workplace, workplaceCodes); err != nil {
		return "", err
	}
	if err := setCoded(q, "f_E", "experience", p.Experience, experienceCodes); err != nil {
		return "", err
	}
	if err := setCoded(q, "f_JT", "jobType", p.JobType, jobTypeCodes); err != nil {
		return "", err
	}
	if err := setCoded(q, "f_TPR", "postedWithin", p.PostedWithin, postedWithinCodes); err != nil {
		return "", err
	}

	return "https://www.linkedin.com/jobs/search/?" + q.Encode(), nil
}

func setCoded(q url.Values, param, name, value string, codes map[string]string) error {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return nil
	}
	code, ok := codes[value]
	if !ok {
		return fmt.Errorf("invalid %s %q", name, value)
	}
	q.Set(param, code)
	return nil
}
