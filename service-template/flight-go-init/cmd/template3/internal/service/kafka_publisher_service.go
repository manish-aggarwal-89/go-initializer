package service

import (
	"context"
	"net/http"
	"time"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	commonModel "github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"

	common "github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/kafka"
)

type (
	KafkaPublisherServiceItf interface {
		Publish(mr commonModel.MandatoryRequest, topic string, msg string, isKafkaGeneral bool) error
	}

	KafkaPublisherServiceImpl struct {
		deps deps.Deps
	}
)

func NewKafkaPublisherService(deps deps.Deps) KafkaPublisherServiceItf {
	return &KafkaPublisherServiceImpl{
		deps: deps,
	}
}

func (svc *KafkaPublisherServiceImpl) Publish(mr commonModel.MandatoryRequest, topic string, msg string, isKafkaGeneral bool) error {
	errorCode := common.Success
	startTime := time.Now()

	defer func() {
		go svc.logKafkaSendMetrics(errorCode, mr, topic, startTime)
	}()

	var err error
	if isKafkaGeneral {
		err = svc.deps.Kafka.KafkaGeneral.Publish(topic, msg, []kafka.Header{})
	} else {
		err = svc.deps.Kafka.KafkaFlight.Publish(topic, msg, []kafka.Header{})
	}
	ctx := logrus.InjectRequestMetadataToContext(context.Background(), mr.RequestId, mr.Username, "")
	if err != nil {
		errorCode = common.Failed
		svc.deps.GetLogger(ctx).Error(util.LogService(mr, "search_publish", shared.ERROR_PUBLISH, err.Error()))

		return err
	}

	return err
}

func (svc *KafkaPublisherServiceImpl) logKafkaSendMetrics(errorCode common.ErrorCode, mr commonModel.MandatoryRequest,
	topic string, startTime time.Time) {
	ctx := logrus.InjectRequestMetadataToContext(context.Background(), mr.RequestId, mr.Username, "")
	tags := map[string]interface{}{
		"CHANNEL_NAME": mr.ChannelId,
	}
	statusCode := http.StatusOK
	if errorCode == common.Failed {
		statusCode = http.StatusInternalServerError
	}
	if err := svc.deps.Metric.CustomMonitorLatency(topic, common.KAFKA_SEND, errorCode, statusCode, tags, time.Since(startTime)); err != nil {
		svc.deps.GetLogger(ctx).Error(util.LogOutbound(mr, "logMetricKafkaSendLatency", "", err.Error()))
	}
	if err := svc.deps.Metric.CustomMonitorCounter(topic, common.KAFKA_SEND, errorCode, statusCode, tags); err != nil {
		svc.deps.GetLogger(ctx).Error(util.LogOutbound(mr, "logMetricKafkaSendCounter", "", err.Error()))
	}
}
