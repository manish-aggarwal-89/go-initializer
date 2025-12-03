package repository

import (
	"context"
	"errors"
	"time"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	credentialMasterDataCollectionName = "credential_master_data"
)

type CredentialMasterDataRepo struct {
	dep deps.Deps
}

type CredentialMasterDataRepository struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Key            string             `json:"key" bson:"key"`
	Label          string             `json:"label" bson:"label"`
	IsReadOnly     bool               `json:"isReadOnly" bson:"isReadOnly"`
	IsRequired     bool               `json:"isRequired" bson:"isRequired"`
	InputFieldType string             `json:"inputFieldType" bson:"inputFieldType"`
	DefaultValues  []string           `json:"defaultValues" bson:"defaultValues"`
	Version        int32              `json:"version" bson:"version"`
	CreateDate     time.Time          `json:"createDate" bson:"createdDate"`
	UpdateDate     time.Time          `json:"updateDate" bson:"updatedDate"`
	CreateBy       string             `json:"createBy" bson:"createdBy"`
	UpdateBy       string             `json:"updateBy" bson:"updatedBy"`
}

func (impl CredentialMasterDataRepository) Dto(data CredentialMasterData, mr common.MandatoryRequest) CredentialMasterDataRepository {
	impl.Key = data.Key
	impl.Label = data.Label
	impl.IsReadOnly = data.IsReadOnly
	impl.IsRequired = data.IsRequired
	impl.InputFieldType = data.InputFieldType
	impl.DefaultValues = data.DefaultValues
	impl.Version++
	impl.UpdateDate = time.Now()
	impl.UpdateBy = mr.Username
	return impl
}

func (impl CredentialMasterDataRepository) TableName() string {
	return credentialMasterDataCollectionName
}

type CredentialMasterDataInterface interface {
	FindAll(ctx context.Context) ([]CredentialMasterDataRepository, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*CredentialMasterDataRepository, error)
	FindByKey(ctx context.Context, key string) (*CredentialMasterDataRepository, error)
	Create(ctx context.Context, data CredentialMasterData, mr common.MandatoryRequest) (*CredentialMasterDataRepository, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, data CredentialMasterData, mr common.MandatoryRequest) (*CredentialMasterDataRepository, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

func NewCredentialMasterDataRepository(deps deps.Deps) CredentialMasterDataInterface {
	return &CredentialMasterDataRepo{dep: deps}
}

func (impl *CredentialMasterDataRepo) FindAll(ctx context.Context) ([]CredentialMasterDataRepository, error) {
	//find records

	dbCursor, err := impl.dep.Mongo.Collection(credentialMasterDataCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer func(dbCursor *mongo.Cursor, ctx context.Context) {
		err = dbCursor.Close(ctx)
		if err != nil {
			impl.dep.GetLogger(ctx).Errorf("[credential][FindAll] error: %s", err.Error())
		}
	}(dbCursor, ctx)
	result := make([]CredentialMasterDataRepository, 0)
	for dbCursor.Next(ctx) {
		var row CredentialMasterDataRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.dep.GetLogger(ctx).Fatal(err.Error())
		}

		result = append(result, row)
	}
	return result, nil
}

func (impl *CredentialMasterDataRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*CredentialMasterDataRepository, error) {
	var result CredentialMasterDataRepository
	err := impl.dep.Mongo.Collection(credentialMasterDataCollectionName).FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			impl.dep.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		}
		return nil, err
	}

	return &result, nil
}

func (impl *CredentialMasterDataRepo) FindByKey(ctx context.Context, key string) (*CredentialMasterDataRepository, error) {
	var result CredentialMasterDataRepository
	err := impl.dep.Mongo.Collection(credentialMasterDataCollectionName).FindOne(ctx, bson.M{"key": key}).Decode(&result)
	if err != nil {
		impl.dep.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return nil, err
	}

	return &result, nil
}

func (impl *CredentialMasterDataRepo) Create(ctx context.Context, data CredentialMasterData, mr common.MandatoryRequest) (*CredentialMasterDataRepository, error) {

	impl.dep.GetLogger(ctx).Info(util.ObjToJson(data.Dto(mr)))
	datainsert := data.Dto(mr)
	_, err := impl.dep.Mongo.Collection(credentialMasterDataCollectionName).InsertOne(ctx, datainsert)
	if err != nil {
		return nil, err
	}
	return &datainsert, nil
}

func (impl *CredentialMasterDataRepo) UpdateByID(ctx context.Context, id primitive.ObjectID, data CredentialMasterData, mr common.MandatoryRequest) (*CredentialMasterDataRepository, error) {

	existingData, err := impl.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newData := existingData.Dto(data, mr)
	_, err = impl.dep.Mongo.Collection(credentialMasterDataCollectionName).UpdateOne(ctx, bson.M{"_id": existingData.ID}, bson.M{"$set": newData})
	if err != nil {
		return nil, err
	}
	return &newData, nil
}

func (impl *CredentialMasterDataRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	existingData, err := impl.FindByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = impl.dep.Mongo.Collection(credentialMasterDataCollectionName).DeleteOne(ctx, bson.M{"_id": existingData.ID})
	if err != nil {
		return err
	}

	return nil
}
