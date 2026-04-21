package scraper

import (
	"os"
	"testing"
)

func TestBuildIndeedURL(t *testing.T) {
	cases := []struct {
		name string
		q    IndeedQuery
		want string
	}{
		{
			name: "solo puesto",
			q:    IndeedQuery{Puesto: "desarrollador .NET"},
			want: "https://mx.indeed.com/jobs?q=desarrollador+.NET",
		},
		{
			name: "puesto+ciudad+estado",
			q:    IndeedQuery{Puesto: "desarrollador .NET", Ciudad: "Monterrey", Estado: "Nuevo León"},
			want: "https://mx.indeed.com/jobs?l=Monterrey%2C+Nuevo+Le%C3%B3n&q=desarrollador+.NET",
		},
		{
			name: "remoto",
			q:    IndeedQuery{Puesto: "desarrollador .NET", Modality: "remote"},
			want: "https://mx.indeed.com/jobs?l=remote&q=desarrollador+.NET",
		},
		{
			name: "salario 300k",
			q:    IndeedQuery{Puesto: "dev", SalaryMinAnnualMXN: 300000},
			want: "https://mx.indeed.com/jobs?q=dev&salaryType=%24300%2C000",
		},
		{
			name: "paginación p3",
			q:    IndeedQuery{Puesto: "dev", Page: 3},
			want: "https://mx.indeed.com/jobs?q=dev&start=20",
		},
		{
			name: "híbrido+FT",
			q:    IndeedQuery{Puesto: "dev", Modality: "hybrid", FullTime: true},
			want: "https://mx.indeed.com/jobs?q=dev&sc=0kf%3Aattr%28PAXZC%29attr%28CF3CP%29%3B",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildIndeedURL(tc.q)
			if got != tc.want {
				t.Errorf("\nwant: %s\n got: %s", tc.want, got)
			}
		})
	}
}

func TestParseIndeedSalary(t *testing.T) {
	cases := []struct {
		input    string
		wantMin  int
		wantMax  int
	}{
		{"$20,000 a $26,000 por mes", 20000, 26000},
		{"Hasta $30,000 por mes", 0, 30000},
		{"$70,000 a $72,000 por mes - Tiempo completo", 70000, 72000},
		{"$15,000.00 por mes", 15000, 15000},
		{"", 0, 0},
		{"A convenir", 0, 0},
	}

	for _, tc := range cases {
		min, max, _ := parseIndeedSalary(tc.input)
		if min != tc.wantMin || max != tc.wantMax {
			t.Errorf("parseIndeedSalary(%q) = (%d, %d), want (%d, %d)",
				tc.input, min, max, tc.wantMin, tc.wantMax)
		}
	}
}

func TestExtractEmbeddedResults(t *testing.T) {
	data, err := os.ReadFile("testdata/indeed_page.html")
	if err != nil {
		t.Fatalf("fixture: %v", err)
	}
	jobs, err := ExtractEmbeddedResults(string(data))
	if err != nil {
		t.Fatalf("ExtractEmbeddedResults: %v", err)
	}
	if len(jobs) < 2 {
		t.Fatalf("want >= 2 jobs, got %d", len(jobs))
	}
	if jobs[0].JobKey == "" {
		t.Error("jobs[0].JobKey is empty")
	}
}

func TestToListing_RemoteLocation(t *testing.T) {
	j := IndeedEmbeddedJob{RemoteLocation: true}
	if got := j.ToListing().Modality; got != ModalityRemoto {
		t.Errorf("want %q, got %q", ModalityRemoto, got)
	}
}

func TestToListing_PubDateUnixMs(t *testing.T) {
	j := IndeedEmbeddedJob{PubDate: 1713600000000}
	want := int64(1713600000)
	if got := j.ToListing().PostedAt.Unix(); got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}
