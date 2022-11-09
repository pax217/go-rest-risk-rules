package metrics

import (
	"context"
	"fmt"

	"github.com/conekta/go_common/datadog"
	"github.com/conekta/go_common/logs"
	"github.com/conekta/risk-rules/pkg/text"
)

type MetricData struct {
	context context.Context
	method  string
	service string

	env        string
	processed  bool
	customTags []string
}

func NewMetricData(ctx context.Context, methodName string,
	serviceName string, environment string) MetricData {
	return MetricData{
		context:    ctx,
		method:     methodName,
		service:    serviceName,
		env:        environment,
		customTags: make([]string, 0),
	}
}

func (metric *MetricData) SetResult(isProcessed bool) {
	metric.processed = isProcessed
}

func (metric *MetricData) AddCustomTags(tags []string) {
	metric.customTags = append(metric.customTags, tags...)
}

func SendAsyncMetrics(dataDog datadog.Metricer, logger logs.Logger, data MetricData, metricName string) {
	go func() {
		data.customTags = append(data.customTags,
			fmt.Sprintf(text.MetricTagSuccess, data.processed),
			fmt.Sprintf(text.MetricTagScope, data.env),
		)
		err := dataDog.Incr(data.context, metricName, data.customTags, 1)
		if err != nil {
			logger.Error(data.context, err.Error(), text.LogTagMethod, fmt.Sprintf(data.service, data.method))
		}
	}()
}
