package repository

import (
	"context"
	"errors"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/dto"
	"{{MODULE_NAME}}/internal/shared/entity"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const credentialCollectionName = "credential"

type CredentialRespositoryImpl struct {
	dep deps.Deps
}

type CredentialRepositoryInterface interface {
	FindAll(ctx context.Context) ([]entity.CredentialRepository, error)
	FindAllPaginate(ctx context.Context, cf dto.CredentialFilter, page int64, size int64, sort string, sortDirect string) ([]entity.CredentialRepository, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.CredentialRepository, error)
	FindBySupplier(ctx context.Context, supplier string) (*entity.CredentialRepository, error)
	FindBySupplierAndIsStaging(ctx context.Context, supplier string, isStaging bool) (*entity.CredentialRepository, error)
	Create(ctx context.Context, data entity.Credential, mr common.MandatoryRequest) (*entity.CredentialRepository, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, data entity.Credential, mr common.MandatoryRequest) (*entity.CredentialRepository, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID, mr common.MandatoryRequest) (*entity.CredentialRepository, error)
	Count(ctx context.Context) (int64, error)
}

func NewCredentialRepository(deps deps.Deps) CredentialRepositoryInterface {
	return &CredentialRespositoryImpl{dep: deps}
}

// Create implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) Create(ctx context.Context, data entity.Credential, mr common.MandatoryRequest) (*entity.CredentialRepository, error) {
	dataInsert := data.DtoInsert(mr)
	_, err := impl.dep.Mongo.Collection(credentialCollectionName).InsertOne(ctx, dataInsert)
	if err != nil {
		return nil, err
	}
	return &dataInsert, nil
}

// DeleteByID implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) DeleteByID(ctx context.Context, id primitive.ObjectID, mr common.MandatoryRequest) (*entity.CredentialRepository, error) {
	resultID, err := impl.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	deleteData := resultID.DtoDelete(mr)
	_, err = impl.dep.Mongo.Collection(credentialCollectionName).UpdateOne(ctx, bson.M{"_id": resultID.ID}, bson.M{"$set": deleteData})
	if err != nil {
		return nil, err
	}
	return &deleteData, nil
}

// FindAll implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) FindAll(ctx context.Context) ([]entity.CredentialRepository, error) {
	//find records
	//pass these options to the Find method
	findOptions := options.Find()

	dbCursor, err := impl.dep.Mongo.Collection(credentialCollectionName).Find(ctx,
		bson.M{"isDeleted": 0, "isStaging": impl.dep.Config.IsStaging}, findOptions)
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer func(dbCursor *mongo.Cursor, ctx context.Context) {
		err = dbCursor.Close(ctx)
		if err != nil {
			impl.dep.GetLogger(ctx).Errorf("[credential][findAll] error: %s", err.Error())
		}
	}(dbCursor, ctx)
	result := make([]entity.CredentialRepository, 0)
	for dbCursor.Next(ctx) {
		var row entity.CredentialRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.dep.GetLogger(ctx).Errorf("[credential][findAll] error: %s", err.Error())
			break
		}

		result = append(result, row)
	}
	return result, nil
}

func (impl *CredentialRespositoryImpl) FindAllPaginate(ctx context.Context, cf dto.CredentialFilter, page int64, size int64, sort string, sortDirect string) ([]entity.CredentialRepository, error) {

	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(size)
	findOptions.SetSkip(page * size)
	findOptions.SetSort(map[string]int{getSortValue(sort): getOrderValue(sortDirect)})
	findOptions.SetCollation(&options.Collation{
		Locale: "en",
	})

	// filters contains all the filters you want to apply:
	filters := []bson.M{
		{"isDeleted": 0},
		{"isStaging": impl.dep.Config.IsStaging},
		{"distributionType": bson.M{"$regex": cf.DistributionType, "$options": "i"}},
		{"supplier": bson.M{"$regex": cf.Supplier, "$options": "i"}},
		{"username": bson.M{"$regex": cf.Username, "$options": "i"}},
	}

	// filter is a single filter document that merges all filters
	filter := bson.M{"$and": filters}

	dbCursor, err := impl.dep.Mongo.Collection(credentialCollectionName).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer func(dbCursor *mongo.Cursor, ctx context.Context) {
		err = dbCursor.Close(ctx)
		if err != nil {
			impl.dep.GetLogger(ctx).Errorf("[credential][FindAllPaginate] error: %s", err.Error())
		}
	}(dbCursor, ctx)
	result := make([]entity.CredentialRepository, 0)
	for dbCursor.Next(ctx) {
		var row entity.CredentialRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.dep.GetLogger(ctx).Errorf("[credential][FindAllPaginate] error: %s", err.Error())
			break
		}

		result = append(result, row)
	}
	return result, nil
}

// FindByID implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.CredentialRepository, error) {
	result := entity.CredentialRepository{}
	err := impl.dep.Mongo.Collection(credentialCollectionName).FindOne(ctx, bson.M{"_id": id, "isDeleted": 0, "isStaging": impl.dep.Config.IsStaging}).Decode(&result)
	if err != nil {
		impl.dep.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return nil, err
	}

	return &result, nil
}

// FindBySupplier implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) FindBySupplier(ctx context.Context, supplier string) (*entity.CredentialRepository, error) {
	result := entity.CredentialRepository{}
	err := impl.dep.Mongo.Collection(credentialCollectionName).FindOne(ctx, bson.M{"supplier": supplier, "isDeleted": 0, "isStaging": impl.dep.Config.IsStaging}).Decode(&result)
	if err != nil {
		impl.dep.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return nil, err
	}

	return &result, nil
}

// FindBySupplierAndIsStaging implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) FindBySupplierAndIsStaging(ctx context.Context,
	supplier string, isStaging bool) (*entity.CredentialRepository, error) {
	result := entity.CredentialRepository{}
	err := impl.dep.Mongo.Collection(credentialCollectionName).FindOne(ctx, bson.M{"supplier": supplier, "isDeleted": 0, "isStaging": isStaging}).Decode(&result)
	if err != nil {
		impl.dep.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return nil, err
	}

	return &result, nil
}

// UpdateByID implements CredentialRepositoryInterface.
func (impl *CredentialRespositoryImpl) UpdateByID(ctx context.Context, id primitive.ObjectID, data entity.Credential, mr common.MandatoryRequest) (*entity.CredentialRepository, error) {
	resultID, err := impl.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := resultID.DtoUpdate(data, mr)
	_, err = impl.dep.Mongo.Collection(credentialCollectionName).UpdateOne(ctx, bson.M{"_id": resultID.ID}, bson.M{"$set": updateData})
	if err != nil {
		return nil, err
	}
	return &updateData, nil
}

func (impl *CredentialRespositoryImpl) Count(ctx context.Context) (int64, error) {
	res, err := impl.dep.Mongo.Collection(credentialCollectionName).CountDocuments(ctx, bson.M{"isDeleted": 0, "isStaging": impl.dep.Config.IsStaging})
	if err != nil {
		impl.dep.GetLogger(ctx).Errorf("[error count data] = %s", err.Error())
		return 0, errors.New("count document error")
	}

	return res, nil
}

func getOrderValue(sortDirect string) int {
	//Order -1 = DESC, 1 = ASC
	if string(dto.ASC) == sortDirect {
		return 1
	}
	return -1
}

func getSortValue(sort string) string {
	if string(dto.DISTRIBUTION) == sort {
		return "distributionType"
	}

	if string(dto.SUPPLIER) == sort {
		return "supplier"
	}

	if string(dto.USERNAME) == sort {
		return "username"
	}

	return "_id"
}
