package db

import "github.com/qdrant/go-client/qdrant"

type QdrantClient struct {
	client *qdrant.Client
}

func NewQdrantClient(host string, port int) (*QdrantClient, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		return nil, err
	}
	return &QdrantClient{client: client}, nil
}
