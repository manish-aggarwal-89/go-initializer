package search

import (
	"{{MODULE_NAME}}/internal/shared/test_var"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
)

var (
	INTEGRATOR_FARE_REQUEST_CGK_DPS_2A2C_RT = fareRQ.IntegratorFareRequest{
		DistributionType: "flyjaya",
		SupplierMappings: []fareRQ.SupplierMapping{test_var.SUPPLIER_MAPPING},
		SupplierRequest: fareRQ.SupplierRequest{
			Accounts: []fareRQ.Account{
				test_var.ACCOUNT,
			},
			Airlines: []string{"KS"},
		},
		Origin:        "CGK",
		Destination:   "YIA",
		Adult:         2,
		Child:         2,
		Infant:        0,
		DepartureDate: "2026-01-01",
		ReturnDate:    "2026-01-03",
		CabinClass:    "ECONOMY",
		TripTypes:     []enum.TripType{enum.ROUND_TRIP},
		Currency: fareRQ.Currency{
			DepartureCurrency: "IDR",
		},
		International: 1,
		Routes: []fareRQ.Route{
			{
				Origin:        "CGK",
				Destination:   "YIA",
				DepartureDate: "2026-01-01",
				Currency:      "IDR",
			},
			{
				Origin:        "YIA",
				Destination:   "CGK",
				DepartureDate: "2026-01-03",
				Currency:      "IDR",
			},
		},
	}
)
