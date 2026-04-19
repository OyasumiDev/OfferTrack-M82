// core/cmd/offertrack/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/QUERTY/OfferTrack-M82/internal/cli"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// loadEnv carga el .env desde la raíz del proyecto.
// Busca hacia arriba desde el ejecutable hasta encontrarlo.
func loadEnv() {
	// Intentar rutas comunes
	candidates := []string{
		".env",
		"../../.env",
		"../../../.env",
	}

	// En desarrollo, usar la ruta del archivo fuente
	if _, file, _, ok := runtime.Caller(0); ok {
		root := filepath.Join(filepath.Dir(file), "..", "..", "..", ".env")
		candidates = append([]string{root}, candidates...)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
	}
	// Si no se encuentra, continuar sin .env (variables de entorno del sistema)
}
