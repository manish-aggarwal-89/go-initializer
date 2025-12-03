package entity

import (
	"testing"
	"time"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

func TestCredentialLog_DtoInsert(t *testing.T) {
	// Create a mock MandatoryRequest object
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	timeNow := time.Now()

	m := make(map[string]string)
	m["string"] = "string"

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

	impl := &CredentialLog{
		Value: repo,
	}

	// Call the Dto method
	credentialLog := impl.DtoInsert(mr)

	// Check if the returned Credential object is correct
	if credentialLog.ID.IsZero() {
		t.Errorf("Expected non-zero ID, but got zero")
	}
	if credentialLog.Value.DistributionType != repo.DistributionType {
		t.Errorf("Expected DistributionType %s, but got %s", repo.DistributionType, credentialLog.Value.DistributionType)
	}
	if credentialLog.Value.Supplier != repo.Supplier {
		t.Errorf("Expected Supplier %s, but got %s", repo.Supplier, credentialLog.Value.Supplier)
	}
	if credentialLog.Value.Username != repo.Username {
		t.Errorf("Expected Username %s, but got %s", repo.Username, credentialLog.Value.Username)
	}
	if credentialLog.Value.Password != repo.Password {
		t.Errorf("Expected Password %s, but got %s", repo.Password, credentialLog.Value.Password)
	}
	if credentialLog.Value.ExpiredDate.Compare(*repo.ExpiredDate) != 0 {
		t.Errorf("Expected ExpiredDate %s, but got %s", repo.ExpiredDate, credentialLog.Value.ExpiredDate)
	}
	if credentialLog.Value.ExtendedData["string"] != repo.ExtendedData["string"] {
		t.Errorf("Expected ExtendedData %s, but got %s", repo.ExtendedData["string"], credentialLog.Value.ExtendedData["string"])
	}
	if credentialLog.CreateBy != mr.Username {
		t.Errorf("Expected create by %s, but got %s", mr.Username, credentialLog.CreateBy)
	}
	if credentialLog.UpdateBy != mr.Username {
		t.Errorf("Expected update by %s, but got %s", mr.Username, credentialLog.UpdateBy)
	}
	if credentialLog.CreateDate.IsZero() {
		t.Errorf("Expected non-zero create date, but got zero")
	}
	if credentialLog.UpdateDate.IsZero() {
		t.Errorf("Expected non-zero update date, but got zero")
	}
}
