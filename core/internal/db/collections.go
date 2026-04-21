// core/internal/db/collections.go
package db

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

const (
	vectorSize = uint64(384)
	distance   = qdrant.Distance_Cosine
)

// CollectionNames agrupa los nombres de las colecciones leídos del config.
type CollectionNames struct {
	Jobs    string
	Profile string
	CVs     string
	Memory  string
}

// InitCollections crea las colecciones necesarias si no existen.
// Es idempotente: se puede llamar en cada arranque sin error.
func InitCollections(ctx context.Context, client *qdrant.Client, names CollectionNames) error {
	collections := []struct {
		name string
	}{
		{names.Jobs},
		{names.Profile},
		{names.CVs},
		{names.Memory},
	}

	for _, col := range collections {
		if col.name == "" {
			continue
		}
		if err := ensureCollection(ctx, client, col.name); err != nil {
			return fmt.Errorf("collections: error creando %q → %w", col.name, err)
		}
	}
	return nil
}

// ensureCollection crea la colección si no existe; no hace nada si ya existe.
func ensureCollection(ctx context.Context, client *qdrant.Client, name string) error {
	exists, err := client.CollectionExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: distance,
			OnDisk:   qdrant.PtrOf(true),
		}),
	})
	if err != nil {
		return err
	}

	fmt.Printf("[qdrant] Colección %q creada (384 dims, Cosine)\n", name)
	return nil
}
