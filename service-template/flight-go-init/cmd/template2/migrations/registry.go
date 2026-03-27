package migrations

import (
	"context"
	"sort"

	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	ID   string
	Name string
	Up   func(ctx context.Context, db *mongo.Database) error
	Down func(ctx context.Context, db *mongo.Database) error
}

var registry []Migration

func Register(m Migration) {
	registry = append(registry, m)
}

func All() []Migration {
	out := make([]Migration, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func FindByID(id string) *Migration {
	for i := range registry {
		if registry[i].ID == id {
			return &registry[i]
		}
	}
	return nil
}
