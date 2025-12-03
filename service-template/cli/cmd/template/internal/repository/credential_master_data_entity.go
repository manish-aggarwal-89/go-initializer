package repository

import (
	"time"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CredentialMasterData struct {
	ID             string   `json:"id"`
	Key            string   `json:"key" bson:"key"`
	Label          string   `json:"label" bson:"label"`
	IsReadOnly     bool     `json:"isReadOnly" bson:"isReadOnly"`
	IsRequired     bool     `json:"isRequired" bson:"isRequired"`
	InputFieldType string   `json:"inputFieldType" bson:"inputFieldType"`
	DefaultValues  []string `json:"defaultValues" bson:"defaultValues"`
}

func (impl *CredentialMasterData) Dto(mr common.MandatoryRequest) CredentialMasterDataRepository {
	recordTime := time.Now()
	credentialMasterData := CredentialMasterDataRepository{}
	credentialMasterData.ID = primitive.NewObjectID()
	credentialMasterData.Key = impl.Key
	credentialMasterData.Label = impl.Label
	credentialMasterData.IsReadOnly = impl.IsReadOnly
	credentialMasterData.IsRequired = impl.IsRequired
	credentialMasterData.DefaultValues = impl.DefaultValues
	credentialMasterData.InputFieldType = impl.InputFieldType
	credentialMasterData.CreateBy = mr.Username
	credentialMasterData.UpdateBy = mr.Username
	credentialMasterData.CreateDate = recordTime
	credentialMasterData.UpdateDate = recordTime

	return credentialMasterData
}
