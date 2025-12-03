package repository

import (
	"context"
	"errors"
	"fmt"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/entity"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const currencyCollectionName = "currency"

type CurrencyRespositoryImpl struct {
	deps deps.Deps
}

type CurrencyRepositoryInterface interface {
	FindAll(ctx context.Context) ([]entity.CurrencyRepository, error)
	FindByCodeFrom(ctx context.Context, codeFrom string) (entity.CurrencyRepository, error)
	UpdateByID(ctx context.Context, currency entity.CurrencyRepository) error
	Create(ctx context.Context, mr common.MandatoryRequest, data entity.Currency) (entity.CurrencyRepository, error)
}

func NewCurrencyRepository(deps deps.Deps) CurrencyRepositoryInterface {
	return &CurrencyRespositoryImpl{deps: deps}
}

func (impl *CurrencyRespositoryImpl) FindAll(ctx context.Context) (ListCurrency []entity.CurrencyRepository, err error) {
	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(5)

	dbCursor, err := impl.deps.Mongo.Collection(currencyCollectionName).Find(ctx, bson.M{"isDeleted": false}, findOptions)
	if err != nil {
		return nil, errors.New("DATA IS EMPTY")
	}

	defer dbCursor.Close(ctx)
	result := make([]entity.CurrencyRepository, 0)
	for dbCursor.Next(ctx) {
		var row entity.CurrencyRepository
		err := dbCursor.Decode(&row)
		if err != nil {
			impl.deps.GetLogger(ctx).Fatal(err.Error())
		}

		result = append(result, row)
	}
	return result, nil
}

// FindByCodeFrom implements CurrencyRepositoryInterface.
func (impl *CurrencyRespositoryImpl) FindByCodeFrom(ctx context.Context, codeFrom string) (entity.CurrencyRepository, error) {
	result := entity.CurrencyRepository{}
	err := impl.deps.Mongo.Collection(currencyCollectionName).FindOne(ctx, bson.M{"codeFrom": codeFrom}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return entity.CurrencyRepository{}, fmt.Errorf("%w: %s", shared.ErrDataNotExist, err.Error())
	}
	if err != nil {
		impl.deps.GetLogger(ctx).Errorf("[error load data] = %s", err.Error())
		return entity.CurrencyRepository{}, err
	}

	return result, nil
}

// UpdateByID implements CurrencyRepositoryInterface.
func (impl *CurrencyRespositoryImpl) UpdateByID(ctx context.Context, currency entity.CurrencyRepository) error {
	_, err := impl.deps.Mongo.Collection(currencyCollectionName).UpdateOne(ctx, bson.M{"_id": currency.ID}, bson.M{"$set": currency})
	if err != nil {
		return err
	}
	return nil
}

// Create implements CurrencyRepositoryInterface.
func (impl *CurrencyRespositoryImpl) Create(ctx context.Context, mr common.MandatoryRequest, data entity.Currency) (entity.CurrencyRepository, error) {

	datainsert := data.DtoInsert(mr)
	_, err := impl.deps.Mongo.Collection(currencyCollectionName).InsertOne(ctx, datainsert)
	if err != nil {
		return entity.CurrencyRepository{}, err
	}
	return datainsert, nil
}
