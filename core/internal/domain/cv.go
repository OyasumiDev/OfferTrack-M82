package domain

import "time"

type CV struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	BaseContent string    `json:"base_content"`
	Adapted     string    `json:"adapted"`
	CreatedAt   time.Time `json:"created_at"`
}
