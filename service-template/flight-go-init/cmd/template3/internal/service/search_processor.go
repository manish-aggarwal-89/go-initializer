package service

import (
	"context"
	"fmt"
	"strings"

	"{{MODULE_NAME}}/config"
	"{{MODULE_NAME}}/internal/service/rule_supplier"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/helper"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/credential"
	commonUtil "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/shared/util"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/constant/enum"
	fareRQ "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/fare"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/kafka"
	fareRS "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/response/fare"
)

type SearchProcessorService interface {
	SearchProcessor(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType) (*fareRS.FlightIntegratorSearchResponse, error)
	SearchAndPublishProcessor(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType) (*fareRS.FlightIntegratorSearchResponse, error)
	SearchKafkaProcessor(ctx context.Context, msg kafka.KafkaIntegratorFareRequest) error
	PublishSearchResponse(mr common.MandatoryRequest, res *fareRS.FlightIntegratorSearchResponse, err error)
}

type searchProcessorImpl struct {
	deps               deps.Deps
	config             *config.Config
	credentialService  credential.Service
	searchAnalyticsSvc AnalyticsService[fareRQ.IntegratorFareRequest]
	publisherService   KafkaPublisherServiceItf
	searchRegistry     *rule_supplier.SearchRuleRegistry
}

func NewSearchProcessorService(
	deps deps.Deps,
	cfg *config.Config,
	credentialService credential.Service,
	searchAnalyticsSvc AnalyticsService[fareRQ.IntegratorFareRequest],
	publisherService KafkaPublisherServiceItf,
	searchRegistry *rule_supplier.SearchRuleRegistry,
) (SearchProcessorService, error) {
	return &searchProcessorImpl{
		deps:               deps,
		config:             cfg,
		credentialService:  credentialService,
		searchAnalyticsSvc: searchAnalyticsSvc,
		publisherService:   publisherService,
		searchRegistry:     searchRegistry,
	}, nil
}

func (p *searchProcessorImpl) SearchProcessor(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType) (*fareRS.FlightIntegratorSearchResponse, error) {
	var err error

	analyticToken := p.searchAnalyticsSvc.CreateToken(mr, req)
	ctx = context.WithValue(ctx, shared.CTX_ANALYTIC_TOKEN_KEY, analyticToken)

	defer func() {
		go p.searchAnalyticsSvc.Complete(ctx, mr, analyticToken, err)
	}()

	distType := p.config.DistributionType
	if distType == "" {
		err = fmt.Errorf("config.DistributionType is required but not set")
		return nil, err
	}
	rule, ok := p.searchRegistry.Get(distType)
	if !ok {
		err = fmt.Errorf("search rule not found for distribution type: %s", distType)
		return nil, err
	}
	if len(req.SupplierMappings) == 0 || req.SupplierMappings[0].Account.Code == "" {
		err = fmt.Errorf("supplier code is empty")
		return nil, err
	}
	supplierCode := req.SupplierMappings[0].Account.Code
	cred, err := p.credentialService.FindBySupplierInCache(ctx, mr, supplierCode)
	if err != nil || cred == nil {
		err = fmt.Errorf("credential not found for supplier %s: %w", supplierCode, err)
		return nil, err
	}
	return rule.Search(ctx, mr, req, tripType, cred, "")
}

func (p *searchProcessorImpl) SearchAndPublishProcessor(ctx context.Context, mr common.MandatoryRequest, req fareRQ.IntegratorFareRequest, tripType enum.TripType) (*fareRS.FlightIntegratorSearchResponse, error) {
	var resp *fareRS.FlightIntegratorSearchResponse
	var err error

	defer func() {
		go p.PublishSearchResponse(mr, resp, err)
	}()

	resp, err = p.SearchProcessor(ctx, mr, req, tripType)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *searchProcessorImpl) SearchKafkaProcessor(ctx context.Context, req kafka.KafkaIntegratorFareRequest) error {
	op := "search_kafka"

	distType := p.config.DistributionType
	rule, ok := p.searchRegistry.Get(distType)
	if !ok {
		return fmt.Errorf("search rule not found for distribution type: %s", distType)
	}

	for _, fareReq := range req.IntegratorFareRequests {
		if !strings.EqualFold(fareReq.DistributionType, shared.DISTRIBUTION_TYPE) {
			continue
		}
		logData := commonUtil.GetDataToLogFromSearch(fareReq)
		ctx = commonUtil.ConvertAndMapToSDCContext(ctx, logData)
		p.deps.GetLogger(ctx).Debug(util.LogService(req.MandatoryRequest, op, shared.REQUEST, fareReq))
		for _, tripType := range fareReq.TripTypes {
			for _, supplierMapping := range fareReq.SupplierMappings {
				supplierCode := supplierMapping.Account.Code
				cred, credErr := p.credentialService.FindBySupplierInCache(ctx, req.MandatoryRequest, supplierCode)
				if credErr != nil || cred == nil {
					p.deps.GetLogger(ctx).Error(util.LogService(req.MandatoryRequest, op, shared.PROCESS, fmt.Sprintf("credential not found for supplier %s: %v", supplierCode, credErr)))
					continue
				}
				if _, err := rule.Search(ctx, req.MandatoryRequest, fareReq, tripType, cred, ""); err != nil {
					p.deps.GetLogger(ctx).Error(util.LogService(req.MandatoryRequest, op, shared.PROCESS, err.Error()))
				}
			}
		}
	}
	return nil
}

func (p *searchProcessorImpl) PublishSearchResponse(mr common.MandatoryRequest, res *fareRS.FlightIntegratorSearchResponse, err error) {
	var response fareRS.FlightIntegratorSearchResponse
	if res != nil {
		response = *res
	}
	publishErr := p.publisherService.Publish(mr, p.deps.Config.KafkaConfig.Topics.IntegratorSearchResponse, util.ObjToJson(helper.ConstructKafkaMessageResponseSearch(mr, response, err)), false)
	if publishErr != nil {
		p.deps.Logger.Error("publish search response failed, err: ", publishErr.Error())
	}
}
