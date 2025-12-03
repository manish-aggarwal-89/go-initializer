package entity

import (
	"time"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListCurrency []Currency

type Currency struct {
	CodeFrom string
	Rate     float32
}

func (impl Currency) DtoInsert(mr common.MandatoryRequest) CurrencyRepository {
	recordTime := time.Now()
	currency := CurrencyRepository{}
	currency.ID = primitive.NewObjectID()
	currency.CodeFrom = impl.CodeFrom
	currency.CodeTo = "IDR"
	currency.BuyRate = impl.Rate
	currency.SellRate = impl.Rate
	currency.CreateDate = recordTime
	currency.UpdateDate = recordTime
	currency.UpdateBy = mr.Username
	currency.CreateBy = mr.Username
	return currency
}

type CurrencyRepository struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	CodeFrom   string             `json:"codeFrom" bson:"codeFrom"`
	CodeTo     string             `json:"codeTo" bson:"codeTo"`
	BuyRate    float32            `json:"buyRate" bson:"buyRate,truncate"`
	SellRate   float32            `json:"sellRate" bson:"sellRate,truncate"`
	IsDeleted  bool               `json:"isDeleted" bson:"isDeleted"`
	Version    int32              `json:"version" bson:"version"`
	CreateDate time.Time          `json:"createDate" bson:"createDate"`
	UpdateDate time.Time          `json:"updateDate" bson:"updateDate"`
	CreateBy   string             `json:"createBy" bson:"createBy"`
	UpdateBy   string             `json:"updateBy" bson:"updateBy"`
}

func (impl CurrencyRepository) Dto(data Currency, mr common.MandatoryRequest) CurrencyRepository {
	impl.CodeFrom = data.CodeFrom
	impl.BuyRate = data.Rate
	impl.SellRate = data.Rate
	impl.Version++
	impl.UpdateDate = time.Now()
	impl.UpdateBy = mr.Username
	return impl
}

func (impl CurrencyRepository) ToCurrency() Currency {
	return Currency{
		CodeFrom: impl.CodeFrom,
		Rate:     impl.SellRate,
	}
}
