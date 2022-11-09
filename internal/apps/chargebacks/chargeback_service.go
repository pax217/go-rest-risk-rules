package chargebacks

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/internal/config"
	"github.com/conekta/risk-rules/internal/entities"
	"github.com/conekta/risk-rules/pkg/strings"
	"github.com/conekta/risk-rules/pkg/text"
)

const (
	serviceName            = "chargeback.service"
	emailEmptyWarning      = "warning, email is empty"
	existChargebackWarning = "warning, a chargeback already exists with this ID"
)

type ChargebackService interface {
	Save(ctx context.Context, payer entities.Payer) error
}

type chargebackService struct {
	chargebackRepository ChargebackRepository
	log                  logs.Logger
	datadog              datadog.Metricer
	configs              config.Config
}

func NewChargebacksService(cfg config.Config, repository ChargebackRepository, logger logs.Logger,
	datadogMetric datadog.Metricer) ChargebackService {
	return &chargebackService{
		configs:              cfg,
		chargebackRepository: repository,
		log:                  logger,
		datadog:              datadogMetric,
	}
}

func (service *chargebackService) Save(ctx context.Context, payer entities.Payer) error {
	if strings.IsEmpty(payer.Email) {
		service.log.Warn(ctx, emailEmptyWarning, text.LogTagMethod, fmt.Sprintf("%s.%s", serviceName, "Save"),
			text.Email, payer.Email, text.Chargeback, payer)
		return nil
	}

	payerFound, err := service.chargebackRepository.Find(ctx, payer)
	if err != nil {
		return err
	}

	if !strings.IsEmpty(payerFound.Email) {
		return service.updateChargeback(ctx, payer, payerFound)
	}

	err = service.chargebackRepository.Save(ctx, payer)
	if err != nil {
		return err
	}

	return nil
}

func (service *chargebackService) updateChargeback(ctx context.Context, payer, payerFound entities.Payer) error {
	if payer.ExistChargeback(payerFound.Chargebacks) {
		service.log.Warn(ctx, existChargebackWarning, text.LogTagMethod, fmt.Sprintf("%s.%s", serviceName, "updateChargeback"),
			text.ChargebackID, payer.Chargebacks[0].ChargebackID, text.Payer, payer)
		return nil
	}

	payerFound.Chargebacks = append(payer.Chargebacks, payerFound.Chargebacks...)

	err := service.chargebackRepository.Update(ctx, payerFound)
	if err != nil {
		return err
	}

	return nil
}
