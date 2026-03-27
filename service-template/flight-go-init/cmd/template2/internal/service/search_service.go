package service

import (
	"context"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	commonUtil "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/util"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/kafka"
	fareRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
	"strings"
	"{{MODULE_NAME}}/internal/outbound"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/helper"
	"{{MODULE_NAME}}/internal/shared/util"
)

type SearchServiceImpl struct {
	deps               deps.Deps
	searchAnalyticsSvc AnalyticsService[fareRQ.IntegratorFareRequest]
	publisherService   KafkaPublisherServiceItf
	credentialService  credential.Service
	outbound           outbound.OutboundItf
}

type SearchServiceItf interface {
	SearchWithMetric(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType, supplierMapping fareRQ.SupplierMapping) (fareRS.FlightIntegratorSearchResponse, error)
	SearchAndPublish(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType, supplierMapping fareRQ.SupplierMapping) (fareRS.FlightIntegratorSearchResponse, error)
	SearchKafka(ctx context.Context, msg kafka.KafkaIntegratorFareRequest) error
	PublishSearchResponse(mr common.MandatoryRequest, res fareRS.FlightIntegratorSearchResponse, err error)
}

func NewSearchServiceImpl(deps deps.Deps,
	searchAnalyticsSvc AnalyticsService[fareRQ.IntegratorFareRequest],
	publisherService KafkaPublisherServiceItf,
	credentialService credential.Service,
	outbound outbound.OutboundItf,
) *SearchServiceImpl {
	return &SearchServiceImpl{
		deps:               deps,
		searchAnalyticsSvc: searchAnalyticsSvc,
		publisherService:   publisherService,
		credentialService:  credentialService,
		outbound:           outbound,
	}
}

func (impl *SearchServiceImpl) SearchWithMetric(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType, supplierMapping fareRQ.SupplierMapping) (fareRS.FlightIntegratorSearchResponse, error) {
	var err error

	analyticToken := impl.searchAnalyticsSvc.CreateToken(mr, req)
	ctx = context.WithValue(ctx, shared.CTX_ANALYTIC_TOKEN_KEY, analyticToken)

	defer func() {
		go impl.searchAnalyticsSvc.Complete(ctx, mr, analyticToken, err)
	}()

	resp, err := impl.search(ctx, mr, req, tripType, supplierMapping)
	if err != nil {
		return fareRS.FlightIntegratorSearchResponse{}, err
	}

	return resp, nil
}

func (impl *SearchServiceImpl) search(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType, supplierMapping fareRQ.SupplierMapping) (fareRS.FlightIntegratorSearchResponse, error) {
	return fareRS.FlightIntegratorSearchResponse{}, nil
}

func (impl *SearchServiceImpl) SearchAndPublish(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType, supplierMapping fareRQ.SupplierMapping) (fareRS.FlightIntegratorSearchResponse, error) {
	var resp fareRS.FlightIntegratorSearchResponse
	var err error

	defer func() {
		go impl.PublishSearchResponse(mr, resp, err)
	}()

	resp, err = impl.SearchWithMetric(ctx, mr, req, tripType, supplierMapping)
	if err != nil {
		return fareRS.FlightIntegratorSearchResponse{}, err
	}

	return resp, nil
}

func (impl *SearchServiceImpl) SearchKafka(ctx context.Context, req kafka.KafkaIntegratorFareRequest) error {
	op := "search_kafka"

	for _, fareReq := range req.IntegratorFareRequests {
		if !strings.EqualFold(fareReq.DistributionType, shared.DISTRIBUTION_TYPE) {
			continue
		}
		logData := commonUtil.GetDataToLogFromSearch(fareReq)
		ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)
		impl.deps.GetLogger(ctx).Debug(util.LogService(req.MandatoryRequest, op, shared.REQUEST, fareReq))
		for _, tripType := range fareReq.TripTypes {
			for _, supplierMapping := range fareReq.SupplierMappings {
				if _, err := impl.search(ctx, req.MandatoryRequest, fareReq, tripType, supplierMapping); err != nil {
					impl.deps.GetLogger(ctx).Error(util.LogService(req.MandatoryRequest, op, shared.PROCESS, err.Error()))
				}
			}
		}
	}
	return nil
}

func (impl *SearchServiceImpl) PublishSearchResponse(mr common.MandatoryRequest, res fareRS.FlightIntegratorSearchResponse, err error) {
	publishErr := impl.publisherService.Publish(mr, impl.deps.Config.KafkaConfig.Topics.IntegratorSearchResponse, util.ObjToJson(helper.ConstructKafkaMessageResponseSearch(mr, res, err)), false)
	if publishErr != nil {
		impl.deps.Logger.Error("publish search response failed, err: ", publishErr.Error())
	}
}
