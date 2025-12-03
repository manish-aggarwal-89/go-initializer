package entity

import (
	"time"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CredentialLog struct {
	Value *CredentialRepository
}

type CredentialLogRepository struct {
	ID         primitive.ObjectID    `bson:"_id" json:"id"`
	Value      *CredentialRepository `json:"value" bson:"value"`
	StoreId    string                `json:"storeId" bson:"storeId"`
	IsDeleted  int16                 `json:"isDeleted" bson:"isDeleted"`
	Version    int32                 `json:"version" bson:"version"`
	CreateDate time.Time             `json:"createDate" bson:"createdDate"`
	UpdateDate time.Time             `json:"updateDate" bson:"updatedDate"`
	CreateBy   string                `json:"createBy" bson:"createdBy"`
	UpdateBy   string                `json:"updateBy" bson:"updatedBy"`
}

func (impl *CredentialLog) DtoInsert(mr common.MandatoryRequest) CredentialLogRepository {
	recordTime := time.Now()
	dto := CredentialLogRepository{}
	dto.ID = primitive.NewObjectID()
	dto.Value = impl.Value
	dto.CreateBy = mr.Username
	dto.UpdateBy = mr.Username
	dto.CreateDate = recordTime
	dto.UpdateDate = recordTime
	return dto
}
