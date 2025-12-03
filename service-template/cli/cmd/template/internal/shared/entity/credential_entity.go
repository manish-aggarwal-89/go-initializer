package entity

import (
	"time"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credential struct {
	DistributionType string            `json:"distributionType" validate:"required,omitempty"`
	Supplier         string            `json:"supplier" validate:"required,omitempty"`
	Username         string            `json:"username" validate:"required,omitempty"`
	Password         string            `json:"password"`
	ExpiredDate      string            `json:"expiredDate"`
	IsStaging        bool              `json:"isStaging"`
	ExtendedData     map[string]string `json:"extendedData"`
}

type CredentialRepository struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	DistributionType string             `json:"distributionType" bson:"distributionType"`
	Supplier         string             `json:"supplier" bson:"supplier"`
	Username         string             `json:"username" bson:"username"`
	Password         string             `json:"password" bson:"password"`
	ExpiredDate      *time.Time         `json:"expiredDate" bson:"expiredDate"`
	IsStaging        bool               `json:"isStaging" bson:"isStaging"`
	ExtendedData     map[string]string  `json:"extendedData" bson:"extendedData"`
	StoreId          string             `json:"storeId" bson:"storeId"`
	ChannelId        string             `json:"channelId" bson:"channelId"`
	IsDeleted        int32              `json:"isDeleted" bson:"isDeleted"`
	Version          int32              `json:"version" bson:"version"`
	CreateDate       time.Time          `json:"createDate" bson:"createdDate"`
	UpdateDate       time.Time          `json:"updateDate" bson:"updatedDate"`
	CreateBy         string             `json:"createBy" bson:"createdBy"`
	UpdateBy         string             `json:"updateBy" bson:"updatedBy"`
}

func (impl *Credential) DtoInsert(mr common.MandatoryRequest) CredentialRepository {
	expiredDate := util.ParseTime(util.DateFormat, impl.ExpiredDate)
	recordTime := time.Now()
	dto := CredentialRepository{}
	dto.ID = primitive.NewObjectID()
	dto.DistributionType = impl.DistributionType
	dto.Supplier = impl.Supplier
	dto.Username = impl.Username
	dto.Password = impl.Password
	dto.ExpiredDate = expiredDate
	dto.IsStaging = impl.IsStaging
	dto.ExtendedData = impl.ExtendedData
	dto.StoreId = mr.StoreId
	dto.ChannelId = mr.ChannelId
	dto.CreateBy = mr.Username
	dto.UpdateBy = mr.Username
	dto.CreateDate = recordTime
	dto.UpdateDate = recordTime
	return dto
}

func (impl CredentialRepository) DtoDelete(mr common.MandatoryRequest) CredentialRepository {
	impl.IsDeleted = 1
	impl.UpdateDate = time.Now()
	impl.UpdateBy = mr.Username
	impl.Version++
	return impl
}

func (impl CredentialRepository) DtoUpdate(data Credential, mr common.MandatoryRequest) CredentialRepository {
	expiredDate := util.ParseTime(util.DateFormat, data.ExpiredDate)
	impl.DistributionType = data.DistributionType
	impl.Supplier = data.Supplier
	impl.Username = data.Username
	impl.Password = data.Password
	impl.ExpiredDate = expiredDate
	impl.IsStaging = data.IsStaging
	impl.ExtendedData = data.ExtendedData
	impl.Version++
	impl.UpdateDate = time.Now()
	impl.UpdateBy = mr.Username
	return impl
}
