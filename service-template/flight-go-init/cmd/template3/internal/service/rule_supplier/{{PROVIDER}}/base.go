package {{PROVIDER}}

import (
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type Base struct {
	Deps     deps.Deps
	Outbound outbound.OutboundItf
}
