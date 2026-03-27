package inbound

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	common "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/kafka"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/system_param"
	kafkaRq "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/request/kafka"
	"{{MODULE_NAME}}/internal/service"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
)

const (
	configKafkaListener = "kafka.search.listener.config"
	ctxOperation        = "listener_search_request"
)

// IntegratorListener handles messages from the integrator search request Kafka topic.
type IntegratorListener interface {
	ListenToIntegratorRequest(data []byte, headers []kafka.Header, _, _ int64) error
}

type integratorListenerImpl struct {
	deps        deps.Deps
	searchProc  service.SearchServiceItf
	systemParam system_param.Service
}

// ListenToIntegratorRequest implements the Kafka callback: check system parameter, then process in worker pool.
func (impl *integratorListenerImpl) ListenToIntegratorRequest(data []byte, headers []kafka.Header, _, _ int64) error {
	ctx := context.Background()
	ctx = logrus.InjectRequestMetadataToContext(ctx, "", "", ctxOperation)

	sysParam, err := impl.systemParam.FindByVariable(ctx, commonModel.MandatoryRequest{}, configKafkaListener)
	if err != nil || sysParam == nil {
		if err != nil {
			impl.deps.GetLogger(ctx).Errorf("failed to get system parameter %s - %s\n", configKafkaListener, err.Error())
		}
		return nil
	}

	if strings.EqualFold(sysParam.Value, "off") {
		return nil
	}

	cfg := impl.deps.Config
	impl.deps.WorkerPoolWrapper.SearchRequestPool.Go(func() {
		var (
			message                = kafkaRq.KafkaIntegratorFareRequest{}
			ctxWithTimeout, cancel = context.WithTimeout(context.Background(), cfg.KafkaConfig.FlightConfig.ReadTimeout)
			errorCode              = common.Success
			startTime              = time.Now()
		)
		tags := map[string]interface{}{"CHANNEL_NAME": ""}

		impl.deps.GetLogger(ctx).Debugf("%s %s %s", shared.INBOUND_IMPL, ctxOperation, util.ObjToJson(message))
		go impl.logKafkaDelayMetric(ctx, headers, tags)
		defer func() {
			cancel()
			go impl.logKafkaConsumerMetrics(ctx, errorCode, tags, startTime)
		}()

		if err := json.Unmarshal(data, &message); err != nil {
			impl.deps.GetLogger(ctx).Errorf("failed to unmarshal json request %s\n", err.Error())
			errorCode = common.Failed
			return
		}
		tags["CHANNEL_NAME"] = message.MandatoryRequest.ChannelId

		if err := impl.searchProc.SearchKafka(ctxWithTimeout, message); err != nil {
			errorCode = common.Failed
			impl.deps.GetLogger(ctx).Errorf("error at searchKafka %v", err)
		}
	})

	return nil
}

func (impl *integratorListenerImpl) logKafkaConsumerMetrics(_ context.Context, errorCode common.ErrorCode, tags map[string]interface{},
	startTime time.Time) {
	go func() {
		statusCode := http.StatusOK
		if errorCode == common.Failed {
			statusCode = http.StatusInternalServerError
		}
		topic := impl.deps.Config.KafkaConfig.Topics.IntegratorSearchRequest
		_ = impl.deps.Metric.CustomMonitorLatency(topic, common.API_IN, errorCode, statusCode, tags, time.Since(startTime))
		_ = impl.deps.Metric.CustomMonitorCounter(topic, common.API_IN, errorCode, statusCode, tags)
	}()
}

func (impl *integratorListenerImpl) logKafkaDelayMetric(ctx context.Context, headers []kafka.Header, tags map[string]interface{}) {
	delay, err := getKafkaDelayDuration(headers)
	if err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogInbound(commonModel.MandatoryRequest{}, "logMetricKafkaDelay", "", err.Error()))
		return
	}
	topic := impl.deps.Config.KafkaConfig.Topics.IntegratorSearchRequest
	if err = impl.deps.Metric.CustomMonitorLatency(topic, common.API_IN, common.Success, http.StatusOK, tags, delay); err != nil {
		impl.deps.GetLogger(ctx).Error(util.LogInbound(commonModel.MandatoryRequest{}, "logMetricKafkaDelay", "", err.Error()))
	}
}

func getKafkaDelayDuration(headers []kafka.Header) (time.Duration, error) {
	for _, h := range headers {
		if strings.EqualFold(h.Key, "timestamp") {
			nano, err := strconv.ParseInt(string(h.Value), 10, 64)
			if err != nil {
				return 0, err
			}
			return time.Since(time.Unix(0, nano)), nil
		}
	}
	return 0, errors.New("no timestamp Kafka header found")
}

// NewIntegratorListener builds the search integrator listener (deps, search processor, system parameter service from common lib).
func NewIntegratorListener(deps deps.Deps, searchProc service.SearchServiceItf, systemParam system_param.Service) (IntegratorListener, error) {
	return &integratorListenerImpl{deps: deps, searchProc: searchProc, systemParam: systemParam}, nil
}
