package {{PROVIDER}}

import (
	"context"

	indexmanager "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/helper/index_manager"
	migrationregistry "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/registry"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	const collectionName = "integrator_error_mapping"
	const indexName = "is_deleted_asc_key_asc_bau_process_name_asc"

	migrationregistry.MigrationRecordRegister(migrationregistry.MigrationRecord{
		ID:   "20260122161658",
		Name: "integrator_error_mapping_key_bau_process_index",
		Up: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)

			fieldOrders, err := indexmanager.NewFieldOrderRefs([]string{
				"isDeleted",
				"key",
				"bauProcessName",
			}, []indexmanager.SortOrder{
				indexmanager.Ascending,
				indexmanager.Ascending,
				indexmanager.Ascending,
			})
			if err != nil {
				return err
			}

			return im.CreateSimpleCompounds(ctx, collectionName, []indexmanager.CompoundParam{
				{
					FieldOrders: fieldOrders,
					Name:        indexName,
				},
			})
		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)
			return im.DropMulti(ctx, collectionName, []string{indexName})
		},
	})
}
