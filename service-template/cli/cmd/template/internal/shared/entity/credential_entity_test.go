package entity

import (
	"testing"
	"time"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

func TestCredential_DtoInsert(t *testing.T) {
	// Create a mock MandatoryRequest object
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	m := make(map[string]string)
	m["string"] = "string"

	// Create a mock Credential object
	impl := &Credential{
		DistributionType: "string",
		Supplier:         "string",
		Username:         "string",
		Password:         "string",
		ExpiredDate:      "2023-09-25",
		ExtendedData:     m,
	}

	// Call the Dto method
	credential := impl.DtoInsert(mr)

	// Check if the returned Credential object is correct
	if credential.ID.IsZero() {
		t.Errorf("Expected non-zero ID, but got zero")
	}
	if credential.DistributionType != impl.DistributionType {
		t.Errorf("Expected DistributionType %s, but got %s", impl.DistributionType, credential.DistributionType)
	}
	if credential.Supplier != impl.Supplier {
		t.Errorf("Expected Supplier %s, but got %s", impl.Supplier, credential.Supplier)
	}
	if credential.Username != impl.Username {
		t.Errorf("Expected Username %s, but got %s", impl.Username, credential.Username)
	}
	if credential.Password != impl.Password {
		t.Errorf("Expected Password %s, but got %s", impl.Password, credential.Password)
	}
	expiredDate := util.FormatTime(util.DateFormat, *credential.ExpiredDate)
	if expiredDate != impl.ExpiredDate {
		t.Errorf("Expected ExpiredDate %s, but got %s", impl.ExpiredDate, expiredDate)
	}
	if credential.ExtendedData["string"] != impl.ExtendedData["string"] {
		t.Errorf("Expected ExtendedData %s, but got %s", impl.ExtendedData["string"], credential.ExtendedData["string"])
	}
	if credential.CreateBy != mr.Username {
		t.Errorf("Expected create by %s, but got %s", mr.Username, credential.CreateBy)
	}
	if credential.UpdateBy != mr.Username {
		t.Errorf("Expected update by %s, but got %s", mr.Username, credential.UpdateBy)
	}
	if credential.CreateDate.IsZero() {
		t.Errorf("Expected non-zero create date, but got zero")
	}
	if credential.UpdateDate.IsZero() {
		t.Errorf("Expected non-zero update date, but got zero")
	}
}

func TestCredentialRepository_DtoUpdate(t *testing.T) {
	// Create a mock MandatoryRequest object
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	timeNow := time.Now()
	timeCreateData := time.Now().Add(-1 * time.Hour)

	m := make(map[string]string)
	m["string"] = "string"

	impl := &Credential{
		DistributionType: "string",
		Supplier:         "string",
		Username:         "string",
		Password:         "string",
		ExpiredDate:      "2023-09-25",
		ExtendedData:     m,
	}

	// Create a mock Credential object
	repo := &CredentialRepository{
		DistributionType: "string",
		Supplier:         "string",
		Username:         "string",
		Password:         "string",
		ExpiredDate:      &timeNow,
		ExtendedData:     m,
		IsDeleted:        0,
		Version:          1,
		CreateDate:       timeNow,
		UpdateDate:       timeNow,
		CreateBy:         mr.Username,
		UpdateBy:         mr.Username,
	}

	// Call the Dto method
	updateRepo := repo.DtoUpdate(*impl, mr)

	// Check if the returned Credential object is correct
	if updateRepo.DistributionType != impl.DistributionType {
		t.Errorf("Expected DistributionType %s, but got %s", impl.DistributionType, updateRepo.DistributionType)
	}
	if updateRepo.Supplier != impl.Supplier {
		t.Errorf("Expected Supplier %s, but got %s", impl.Supplier, updateRepo.Supplier)
	}
	if updateRepo.Username != impl.Username {
		t.Errorf("Expected Username %s, but got %s", impl.Username, updateRepo.Username)
	}
	if updateRepo.Password != impl.Password {
		t.Errorf("Expected Password %s, but got %s", impl.Password, updateRepo.Password)
	}
	expiredDate := util.FormatTime(util.DateFormat, *updateRepo.ExpiredDate)
	if expiredDate != impl.ExpiredDate {
		t.Errorf("Expected ExpiredDate %s, but got %s", impl.ExpiredDate, expiredDate)
	}
	if updateRepo.ExtendedData["string"] != impl.ExtendedData["string"] {
		t.Errorf("Expected ExtendedData %s, but got %s", impl.ExtendedData["string"], updateRepo.ExtendedData["string"])
	}

	// Check that the Version field was incremented correctly
	if updateRepo.Version != 2 {
		t.Errorf("Version field was not incremented correctly")
	}

	// Check that the UpdateDate and UpdateBy fields were updated correctly
	if updateRepo.UpdateDate.Before(timeCreateData) {
		t.Errorf("UpdateDate field was not updated correctly")
	}
	if updateRepo.UpdateBy != mr.Username {
		t.Errorf("Expected update by %s, but got %s", mr.Username, updateRepo.UpdateBy)
	}
	if updateRepo.CreateBy != mr.Username {
		t.Errorf("Expected create by %s, but got %s", mr.Username, updateRepo.CreateBy)
	}
}

func TestCredentialRepository_DtoDelete(t *testing.T) {
	// Create a mock MandatoryRequest object
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	timeNow := time.Now()
	timeCreateData := time.Now().Add(-1 * time.Hour)

	// Create a mock Credential object
	repo := &CredentialRepository{
		IsDeleted:  0,
		Version:    1,
		CreateDate: timeNow,
		UpdateDate: timeNow,
		CreateBy:   mr.Username,
		UpdateBy:   mr.Username,
	}

	// Call the Dto method
	deleteRepo := repo.DtoDelete(mr)

	// Check that the Version field was incremented correctly
	if deleteRepo.Version != 2 {
		t.Errorf("Version field was not incremented correctly")
	}

	if deleteRepo.IsDeleted != 1 {
		t.Errorf("IsDeleted field was not updated correctly")
	}

	// Check that the UpdateDate and UpdateBy fields were updated correctly
	if deleteRepo.UpdateDate.Before(timeCreateData) {
		t.Errorf("UpdateDate field was not updated correctly")
	}
	if deleteRepo.UpdateBy != mr.Username {
		t.Errorf("UpdateBy field was not updated correctly")
	}
}
