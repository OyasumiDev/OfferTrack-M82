package domain

import "time"

type Job struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Description string    `json:"description"`
	Salary      string    `json:"salary"`
	Modality    string    `json:"modality"`
	Location    string    `json:"location"`
	Portal      string    `json:"portal"`
	URL         string    `json:"url"`
	ScrapedAt   time.Time `json:"scraped_at"`
	CompatScore int       `json:"compat_score"`
	Status      string    `json:"status"`
}
