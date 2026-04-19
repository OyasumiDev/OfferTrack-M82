// core/internal/scraper/occ_url_test.go
package scraper

import (
	"testing"
)

func TestOccSlugify(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Desarrollador .NET", "desarrollador-net"},
		{".NET Developer", "net-developer"},
		{"Nuevo León", "nuevo-leon"},
		{"San Pedro Garza García", "san-pedro-garza-garcia"},
		{"C# / ASP.NET", "c-asp-net"},
		{"Monterrey", "monterrey"},
	}

	for _, c := range cases {
		got := occSlugify(c.input)
		if got != c.expected {
			t.Errorf("occSlugify(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestBuildOccURL(t *testing.T) {
	cases := []struct {
		name     string
		params   SearchParams
		contains string
	}{
		{
			"puesto nacional",
			SearchParams{Keywords: "desarrollador .NET"},
			"/empleos/de-desarrollador-net/",
		},
		{
			"puesto + ciudad conocida",
			SearchParams{Keywords: "desarrollador .NET", City: "Monterrey"},
			"/en-nuevo-leon/en-la-ciudad-de-monterrey/",
		},
		{
			"puesto + estado + ciudad",
			SearchParams{Keywords: "desarrollador .NET", State: "Nuevo León", City: "Monterrey"},
			"/en-nuevo-leon/en-la-ciudad-de-monterrey/",
		},
		{
			"modalidad remoto",
			SearchParams{Keywords: "desarrollador .NET", Modality: ModalityRemoto},
			"tipo-home-office-remoto",
		},
		{
			"salario",
			SearchParams{Keywords: "desarrollador .NET", SalaryMin: 25000, SalaryMax: 70000},
			"?salary=25000,70000",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			url := buildOccURL(c.params)
			if !contains(url, c.contains) {
				t.Errorf("buildOccURL(%+v) = %q, expected to contain %q", c.params, url, c.contains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
