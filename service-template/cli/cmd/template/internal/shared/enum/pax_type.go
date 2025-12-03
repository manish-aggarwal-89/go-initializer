package enum

import "strings"

type PaxType struct {
	Code string
	Name string
}

type SupplierPaxType string

const (
	PaxTypeAdult  SupplierPaxType = "ADULT"
	PaxTypeChild  SupplierPaxType = "CHD"
	PaxTypeInfant SupplierPaxType = "INF"
)

var (
	ADT = PaxType{"ADT", "ADULT"}
	CHD = PaxType{"CHD", "CHILD"}
	INF = PaxType{"INF", "INFANT"}
)

func ConvertPaxType(paxType string) SupplierPaxType {
	paxTypeUppercase := strings.ToUpper(paxType)
	switch paxTypeUppercase {
	case ADT.Name:
		return PaxTypeAdult
	case CHD.Name:
		return PaxTypeChild
	case INF.Name:
		return PaxTypeInfant
	default:
		return ""
	}
}
