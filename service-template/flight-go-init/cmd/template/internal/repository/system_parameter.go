package repository

import (
	"context"
	"errors"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SysParamRepo struct {
	deps deps.Deps
}

type SystemParameterInterface interface {
	FindAll(ctx context.Context) ([]entity.SystemParameterRepository, error)
	FindAllPaginate(ctx context.Context, limit int64, page int64) ([]entity.SystemParameterRepository, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.SystemParameterRepository, error)
	FindByVariable(ctx context.Context, variable string) (*entity.SystemParameterRepository, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, data entity.SystemParameter, mr common.MandatoryRequest) (*entity.SystemParameterRepository, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, data entity.SystemParameter, mr common.MandatoryRequest) (*entity.SystemParameterRepository, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID, mr common.MandatoryRequest) error
}

func NewSystemParameterRepository(deps deps.Deps) SystemParameterInterface {
	return &SysParamRepo{deps: deps}
}

func (impl *SysParamRepo) FindAll(ctx context.Context) ([]entity.SystemParameterRepository, error) {
	dbCursor, err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).Find(ctx, bson.M{"isDeleted": 0})
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer dbCursor.Close(ctx)
	result := make([]entity.SystemParameterRepository, 0)
	for dbCursor.Next(ctx) {
		var row entity.SystemParameterRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.deps.GetLogger(ctx).Fatal(err.Error())
		}

		result = append(result, row)
	}
	return result, nil
}
func (impl *SysParamRepo) FindAllPaginate(ctx context.Context, limit int64, page int64) ([]entity.SystemParameterRepository, error) {
	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(limit)
	findOptions.SetSkip(page * limit)

	dbCursor, err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).Find(ctx, bson.M{"isDeleted": 0}, findOptions)
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer func(dbCursor *mongo.Cursor, ctx context.Context) {
		if err = dbCursor.Close(ctx); err != nil {
			impl.deps.GetLogger(ctx).Error("[system_param][FindAllPaginate] error closing db cursor:", err)
		}
	}(dbCursor, ctx)
	result := make([]entity.SystemParameterRepository, 0)
	for dbCursor.Next(ctx) {
		var row entity.SystemParameterRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.deps.GetLogger(ctx).Error(err.Error())
		}

		result = append(result, row)
	}
	return result, nil
}

func (impl *SysParamRepo) Count(ctx context.Context) (int64, error) {
	countOptions := options.EstimatedDocumentCountOptions{}
	res, err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).EstimatedDocumentCount(ctx, &countOptions)
	if err != nil {
		return 0, errors.New("count document error")
	}

	return res, nil
}

func (impl *SysParamRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.SystemParameterRepository, error) {
	result := entity.SystemParameterRepository{}
	err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		impl.deps.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return nil, err
	}

	return &result, nil
}

func (impl *SysParamRepo) FindByVariable(ctx context.Context, variable string) (*entity.SystemParameterRepository, error) {
	result := entity.SystemParameterRepository{}
	err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).FindOne(ctx, bson.M{"variable": variable, "isDeleted": 0}).Decode(&result)
	if err != nil {
		impl.deps.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return nil, err
	}

	return &result, nil
}

func (impl *SysParamRepo) Create(ctx context.Context, data entity.SystemParameter, mr common.MandatoryRequest) (*entity.SystemParameterRepository, error) {

	dataInsert := data.Dto(mr)
	_, err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).InsertOne(ctx, dataInsert)
	if err != nil {
		return nil, err
	}

	return &dataInsert, nil
}

func (impl *SysParamRepo) UpdateByID(ctx context.Context, id primitive.ObjectID, data entity.SystemParameter, mr common.MandatoryRequest) (*entity.SystemParameterRepository, error) {

	resultID, err := impl.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := resultID.Dto(data, mr)
	_, err = impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).UpdateOne(ctx, bson.M{"_id": resultID.ID}, bson.M{"$set": updateData})
	if err != nil {
		return nil, err
	}
	return &updateData, nil
}

func (impl *SysParamRepo) DeleteByID(ctx context.Context, id primitive.ObjectID, mr common.MandatoryRequest) error {
	_, err := impl.deps.Mongo.Collection(entity.SystemParameterCollectionName).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
