package {{PROVIDER}}

import (
	"context"

	indexmanager "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/helper/index_manager"
	migrationregistry "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/registry"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	const collectionName = "credentials"

	migrationregistry.MigrationRecordRegister(migrationregistry.MigrationRecord{
		ID:   "20260128234925",
		Name: "credential_index",
		Up: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)

			fieldSupplier, err := indexmanager.NewFieldOrderRef("supplier", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldUsername, err := indexmanager.NewFieldOrderRef("username", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldDistributionType, err := indexmanager.NewFieldOrderRef("distributionType", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldIsDeleted, err := indexmanager.NewFieldOrderRef("isDeleted", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldStoreId, err := indexmanager.NewFieldOrderRef("storeId", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldUpdatedDate, err := indexmanager.NewFieldOrderRef("updateDate", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldExtendedData, err := indexmanager.NewFieldOrderRef("extendedData", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldPassword, err := indexmanager.NewFieldOrderRef("password", indexmanager.Ascending)
			if err != nil {
				return err
			}

			fieldIsStaging, err := indexmanager.NewFieldOrderRef("isStaging", indexmanager.Ascending)
			if err != nil {
				return err
			}

			return im.CreateSimpleSingles(ctx, collectionName, []indexmanager.SingleParam{
				{
					FieldOrder: fieldSupplier,
					Name:       "supplier_1",
				},
				{
					FieldOrder: fieldUsername,
					Name:       "username_1",
				},
				{
					FieldOrder: fieldDistributionType,
					Name:       "distributionType_1",
				},
				{
					FieldOrder: fieldIsDeleted,
					Name:       "isDeleted_1",
				},
				{
					FieldOrder: fieldStoreId,
					Name:       "storeId_1",
				},
				{
					FieldOrder: fieldUpdatedDate,
					Name:       "updateDate_1",
				},
				{
					FieldOrder: fieldExtendedData,
					Name:       "extendedData_1",
				},
				{
					FieldOrder: fieldPassword,
					Name:       "password_1",
				},
				{
					FieldOrder: fieldIsStaging,
					Name:       "isStaging_1",
				},
			})
		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			im := indexmanager.NewIndexManager(db)
			return im.DropMulti(ctx, collectionName, []string{
				"supplier_1",
				"username_1",
				"distributionType_1",
				"isDeleted_1",
				"storeId_1",
				"updateDate_1",
				"extendedData_1",
				"password_1",
				"isStaging_1",
			})
		},
	})
}
