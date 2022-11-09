package chargebacks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/kafka"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/metrics"
	"github.com/conekta/risk-rules/pkg/text"
)

const (
	handlerName = "chargeback.handler"
)

type ChargebackHandler interface {
	ListenChargebacks()
}
type chargebackHandler struct {
	config  config.Config
	logs    logs.Logger
	service ChargebackService
	metrics datadog.Metricer
}

func NewChargebackHandler(service ChargebackService, cfg config.Config, logger logs.Logger,
	metric datadog.Metricer) ChargebackHandler {
	return &chargebackHandler{
		config:  cfg,
		logs:    logger,
		service: service,
		metrics: metric,
	}
}

func (handler *chargebackHandler) ListenChargebacks() {
	ctx := context.Background()

	consumer, err := kafka.NewFactoryConsumer(handler.logs, handler.config.EventBus.Chargebacks.BoostrapServers,
		kafka.SetSaslAuth(handler.config.EventBus.Chargebacks.EnabledAuth),
		kafka.SetSaslPassword(handler.config.EventBus.Chargebacks.Password),
		kafka.SetSaslUserName(handler.config.EventBus.Chargebacks.User),
		kafka.SetServiceName(handler.config.ProjectName),
		kafka.SetSaslMechanism(handler.config.EventBus.Chargebacks.Mechanism),
		kafka.SetSecurityProtocol(handler.config.EventBus.Chargebacks.SecurityProtocol),
		kafka.SetEnableSslCertificateVerification(handler.config.EventBus.Chargebacks.EnabledSslCertification),
	)
	if err != nil {
		handler.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", handlerName, "NewFactoryConsumer"))
		return
	}

	err = consumer.Listen(ctx, handler.config.EventBus.Chargebacks.Topic,
		handler.config.EventBus.Chargebacks.GroupID, handler.readChargebacks)
	if err != nil {
		handler.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", handlerName, "Listen"))
		return
	}
}

func (handler *chargebackHandler) readChargebacks(ctx context.Context, message []byte) {
	handler.logs.Info(ctx, fmt.Sprintf("reading message [%s]", message), text.LogTagMethod,
		fmt.Sprintf("%s.%s", handlerName, "readChargebacks"))

	defer ctx.Done()

	var chargebackRequest entities.ChargebackRequest
	err := json.Unmarshal(message, &chargebackRequest)
	if err != nil {
		handler.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", handlerName, "readChargebacks"))
		handler.sendChargebackMetrics(ctx, chargebackRequest, false)
		return
	}

	err = handler.service.Save(ctx, chargebackRequest.NewPayerFromPostRequest())
	if err != nil {
		handler.logs.Error(ctx, err.Error(), text.LogTagMethod, fmt.Sprintf("%s.%s", handlerName, "readChargebacks"))
		handler.sendChargebackMetrics(ctx, chargebackRequest, false)
		return
	}
	handler.sendChargebackMetrics(ctx, chargebackRequest, true)
}

func (handler *chargebackHandler) sendChargebackMetrics(ctx context.Context,
	chargebackRequest entities.ChargebackRequest, result bool) {
	metricData := metrics.NewMetricData(ctx, "readChargebacks", handlerName, handler.config.Env)
	metricData.AddCustomTags([]string{
		fmt.Sprintf(text.MetricStatus, chargebackRequest.Status),
	})
	metricData.SetResult(result)
	metrics.SendAsyncMetrics(handler.metrics, handler.logs, metricData, text.SaveChargebackMetricName)
}
