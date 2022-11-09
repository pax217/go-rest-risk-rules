package datadog

import (
	"context"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/stretchr/testify/mock"
)

type MetricsDogMock struct {
	mock.Mock
}

func (m *MetricsDogMock) Incr(ctx context.Context, name string, tags []string, rate float64) error {
	return nil
}
func (m *MetricsDogMock) Decr(ctx context.Context, name string, tags []string, rate float64) error {
	return nil
}
func (m *MetricsDogMock) Count(ctx context.Context, name string, value int64, tags []string, rate float64) error {
	return nil
}
func (m *MetricsDogMock) Gauge(ctx context.Context, name string, value float64, tags []string, rate float64) error {
	return nil
}
func (m *MetricsDogMock) Histogram(ctx context.Context, name string, value float64, tags []string, rate float64) error {
	return nil
}
func (m *MetricsDogMock) Distribution(ctx context.Context, name string, value float64, tags []string, rate float64) error {
	return nil
}
func (m *MetricsDogMock) Client() *statsd.Client {
	return nil
}
