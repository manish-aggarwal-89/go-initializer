package repository

import (
	"context"
	"errors"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const credentialLogCollectionName = "credential_log"

type CredentialLogRespositoryImpl struct {
	dep deps.Deps
}

type CredentialLogRepositoryInterface interface {
	FindByID(ctx context.Context, id primitive.ObjectID) ([]entity.CredentialLogRepository, error)
	Create(ctx context.Context, data entity.CredentialLog, mr common.MandatoryRequest) (*entity.CredentialLogRepository, error)
}

func NewCredentialLogRepository(deps deps.Deps) CredentialLogRepositoryInterface {
	return &CredentialLogRespositoryImpl{dep: deps}
}

// Create implements CredentialLogRepositoryInterface.
func (impl *CredentialLogRespositoryImpl) Create(ctx context.Context, data entity.CredentialLog, mr common.MandatoryRequest) (*entity.CredentialLogRepository, error) {
	datainsert := data.DtoInsert(mr)
	_, err := impl.dep.Mongo.Collection(credentialLogCollectionName).InsertOne(ctx, datainsert)
	if err != nil {
		return nil, err
	}
	return &datainsert, nil
}

// FindByID implements CredentialLogRepositoryInterface.
func (impl *CredentialLogRespositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) ([]entity.CredentialLogRepository, error) {
	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	// findOptions.SetSort(bson.D{{"updateDate", -1}})
	findOptions.SetSort(map[string]int{"updateDate": -1})

	dbCursor, err := impl.dep.Mongo.Collection(credentialLogCollectionName).Find(ctx, bson.M{"value._id": id}, findOptions)
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer func(dbCursor *mongo.Cursor, ctx context.Context) {
		err = dbCursor.Close(ctx)
		if err != nil {
			impl.dep.GetLogger(ctx).Errorf("[credentialLog][FindByID] error: %s", err.Error())
		}
	}(dbCursor, ctx)
	result := make([]entity.CredentialLogRepository, 0)
	for dbCursor.Next(ctx) {
		var row entity.CredentialLogRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.dep.GetLogger(ctx).Fatal(err.Error())
		}

		result = append(result, row)
	}
	return result, nil
}
