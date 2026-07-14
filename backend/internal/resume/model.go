// Package resume defines the structured resume document and the ATS keyword
// matcher. Resumes are stored as JSONB per user; rendering stays in the
// frontend (single-column, ATS-safe layout).
package resume

import (
	"errors"
	"strings"
)

type Contact struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	LinkedIn string `json:"linkedin"`
	Website  string `json:"website"`
}

type Experience struct {
	Company   string   `json:"company"`
	Title     string   `json:"title"`
	Location  string   `json:"location"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Current   bool     `json:"current"`
	Bullets   []string `json:"bullets"`
}

type Education struct {
	School   string `json:"school"`
	Degree   string `json:"degree"`
	Field    string `json:"field"`
	GradYear string `json:"gradYear"`
}

type Document struct {
	TargetTitle    string       `json:"targetTitle"`
	JobDescription string       `json:"jobDescription"`
	Contact        Contact      `json:"contact"`
	Summary        string       `json:"summary"`
	Experience     []Experience `json:"experience"`
	Education      []Education  `json:"education"`
	Skills         []string     `json:"skills"`
	Certifications []string     `json:"certifications"`
}

// Validate enforces the minimum needed to render a useful resume.
func (d Document) Validate() error {
	if strings.TrimSpace(d.TargetTitle) == "" {
		return errors.New("targetTitle is required")
	}
	if strings.TrimSpace(d.Contact.FullName) == "" {
		return errors.New("contact.fullName is required")
	}
	if strings.TrimSpace(d.Contact.Email) == "" {
		return errors.New("contact.email is required")
	}
	return nil
}

// PlainText flattens the searchable content of the resume for keyword
// matching - the same text an ATS would parse.
func (d Document) PlainText() string {
	var b strings.Builder
	write := func(parts ...string) {
		for _, p := range parts {
			if p = strings.TrimSpace(p); p != "" {
				b.WriteString(p)
				b.WriteString("\n")
			}
		}
	}
	write(d.TargetTitle, d.Summary)
	for _, e := range d.Experience {
		write(e.Title, e.Company)
		write(e.Bullets...)
	}
	for _, e := range d.Education {
		write(e.Degree, e.Field, e.School)
	}
	write(d.Skills...)
	write(d.Certifications...)
	return b.String()
}
