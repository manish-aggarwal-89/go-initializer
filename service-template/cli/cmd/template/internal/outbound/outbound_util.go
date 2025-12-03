package outbound

import (
	"context"
	"encoding/xml"
	"fmt"
	commonbau "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"io"
	"net/http"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/bau"
	"{{MODULE_NAME}}/internal/shared/util"
)

func (o *OutboundImpl) logData(ctx context.Context, mr common.MandatoryRequest, client Client, url string, marshalledReq []byte) {

	msg := fmt.Sprintf("url:%s - request:%s", url, marshalledReq)
	o.dep.GetLogger(ctx).Info(util.LogOutbound(mr, client.opName, shared.REQUEST, msg))
}

func handleResponse[V any](o *OutboundImpl, ctx context.Context, mr common.MandatoryRequest, client Client, httpResp *http.Response) (*V, error) {

	// try to always unmarshal the response because the response contains error message that might helpful.
	response, err := unmarshalResponse[V](o, ctx, mr, client.opName, httpResp, client)
	if err != nil {
		o.dep.GetLogger(ctx).Error(util.LogOutbound(mr, client.opName, shared.ERROR_UNMARSHAL, fmt.Sprintf("%s: with http status code: %d", err.Error(), httpResp.StatusCode)))
		return nil, err
	}

	if response == nil || isHttpError(httpResp.StatusCode) {
		err = fmt.Errorf("outbound %s failed, status code: %d, status: %s. with response: %+v", client.opName, httpResp.StatusCode, httpResp.Status, response)
		o.dep.GetLogger(ctx).Error(util.LogOutbound(mr, client.opName, shared.ERROR_REQUEST_API, err.Error()))
		return nil, err
	}

	return response, nil
}

func isHttpError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}

func unmarshalResponse[V any](o *OutboundImpl, ctx context.Context, mr common.MandatoryRequest, opName string, resp *http.Response, client Client) (*V, error) {
	var res V
	bodyRes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			o.dep.GetLogger(ctx).Error(util.LogOutbound(mr, client.opName, shared.ERROR_UNMARSHAL, fmt.Sprintf("error closing resp body: %s", err.Error())))
		}
	}(resp.Body)

	o.dep.GetLogger(ctx).Info(util.LogOutbound(mr, opName, shared.RESPONSE, string(bodyRes)))
	if err = xml.Unmarshal(bodyRes, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func createAnalyticSegment(ctx context.Context, name string, isParallel bool) *commonbau.AnalyticSegment {
	v := ctx.Value(shared.CTX_ANALYTIC_TOKEN_KEY)

	if token, ok := v.(bau.AnalyticsTokenModifier); ok {
		analyticSegment := token.CreateSegment(name, isParallel)
		return analyticSegment
	}

	return nil
}

func completeAnalyticsSegment(analyticSegment *commonbau.AnalyticSegment) {
	if analyticSegment != nil {
		analyticSegment.Complete()
	}
}
