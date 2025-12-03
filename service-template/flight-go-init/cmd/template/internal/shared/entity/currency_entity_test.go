package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCurrency_DtoInsert(t *testing.T) {
	// Create a mock MandatoryRequest object
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	// Create a mock Currency object
	impl := Currency{
		CodeFrom: "USD",
		Rate:     1.23,
	}

	// Call the DtoInsert function
	currency := impl.DtoInsert(mr)

	// Check that the ID field is not empty
	if currency.ID.IsZero() {
		t.Errorf("Expected ID to be non-zero, but got zero")
	}

	// Check that the CodeFrom field is correct
	if currency.CodeFrom != "USD" {
		t.Errorf("Expected CodeFrom to be 'USD', but got '%s'", currency.CodeFrom)
	}

	// Check that the CodeTo field is correct
	if currency.CodeTo != "IDR" {
		t.Errorf("Expected CodeTo to be 'IDR', but got '%s'", currency.CodeTo)
	}

	// Check that the BuyRate field is correct
	if currency.BuyRate != 1.23 {
		t.Errorf("Expected BuyRate to be 1.23, but got %f", currency.BuyRate)
	}

	// Check that the SellRate field is correct
	if currency.SellRate != 1.23 {
		t.Errorf("Expected SellRate to be 1.23, but got %f", currency.SellRate)
	}

	// Check that the CreateDate field is not zero
	if currency.CreateDate.IsZero() {
		t.Errorf("Expected CreateDate to be non-zero, but got zero")
	}

	// Check that the UpdateDate field is not zero
	if currency.UpdateDate.IsZero() {
		t.Errorf("Expected UpdateDate to be non-zero, but got zero")
	}

	// Check that the UpdateBy field is correct
	if currency.UpdateBy != "testuser" {
		t.Errorf("Expected UpdateBy to be 'testuser', but got '%s'", currency.UpdateBy)
	}

	// Check that the CreateBy field is correct
	if currency.CreateBy != "testuser" {
		t.Errorf("Expected CreateBy to be 'testuser', but got '%s'", currency.CreateBy)
	}
}

func TestCurrencyRepository_Dto(t *testing.T) {
	dataCreateTime := time.Now().Add(-1 * time.Hour)
	// Create a new instance of CurrencyRepository
	repo := CurrencyRepository{
		ID:         primitive.NewObjectID(),
		CodeFrom:   "USD",
		CodeTo:     "EUR",
		BuyRate:    1.23,
		SellRate:   1.24,
		IsDeleted:  false,
		Version:    1,
		CreateDate: time.Now(),
		UpdateDate: time.Now(),
		CreateBy:   "testuser",
		UpdateBy:   "testuser",
	}

	// Create a new instance of Currency
	data := Currency{
		CodeFrom: "USD",
		Rate:     1.23,
	}

	// Create a new instance of MandatoryRequest
	mr := common.MandatoryRequest{
		Username: "testuser",
	}

	// Call the Dto method on the repository
	updatedRepo := repo.Dto(data, mr)

	// Check that the CodeFrom field was updated correctly
	if updatedRepo.CodeFrom != "USD" {
		t.Errorf("CodeFrom field was not updated correctly")
	}

	// Check that the BuyRate and SellRate fields were updated correctly
	if updatedRepo.BuyRate != 1.23 {
		t.Errorf("BuyRate field was not updated correctly")
	}
	if updatedRepo.SellRate != 1.23 {
		t.Errorf("SellRate field was not updated correctly")
	}

	// Check that the Version field was incremented correctly
	if updatedRepo.Version != 2 {
		t.Errorf("Version field was not incremented correctly")
	}

	// Check that the UpdateDate and UpdateBy fields were updated correctly
	if updatedRepo.UpdateDate.Before(dataCreateTime) {
		t.Errorf("UpdateDate field was not updated correctly")
	}
	if updatedRepo.UpdateBy != "testuser" {
		t.Errorf("UpdateBy field was not updated correctly")
	}
}

func TestCurrencyRepository_ToCurrency(t *testing.T) {
	t.Run("it should convert repository to currency", func(t *testing.T) {
		// Arrange
		repo := CurrencyRepository{
			CodeFrom: "USD",
			SellRate: 1.2,
		}

		expectedCurrency := Currency{
			CodeFrom: "USD",
			Rate:     1.2,
		}

		// Act
		actualCurrency := repo.ToCurrency()

		// Assert
		assert.Equal(t, expectedCurrency, actualCurrency)
	})
}
