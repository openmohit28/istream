package quiz

// Dimension is one of the six RIASEC career-interest dimensions
// (Holland Codes), the standard model for interest-to-occupation matching.
type Dimension string

const (
	Realistic     Dimension = "R"
	Investigative Dimension = "I"
	Artistic      Dimension = "A"
	Social        Dimension = "S"
	Enterprising  Dimension = "E"
	Conventional  Dimension = "C"
)

var Dimensions = []Dimension{Realistic, Investigative, Artistic, Social, Enterprising, Conventional}

var DimensionLabels = map[Dimension]string{
	Realistic:     "Realistic (hands-on)",
	Investigative: "Investigative (analytical)",
	Artistic:      "Artistic (creative)",
	Social:        "Social (helping)",
	Enterprising:  "Enterprising (leading)",
	Conventional:  "Conventional (organizing)",
}

type Question struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Dimension Dimension `json:"-"` // hidden from clients so answers stay unbiased
}

type DemandOutlook string

const (
	Growing   DemandOutlook = "growing"
	Stable    DemandOutlook = "stable"
	Declining DemandOutlook = "declining"
)

type AIRisk string

const (
	LowRisk    AIRisk = "low"
	MediumRisk AIRisk = "medium"
	HighRisk   AIRisk = "high"
)

// Job is a catalog entry. HollandCode is the job's dominant dimensions in
// priority order (e.g. "IRE"); demand and AI-risk come from 2026 labor
// research (WEF Future of Jobs, PwC AI Jobs Barometer).
type Job struct {
	Title       string        `json:"title"`
	HollandCode string        `json:"hollandCode"`
	Category    string        `json:"category"`
	Demand      DemandOutlook `json:"demand"`
	AIRisk      AIRisk        `json:"aiRisk"`
	Blurb       string        `json:"blurb"`
}

// Scores maps each dimension to a 0-100 normalized score.
type Scores map[Dimension]int

type JobMatch struct {
	Job
	Fit int `json:"fit"` // 0-100 cosine similarity with the user profile
}
