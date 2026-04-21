package scraper

import (
	"fmt"
	"net/url"
)

const indeedBaseURL = "https://mx.indeed.com/jobs"

// IndeedQuery agrupa todos los parámetros de búsqueda de Indeed México.
type IndeedQuery struct {
	Puesto             string
	Ciudad             string
	Estado             string
	Page               int
	RadiusKm           int
	SalaryMinAnnualMXN int
	FromAgeDays        int
	Modality           string // "remote", "hybrid", ""
	FullTime           bool
	SortByDate         bool
}

// BuildIndeedURL construye la URL de búsqueda de Indeed México a partir de los parámetros dados.
func BuildIndeedURL(q IndeedQuery) string {
	v := url.Values{}
	v.Set("q", q.Puesto)

	if q.Modality == "remote" {
		v.Set("l", "remote")
	} else if q.Ciudad != "" && q.Estado != "" {
		v.Set("l", q.Ciudad+", "+q.Estado)
	} else if q.Ciudad != "" {
		v.Set("l", q.Ciudad)
	} else if q.Estado != "" {
		v.Set("l", q.Estado)
	}

	if q.RadiusKm > 0 && q.RadiusKm != 25 {
		v.Set("radius", fmt.Sprintf("%d", q.RadiusKm))
	}

	if q.SalaryMinAnnualMXN > 0 {
		v.Set("salaryType", "$"+formatWithCommas(q.SalaryMinAnnualMXN))
	}

	if q.FromAgeDays > 0 {
		v.Set("fromage", fmt.Sprintf("%d", q.FromAgeDays))
	}

	if q.Page > 1 {
		v.Set("start", fmt.Sprintf("%d", (q.Page-1)*10))
	}

	if q.SortByDate {
		v.Set("sort", "date")
	}

	var attrs []string
	if q.Modality == "hybrid" {
		attrs = append(attrs, "PAXZC")
	}
	if q.FullTime {
		attrs = append(attrs, "CF3CP")
	}
	if len(attrs) > 0 {
		sc := "0kf:"
		for _, a := range attrs {
			sc += "attr(" + a + ")"
		}
		sc += ";"
		v.Set("sc", sc)
	}

	return indeedBaseURL + "?" + v.Encode()
}

// BuildIndeedDetailURL devuelve la URL canónica de detalle para un jobkey dado.
func BuildIndeedDetailURL(jobkey string) string {
	v := url.Values{}
	v.Set("jk", jobkey)
	return "https://mx.indeed.com/viewjob?" + v.Encode()
}

func formatWithCommas(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	result := make([]byte, 0, len(s)+len(s)/3)
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
