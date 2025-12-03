package dto

type SortType string

const (
	ID           SortType = "ID"
	SUPPLIER     SortType = "SUPPLIER"
	USERNAME     SortType = "USERNAME"
	DISTRIBUTION SortType = "DISTRIBUTION"
)

type OrderType string

const (
	ASC  OrderType = "ASC"
	DESC OrderType = "DESC"
)

type CredentialFilter struct {
	DistributionType string `json:"distributionType,omitempty"`
	Supplier         string `json:"supplier,omitempty"`
	Username         string `json:"username,omitempty"`
}

type CredentialLogResponse struct {
	LastUpdatedDate string `json:"lastUpdatedDate,omitempty"`
	LastUpdatedBy   string `json:"lastUpdatedBy,omitempty"`
	LogData         string `json:"logData,omitempty"`
}
