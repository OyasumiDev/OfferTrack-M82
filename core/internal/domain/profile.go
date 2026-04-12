package domain

type Profile struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Skills     []string `json:"skills"`
	Experience string   `json:"experience"`
	Education  string   `json:"education"`
	SalaryMin  int      `json:"salary_min"`
	Modalities []string `json:"modalities"`
	Locations  []string `json:"locations"`
}
