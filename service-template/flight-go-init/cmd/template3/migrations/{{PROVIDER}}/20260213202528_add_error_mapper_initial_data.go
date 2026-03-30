package {{PROVIDER}}

import (
	"context"
	"log"
	"time"

	migrationregistry "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/db-tools/mongo/migration/registry"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	const collectionName = "integrator_error_mapping"
	timeNow := time.Now()

	type integratorErrorMappingRecord struct {
		Key            string    `bson:"key"`
		ErrorMessage   string    `bson:"errorMessage"`
		HealthStatus   string    `bson:"healthStatus"`
		BauProcessName string    `bson:"bauProcessName"`
		Version        int       `bson:"version"`
		CreatedDate    time.Time `bson:"createdDate"`
		UpdatedDate    time.Time `bson:"updatedDate"`
		CreatedBy      string    `bson:"createdBy"`
		UpdatedBy      string    `bson:"updatedBy"`
		StoreId        string    `bson:"storeId"`
		ChannelId      string    `bson:"channelId"`
		IsDeleted      int       `bson:"isDeleted"`
	}

	records := []integratorErrorMappingRecord{
		{
			Key:            "flight not found",
			ErrorMessage:   "flight not found",
			HealthStatus:   "FLIGHT_NOT_FOUND",
			BauProcessName: "SEARCH",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
		{
			Key:            "invalid route",
			ErrorMessage:   "invalid route",
			HealthStatus:   "INVALID_ROUTE",
			BauProcessName: "SEARCH",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
		{
			Key:            "flight not found",
			ErrorMessage:   "flight not found",
			HealthStatus:   "NO_SEAT",
			BauProcessName: "BOOKING",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
		{
			Key:            "no seat",
			ErrorMessage:   "no seat",
			HealthStatus:   "NO_SEAT",
			BauProcessName: "BOOKING",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
		{
			Key:            "baggage not available",
			ErrorMessage:   "baggage not available",
			HealthStatus:   "BAGGAGE_NOT_AVAILABLE",
			BauProcessName: "BOOKING",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
		{
			Key:            "meals not available",
			ErrorMessage:   "meals not available",
			HealthStatus:   "MEALS_NOT_AVAILABLE",
			BauProcessName: "BOOKING",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
		{
			Key:            "seat not available",
			ErrorMessage:   "seat not available",
			HealthStatus:   "SEAT_NOT_AVAILABLE",
			BauProcessName: "BOOKING",
			Version:        1,
			CreatedDate:    timeNow,
			UpdatedDate:    timeNow,
			CreatedBy:      "migrations_script",
			UpdatedBy:      "migrations_script",
			StoreId:        "TIKETCOM",
			ChannelId:      "WEB",
			IsDeleted:      0,
		},
	}

	migrationregistry.MigrationRecordRegister(migrationregistry.MigrationRecord{
		ID:   "20260213202528",
		Name: "add_error_mapper_initial_data",
		Up: func(ctx context.Context, db *mongo.Database) error {
			mongoRecords := make([]interface{}, len(records))
			for i, r := range records {
				mongoRecords[i] = r
			}
			_, err := db.Collection(collectionName).InsertMany(ctx, mongoRecords)
			return err
		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			for _, r := range records {
				filter := bson.M{
					"key":            r.Key,
					"errorMessage":   r.ErrorMessage,
					"healthStatus":   r.HealthStatus,
					"bauProcessName": r.BauProcessName,
				}
				_, err := db.Collection(collectionName).DeleteOne(ctx, filter)
				if err != nil {
					log.Println("Failed to delete record during rollback:", err)
				}
			}
			return nil
		},
	})
}
