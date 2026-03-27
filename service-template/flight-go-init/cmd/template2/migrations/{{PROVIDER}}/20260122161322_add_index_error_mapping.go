package {{PROVIDER}}

import (
	"context"

	indexmanager "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/helper/index_manager"
	migrationregistry "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/registry"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	const collectionName = "integrator_error_mapping"
	const indexName = "is_deleted_asc"

	migrationregistry.MigrationRecordRegister(migrationregistry.MigrationRecord{
		ID:   "20260122161322",
		Name: "integrator_error_mapping_is_deleted_index",
		Up: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)

			fieldOrder, err := indexmanager.NewFieldOrderRef("isDeleted", indexmanager.Ascending)
			if err != nil {
				return err
			}

			return im.CreateSimpleSingles(ctx, collectionName, []indexmanager.SingleParam{
				{
					FieldOrder: fieldOrder,
					Name:       indexName,
				},
			})
		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)
			return im.DropMulti(ctx, collectionName, []string{indexName})
		},
	})
}
