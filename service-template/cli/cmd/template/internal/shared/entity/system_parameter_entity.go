package entity

import (
	"time"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SystemParameterCollectionName = "system_parameter"
)

type SystemParameterRepository struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Variable    string             `json:"variable" bson:"variable"`
	Value       string             `json:"value" bson:"value"`
	Description string             `json:"description" bson:"description"`
	IsDeleted   int16              `json:"isDeleted" bson:"isDeleted"`
	Version     int32              `json:"version" bson:"version"`
	CreateDate  time.Time          `json:"createDate" bson:"createdDate"`
	UpdateDate  time.Time          `json:"updateDate" bson:"updatedDate"`
	CreateBy    string             `json:"createBy" bson:"createdBy"`
	UpdateBy    string             `json:"updateBy" bson:"updatedBy"`
}

func (impl SystemParameterRepository) Dto(data SystemParameter, mr common.MandatoryRequest) SystemParameterRepository {
	impl.Value = data.Value
	impl.Description = data.Description
	impl.Variable = data.Variable
	impl.Version++
	impl.UpdateDate = time.Now()
	impl.UpdateBy = mr.Username
	return impl
}

func (impl SystemParameterRepository) DtoDeleted(mr common.MandatoryRequest) SystemParameterRepository {
	impl.IsDeleted = 1
	impl.UpdateDate = time.Now()
	impl.UpdateBy = mr.Username
	impl.Version++
	return impl
}

func (impl SystemParameterRepository) TableName() string {
	return SystemParameterCollectionName
}

type SystemParameter struct {
	Variable    string `json:"variable" bson:"variable"`
	Value       string `json:"value" bson:"value"`
	Description string `json:"description" bson:"description"`
}

func (impl *SystemParameter) Dto(mr common.MandatoryRequest) SystemParameterRepository {
	recordTime := time.Now()
	sysParam := SystemParameterRepository{}
	sysParam.ID = primitive.NewObjectID()
	sysParam.Description = impl.Description
	sysParam.Value = impl.Value
	sysParam.Variable = impl.Variable
	sysParam.CreateBy = mr.Username
	sysParam.UpdateBy = mr.Username
	sysParam.CreateDate = recordTime
	sysParam.UpdateDate = recordTime

	return sysParam
}
