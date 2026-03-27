package outbound

import (
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/config"
	httpclient "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/httpclient/util"
	"{{MODULE_NAME}}/internal/shared/deps"
)

type OutboundImpl struct {
	dep         deps.Deps
	dummyClient httpclient.BaseClient
}

func NewOutboundImpl(dep deps.Deps) *OutboundImpl {
	httpConfig := dep.Config.OutboundConfig.SupplierOutboundHttpConfig
	baseURL := httpConfig.BaseUrl
	createClient := func(cfg config.ClientHttpConfig, opName string) httpclient.BaseClient {
		options := []httpclient.Option{
			httpclient.SetStatsd(dep.Metric),
		}
		return httpclient.BaseClient{
			HttpClient: httpclient.NewHTTPRequest(cfg.Retries, cfg.Timeout, options...),
			OpName:     opName,
			HttpCfg:    cfg,
			BaseUrl:    baseURL,
		}
	}
	return &OutboundImpl{
		dep:         dep,
		dummyClient: createClient(config.ClientHttpConfig{}, "dummy"), // added as dummy
	}
}
