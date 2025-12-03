package repository

import (
	"testing"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

func TestCredentialMasterData_Dto(t *testing.T) {
	// Create a mock MandatoryRequest object
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	// Create a mock SystemParameter object
	impl := &CredentialMasterData{
		DefaultValues:  []string{"string"},
		InputFieldType: "string",
		Key:            "string",
		Label:          "string",
		IsReadOnly:     true,
		IsRequired:     true,
	}

	// Call the Dto method
	credentialMasterData := impl.Dto(mr)

	// Check if the returned SystemParameter object is correct
	if credentialMasterData.ID.IsZero() {
		t.Errorf("Expected non-zero ID, but got zero")
	}
	if credentialMasterData.DefaultValues[0] != impl.DefaultValues[0] {
		t.Errorf("Expected DefaultValues %s, but got %s", impl.DefaultValues, credentialMasterData.DefaultValues)
	}
	if credentialMasterData.InputFieldType != impl.InputFieldType {
		t.Errorf("Expected InputFieldType %s, but got %s", impl.InputFieldType, credentialMasterData.InputFieldType)
	}
	if credentialMasterData.Key != impl.Key {
		t.Errorf("Expected Key %s, but got %s", impl.Key, credentialMasterData.Key)
	}
	if credentialMasterData.Label != impl.Label {
		t.Errorf("Expected Label %s, but got %s", impl.Label, credentialMasterData.Label)
	}
	if credentialMasterData.IsReadOnly != impl.IsReadOnly {
		t.Errorf("Expected IsReadOnly %t, but got %t", impl.IsReadOnly, credentialMasterData.IsReadOnly)
	}
	if credentialMasterData.IsRequired != impl.IsRequired {
		t.Errorf("Expected IsRequired %t, but got %t", impl.IsRequired, credentialMasterData.IsRequired)
	}
	if credentialMasterData.CreateBy != mr.Username {
		t.Errorf("Expected create by %s, but got %s", mr.Username, credentialMasterData.CreateBy)
	}
	if credentialMasterData.UpdateBy != mr.Username {
		t.Errorf("Expected update by %s, but got %s", mr.Username, credentialMasterData.UpdateBy)
	}
	if credentialMasterData.CreateDate.IsZero() {
		t.Errorf("Expected non-zero create date, but got zero")
	}
	if credentialMasterData.UpdateDate.IsZero() {
		t.Errorf("Expected non-zero update date, but got zero")
	}
}
