package {{PROVIDER}}

import (
	"context"

	indexmanager "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/helper/index_manager"
	migrationregistry "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/registry"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	const collectionName = "system_parameter"

	migrationregistry.MigrationRecordRegister(migrationregistry.MigrationRecord{
		ID:   "20260129000326",
		Name: "system_parameter_index_2",
		Up: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)

			fieldVariableIsDeleted, err := indexmanager.NewFieldOrderRefs([]string{
				"variable",
				"isDeleted",
			}, []indexmanager.SortOrder{
				indexmanager.Ascending,
				indexmanager.Ascending,
			})
			if err != nil {
				return err
			}

			return im.CreateSimpleCompounds(ctx, collectionName, []indexmanager.CompoundParam{
				{
					FieldOrders: fieldVariableIsDeleted,
					Name:        "variable_1_isDeleted_1",
				},
			})
		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)
			return im.DropMulti(ctx, collectionName, []string{
				"variable_1_isDeleted_1",
			})
		},
	})
}
