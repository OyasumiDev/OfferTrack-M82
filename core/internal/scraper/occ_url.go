// core/internal/scraper/occ_url.go
// Construcción de URL canónica para OCC Mundial — sin estado, sin I/O
package scraper

import (
	"fmt"
	"regexp"
	"strings"
)

// occSlugify aplica el algoritmo del Prompt Maestro §2.3:
// lowercase → quitar tildes → "." = espacio → [^a-z0-9]+→"-" → colapsar → trim
func occSlugify(s string) string {
	s = strings.ToLower(s)
	s = removeAccentsOCC(s)
	s = strings.ReplaceAll(s, ".", " ")
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// removeAccentsOCC reemplaza vocales acentuadas y ñ del español.
func removeAccentsOCC(s string) string {
	r := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"ü", "u", "ñ", "n",
		"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u",
		"Ü", "u", "Ñ", "n",
	)
	return r.Replace(s)
}

// occCityToState mapea ciudad→estado (slugs ya normalizados).
var occCityToState = map[string]string{
	"monterrey": "nuevo-leon", "san-nicolas-de-los-garza": "nuevo-leon",
	"guadalupe": "nuevo-leon", "apodaca": "nuevo-leon", "san-pedro-garza-garcia": "nuevo-leon",
	"guadalajara": "jalisco", "zapopan": "jalisco", "tlaquepaque": "jalisco", "tonala": "jalisco",
	"mexico": "ciudad-de-mexico", "cdmx": "ciudad-de-mexico", "ciudad-de-mexico": "ciudad-de-mexico",
	"puebla": "puebla",
	"tijuana": "baja-california", "mexicali": "baja-california", "ensenada": "baja-california",
	"leon": "guanajuato", "irapuato": "guanajuato", "celaya": "guanajuato",
	"queretaro": "queretaro", "san-juan-del-rio": "queretaro",
	"merida": "yucatan",
	"cancun": "quintana-roo", "playa-del-carmen": "quintana-roo", "chetumal": "quintana-roo",
	"hermosillo": "sonora", "ciudad-obregon": "sonora",
	"chihuahua": "chihuahua", "ciudad-juarez": "chihuahua",
	"culiacan": "sinaloa", "mazatlan": "sinaloa",
	"torreon": "coahuila", "saltillo": "coahuila", "monclova": "coahuila",
	"san-luis-potosi": "san-luis-potosi",
	"aguascalientes": "aguascalientes",
	"morelia": "michoacan", "uruapan": "michoacan",
	"toluca": "estado-de-mexico", "ecatepec": "estado-de-mexico", "naucalpan": "estado-de-mexico",
	"veracruz": "veracruz", "xalapa": "veracruz",
	"acapulco": "guerrero",
	"oaxaca": "oaxaca",
	"villahermosa": "tabasco",
	"tuxtla-gutierrez": "chiapas",
	"durango": "durango",
	"zacatecas": "zacatecas",
	"tepic": "nayarit",
	"colima": "colima",
	"campeche": "campeche",
	"la-paz": "baja-california-sur", "los-cabos": "baja-california-sur",
	"pachuca": "hidalgo",
	"cuernavaca": "morelos",
	"tlaxcala": "tlaxcala",
}

const occBase = "https://www.occ.com.mx"

// buildOccURL construye la URL canónica de OCC per Prompt Maestro §2.2.
// SearchParams.City puede ser "Monterrey", "monterrey", o vacío.
// SearchParams.State puede ser "Nuevo León", "nuevo-leon", o vacío.
func buildOccURL(p SearchParams) string {
	roleSlug := occSlugify(p.Keywords)
	url := fmt.Sprintf("%s/empleos/de-%s/", occBase, roleSlug)

	// Determinar segmentos de estado y ciudad
	citySlug := occSlugify(p.City)
	stateSlug := occSlugify(p.State)

	if stateSlug != "" && citySlug != "" {
		url += fmt.Sprintf("en-%s/en-la-ciudad-de-%s/", stateSlug, citySlug)
	} else if stateSlug != "" {
		url += fmt.Sprintf("en-%s/", stateSlug)
	} else if citySlug != "" {
		// Autodetección ciudad → estado
		if knownState, ok := occCityToState[citySlug]; ok {
			url += fmt.Sprintf("en-%s/en-la-ciudad-de-%s/", knownState, citySlug)
		} else {
			url += fmt.Sprintf("en-%s/", citySlug)
		}
	}

	// Segmento de modalidad
	if p.Modality == ModalityRemoto {
		url += "tipo-home-office-remoto"
	}

	// Query string de salario
	if p.SalaryMin > 0 && p.SalaryMax > 0 {
		url += fmt.Sprintf("?salary=%d,%d", p.SalaryMin, p.SalaryMax)
	} else if p.SalaryMin > 0 {
		url += fmt.Sprintf("?salary=%d,", p.SalaryMin)
	}

	return url
}
