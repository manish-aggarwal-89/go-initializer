package migrations

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexHelper provides utility functions for creating and managing MongoDB indices
type IndexHelper struct {
	db *mongo.Database
}

// NewIndexHelper creates a new index helper instance
func NewIndexHelper(db *mongo.Database) *IndexHelper {
	return &IndexHelper{db: db}
}

// CreateIndex creates a single index on specified fields with existence check
func (h *IndexHelper) CreateIndex(ctx context.Context, collectionName string, keys bson.D, opts ...*options.IndexOptions) error {
	collection := h.db.Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys: keys,
	}

	if len(opts) > 0 {
		indexModel.Options = opts[0]
	}

	// Check if index already exists by name (if name is provided)
	var indexName string
	if len(opts) > 0 && opts[0].Name != nil {
		indexName = *opts[0].Name
		exists, err := h.IndexExists(ctx, collectionName, indexName)
		if err != nil {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		if exists {
			fmt.Printf("Index '%s' already exists on collection '%s', skipping creation\n", indexName, collectionName)
			return nil
		}
	}

	createdIndexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		errMsg := err.Error()
		if mongo.IsDuplicateKeyError(err) ||
			strings.Contains(errMsg, "already exists") ||
			strings.Contains(errMsg, "dup key") ||
			strings.Contains(errMsg, "duplicate key") {
			fmt.Printf("Index creation skipped on collection '%s' - index already exists or duplicate data found\n", collectionName)
			return nil
		}
		return fmt.Errorf("failed to create index on %s: %w", collectionName, err)
	}

	fmt.Printf("Created index '%s' on collection '%s'\n", createdIndexName, collectionName)
	return nil
}

// CreateUniqueIndex creates a unique index on specified fields
func (h *IndexHelper) CreateUniqueIndex(ctx context.Context, collectionName string, keys bson.D, name string) error {
	opts := options.Index().SetUnique(true)
	if name != "" {
		opts.SetName(name)
	}
	return h.CreateIndex(ctx, collectionName, keys, opts)
}

// CreateCompoundIndex creates a compound index on multiple fields
func (h *IndexHelper) CreateCompoundIndex(ctx context.Context, collectionName string, keys bson.D, name string, unique bool) error {
	opts := options.Index()
	if name != "" {
		opts.SetName(name)
	}
	if unique {
		opts.SetUnique(true)
	}
	return h.CreateIndex(ctx, collectionName, keys, opts)
}

// DropIndex drops an index by name
func (h *IndexHelper) DropIndex(ctx context.Context, collectionName string, indexName string) error {
	collection := h.db.Collection(collectionName)
	_, err := collection.Indexes().DropOne(ctx, indexName)
	if err != nil {
		return fmt.Errorf("failed to drop index '%s' on collection '%s': %w", indexName, collectionName, err)
	}
	fmt.Printf("Dropped index '%s' from collection '%s'\n", indexName, collectionName)
	return nil
}

// IndexExists checks if an index with the given name exists
func (h *IndexHelper) IndexExists(ctx context.Context, collectionName string, indexName string) (bool, error) {
	indexes, err := h.ListIndexes(ctx, collectionName)
	if err != nil {
		return false, err
	}
	for _, idx := range indexes {
		if name, ok := idx["name"].(string); ok && name == indexName {
			return true, nil
		}
	}
	return false, nil
}

// ListIndexes returns all indexes for a collection
func (h *IndexHelper) ListIndexes(ctx context.Context, collectionName string) ([]bson.M, error) {
	collection := h.db.Collection(collectionName)
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes for collection '%s': %w", collectionName, err)
	}
	defer cursor.Close(ctx)
	var indexes []bson.M
	if err := cursor.All(ctx, &indexes); err != nil {
		return nil, fmt.Errorf("failed to decode indexes for collection '%s': %w", collectionName, err)
	}
	return indexes, nil
}

// CreateSimpleIndex creates an ascending index on a single field
func (h *IndexHelper) CreateSimpleIndex(ctx context.Context, collectionName, field, indexName string) error {
	if indexName != "" {
		exists, err := h.IndexExists(ctx, collectionName, indexName)
		if err != nil {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		if exists {
			fmt.Printf("Index '%s' already exists on collection '%s', skipping creation\n", indexName, collectionName)
			return nil
		}
	}
	keys := bson.D{{Key: field, Value: 1}}
	opts := options.Index()
	if indexName != "" {
		opts.SetName(indexName)
	}
	return h.CreateIndex(ctx, collectionName, keys, opts)
}

// CreateSimpleUniqueIndex creates a unique ascending index on a single field
func (h *IndexHelper) CreateSimpleUniqueIndex(ctx context.Context, collectionName, field, indexName string) error {
	if indexName != "" {
		exists, err := h.IndexExists(ctx, collectionName, indexName)
		if err != nil {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		if exists {
			fmt.Printf("Index '%s' already exists on collection '%s', skipping creation\n", indexName, collectionName)
			return nil
		}
	}
	keys := bson.D{{Key: field, Value: 1}}
	return h.CreateUniqueIndex(ctx, collectionName, keys, indexName)
}

// CreateMultiFieldIndex creates a compound index on multiple fields (all ascending)
func (h *IndexHelper) CreateMultiFieldIndex(ctx context.Context, collectionName string, fields []string, indexName string, unique bool) error {
	if indexName != "" {
		exists, err := h.IndexExists(ctx, collectionName, indexName)
		if err != nil {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		if exists {
			fmt.Printf("Index '%s' already exists on collection '%s', skipping creation\n", indexName, collectionName)
			return nil
		}
	}
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: 1})
	}
	if unique {
		return h.CreateUniqueIndex(ctx, collectionName, keys, indexName)
	}
	opts := options.Index()
	if indexName != "" {
		opts.SetName(indexName)
	}
	return h.CreateIndex(ctx, collectionName, keys, opts)
}
