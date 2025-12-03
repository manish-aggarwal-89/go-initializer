package outbound

import (
	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"
)

type Client struct {
	httpClient util.HTTPRequest
	opName     string
	httpCfg    config.ClientHttpConfig
	baseUrl    string
}

type OutboundImpl struct {
	dep deps.Deps
}

func NewOutboundImpl(dep deps.Deps) *OutboundImpl {
	//httpConfig := dep.Config.OutboundConfig.SupplierOutboundHttpConfig
	//baseURL := httpConfig.BaseUrl
	//createClient := func(cfg config.ClientHttpConfig, opName string) Client {
	//	options := []util.Option{
	//		util.SetStatsd(dep.Metric),
	//	}
	//	return Client{
	//		httpClient: util.NewHTTPRequest(cfg.Retries, cfg.Timeout, options...),
	//		opName:     opName,
	//		httpCfg:    cfg,
	//		baseUrl:    baseURL,
	//	}
	//}
	return &OutboundImpl{
		dep: dep,
	}
}
