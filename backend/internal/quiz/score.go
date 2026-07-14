package quiz

import (
	"fmt"
	"math"
	"sort"
)

// Score converts a complete answer set (question ID -> 1..5) into 0-100
// per-dimension scores. Every question in the bank must be answered.
func Score(answers map[string]int) (Scores, error) {
	if len(answers) != len(Questions) {
		return nil, fmt.Errorf("expected %d answers, got %d", len(Questions), len(answers))
	}

	sums := map[Dimension]int{}
	counts := map[Dimension]int{}
	for id, value := range answers {
		q, ok := QuestionByID[id]
		if !ok {
			return nil, fmt.Errorf("unknown question id %q", id)
		}
		if value < AnswerMin || value > AnswerMax {
			return nil, fmt.Errorf("answer for %q must be between %d and %d", id, AnswerMin, AnswerMax)
		}
		sums[q.Dimension] += value
		counts[q.Dimension]++
	}

	scores := Scores{}
	for _, d := range Dimensions {
		n := counts[d]
		min, max := n*AnswerMin, n*AnswerMax
		scores[d] = int(math.Round(float64(sums[d]-min) / float64(max-min) * 100))
	}
	return scores, nil
}

// jobVector expands a Holland code like "IRE" into dimension weights:
// primary 3, secondary 2, tertiary 1.
func jobVector(hollandCode string) map[Dimension]float64 {
	v := map[Dimension]float64{}
	for i, ch := range hollandCode {
		v[Dimension(ch)] = float64(3 - i)
	}
	return v
}

// demandBoost nudges ranking toward growing fields without hiding fit.
func demandBoost(d DemandOutlook) float64 {
	switch d {
	case Growing:
		return 5
	case Declining:
		return -5
	default:
		return 0
	}
}

// Match ranks the catalog against the user's profile. Fit is the cosine
// similarity (0-100) between the score vector and the job's Holland vector;
// ordering additionally rewards growing demand.
func Match(scores Scores, limit int) []JobMatch {
	matches := make([]JobMatch, 0, len(Jobs))
	for _, job := range Jobs {
		fit := int(math.Round(cosine(scores, jobVector(job.HollandCode)) * 100))
		matches = append(matches, JobMatch{Job: job, Fit: fit})
	}

	sort.SliceStable(matches, func(i, j int) bool {
		ri := float64(matches[i].Fit) + demandBoost(matches[i].Demand)
		rj := float64(matches[j].Fit) + demandBoost(matches[j].Demand)
		if ri != rj {
			return ri > rj
		}
		return matches[i].Title < matches[j].Title
	})

	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	return matches
}

func cosine(scores Scores, job map[Dimension]float64) float64 {
	var dot, normUser, normJob float64
	for _, d := range Dimensions {
		u := float64(scores[d])
		j := job[d]
		dot += u * j
		normUser += u * u
		normJob += j * j
	}
	if normUser == 0 || normJob == 0 {
		return 0
	}
	return dot / (math.Sqrt(normUser) * math.Sqrt(normJob))
}
