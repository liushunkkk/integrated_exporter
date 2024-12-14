package metricx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestICounter(t *testing.T) {
	counter := GetOrRegisterICounter(&IOpts{
		Namespace: "aaa",
		Name:      "test_counter",
		Labels:    map[string]string{"instance": "localhost"},
	}, nil)

	counter.Inc()

	metrics, err := ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_counter{instance="localhost"} 1`)

	counter.Inc()

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_counter{instance="localhost"} 2`)

	counter.Add(10)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_counter{instance="localhost"} 12`)

	originCounter := GetOrRegisterICounter(&IOpts{
		Namespace: "aaa",
		Name:      "test_counter",
		Labels:    map[string]string{"instance": "localhost"},
	}, nil)
	originCounter.Inc()

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)
	assert.Contains(t, metrics.String(), `aaa_test_counter{instance="localhost"} 13`)

	nilCounter := GetOrRegisterICounter(nil, nil)

	assert.Nil(t, nilCounter)
}
