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
		ID:   "20260128235948",
		Name: "system_parameter_index",
		Up: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)

			fieldVariable, err := indexmanager.NewFieldOrderRef("variable", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldIsDeleted, err := indexmanager.NewFieldOrderRef("isDeleted", indexmanager.Ascending)
			if err != nil {
				return err
			}

			return im.CreateSimpleSingles(ctx, collectionName, []indexmanager.SingleParam{
				{
					FieldOrder: fieldVariable,
					Name:       "variable_1",
				},
				{
					FieldOrder: fieldIsDeleted,
					Name:       "isDeleted_1",
				},
			})
		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)
			return im.DropMulti(ctx, collectionName, []string{
				"variable_1",
				"isDeleted_1",
			})
		},
	})
}
