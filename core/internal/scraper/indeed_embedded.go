package scraper

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IndeedEmbeddedJob representa la estructura del JSON embebido en la página de resultados de Indeed.
type IndeedEmbeddedJob struct {
	JobKey            string   `json:"jobkey"`
	DisplayTitle      string   `json:"displayTitle"`
	NormTitle         string   `json:"normTitle"`
	Company           string   `json:"company"`
	FormattedLocation string   `json:"formattedLocation"`
	JobLocationCity   string   `json:"jobLocationCity"`
	JobLocationState  string   `json:"jobLocationState"`
	CreateDate        int64    `json:"createDate"`
	PubDate           int64    `json:"pubDate"`
	RemoteLocation    bool     `json:"remoteLocation"`
	JobTypes          []string `json:"jobTypes"`
	SalarySnippet     struct {
		Currency string `json:"currency"`
		Text     string `json:"text"`
		Source   string `json:"source"`
	} `json:"salarySnippet"`
	Snippet   string `json:"snippet"`
	Sponsored bool   `json:"sponsored"`
	Expired   bool   `json:"expired"`
}

var resultsRegex = regexp.MustCompile(`(?s)"results"\s*:\s*(\[\s*\{.*?\}\s*\])`)

// ExtractEmbeddedResults extrae el array de vacantes del JSON embebido en el HTML de Indeed.
func ExtractEmbeddedResults(html string) ([]IndeedEmbeddedJob, error) {
	matches := resultsRegex.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		var jobs []IndeedEmbeddedJob
		if err := json.Unmarshal([]byte(m[1]), &jobs); err != nil {
			continue
		}
		if len(jobs) > 0 && jobs[0].JobKey != "" {
			return jobs, nil
		}
	}
	return nil, errors.New("indeed: no se encontró array results en HTML")
}

// ToListing convierte un IndeedEmbeddedJob al contrato Listing del core.
func (j IndeedEmbeddedJob) ToListing() Listing {
	l := Listing{
		ExternalID:  j.JobKey,
		Title:       j.DisplayTitle,
		Company:     j.Company,
		Location:    j.FormattedLocation,
		Description: j.Snippet,
		SourceURL:   BuildIndeedDetailURL(j.JobKey),
		FetchedAt:   time.Now(),
		Portal:      PortalIndeed,
	}
	if j.RemoteLocation {
		l.Modality = ModalityRemoto
	}
	l.SalaryMin, l.SalaryMax, l.SalaryRaw = parseIndeedSalary(j.SalarySnippet.Text)
	if j.PubDate > 0 {
		l.PostedAt = time.UnixMilli(j.PubDate)
	}
	return l
}

var salaryNumRegex = regexp.MustCompile(`\$([\d,]+(?:\.\d+)?)`)

func parseIndeedSalary(text string) (min, max int, raw string) {
	if text == "" {
		return 0, 0, ""
	}
	raw = text
	matches := salaryNumRegex.FindAllStringSubmatch(text, -1)
	var nums []int
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		s := strings.ReplaceAll(m[1], ",", "")
		if idx := strings.Index(s, "."); idx >= 0 {
			s = s[:idx]
		}
		n, err := strconv.Atoi(s)
		if err == nil {
			nums = append(nums, n)
		}
	}
	lower := strings.ToLower(text)
	switch {
	case len(nums) == 0:
		return 0, 0, raw
	case strings.Contains(lower, "hasta") && len(nums) == 1:
		return 0, nums[0], raw
	case len(nums) >= 2:
		return nums[0], nums[1], raw
	default:
		return nums[0], nums[0], raw
	}
}
