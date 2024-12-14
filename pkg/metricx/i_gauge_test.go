package metricx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIGauge(t *testing.T) {
	gauge := GetOrRegisterIGauge(&IOpts{
		Namespace: "aaa",
		Name:      "test_gauge",
		Labels:    map[string]string{"instance": "localhost"},
	}, nil)

	gauge.Inc()

	metrics, err := ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 1`)

	gauge.Inc()

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 2`)

	gauge.Add(10)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 12`)

	originGauge := GetOrRegisterIGauge(&IOpts{
		Namespace: "aaa",
		Name:      "test_gauge",
		Labels:    map[string]string{"instance": "localhost"},
	}, nil)

	originGauge.Inc()

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 13`)

	originGauge.Dec()

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 12`)

	originGauge.Sub(5)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 7`)

	originGauge.Set(100)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_gauge{instance="localhost"} 100`)

	nilGauge := GetOrRegisterIGauge(nil, nil)

	assert.Nil(t, nilGauge)
}
