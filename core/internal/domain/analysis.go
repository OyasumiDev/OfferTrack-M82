package domain

type Analysis struct {
	JobID              string   `json:"job_id"`
	CompatibilityScore int      `json:"compatibility_score"`
	Strengths          []string `json:"strengths"`
	Gaps               []string `json:"gaps"`
	SalaryEstimate     string   `json:"salary_estimate"`
	Recommendation     string   `json:"recommendation"`
	RawAnalysis        string   `json:"raw_analysis"`
}
