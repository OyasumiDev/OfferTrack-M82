package domain

import "time"

type Job struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Company         string    `json:"company"`
	Description     string    `json:"description"`
	SalaryMin       int       `json:"salary_min"`
	SalaryMax       int       `json:"salary_max"`
	Currency        string    `json:"currency"`
	Modality        string    `json:"modality"` // "remote" | "hybrid" | "onsite"
	Location        string    `json:"location"`
	Portal          string    `json:"portal"` // "occ" | "computrabajo" | "indeed"
	URL             string    `json:"url"`
	PostedAt        time.Time `json:"posted_at"`
	ScrapedAt       time.Time `json:"scraped_at"`
	CompatScore     int       `json:"compat_score"`
	EmbeddingVector []float32 `json:"embedding_vector,omitempty"`
}
