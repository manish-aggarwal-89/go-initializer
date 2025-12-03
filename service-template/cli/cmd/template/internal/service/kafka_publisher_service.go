package service

import (
	"context"
	"net/http"
	"time"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/deps"
	"{{MODULE_NAME}}/internal/shared/util"

	"github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO/log/logrus"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/messaging"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

type (
	KafkaPublisherServiceItf interface {
		Publish(mr common.MandatoryRequest, topic string, msg string, isKafkaGeneral bool) error
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

func (svc *KafkaPublisherServiceImpl) Publish(mr common.MandatoryRequest, topic string, msg string, isKafkaGeneral bool) error {
	errorCode := metrics.Success
	startTime := time.Now()

	defer func() {
		go svc.logKafkaSendMetrics(errorCode, mr, topic, startTime)
	}()

	var err error
	if isKafkaGeneral {
		err = svc.deps.Kafka.KafkaGeneral.Publish(topic, msg, []messaging.Header{})
	} else {
		err = svc.deps.Kafka.KafkaFlight.Publish(topic, msg, []messaging.Header{})
	}
	ctx := logrus.InjectRequestMetadataToContext(context.Background(), mr.RequestId, mr.Username, "")
	if err != nil {
		errorCode = metrics.Failed
		svc.deps.GetLogger(ctx).Error(util.LogService(mr, "search_publish", shared.ERROR_PUBLISH, err.Error()))

		return err
	}

	return err
}

func (svc *KafkaPublisherServiceImpl) logKafkaSendMetrics(errorCode metrics.ErrorCode, mr common.MandatoryRequest,
	topic string, startTime time.Time) {
	ctx := logrus.InjectRequestMetadataToContext(context.Background(), mr.RequestId, mr.Username, "")
	tags := map[string]interface{}{
		"CHANNEL_NAME": mr.ChannelId,
	}
	statusCode := http.StatusOK
	if errorCode == metrics.Failed {
		statusCode = http.StatusInternalServerError
	}
	if err := svc.deps.Metric.CustomMonitorLatency(topic, metrics.KAFKA_SEND, errorCode, statusCode, tags, time.Since(startTime)); err != nil {
		svc.deps.GetLogger(ctx).Errorf(util.LogOutbound(mr, "logMetricKafkaSendLatency", "", err.Error()))
	}
	if err := svc.deps.Metric.CustomMonitorCounter(topic, metrics.KAFKA_SEND, errorCode, statusCode, tags); err != nil {
		svc.deps.GetLogger(ctx).Errorf(util.LogOutbound(mr, "logMetricKafkaSendCounter", "", err.Error()))
	}
}
