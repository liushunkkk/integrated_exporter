package metricx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIHistogram(t *testing.T) {
	opts := &IHistogramOpts{
		IOpts: &IOpts{
			Namespace: "aaa",
			Name:      "test_histogram",
			Labels:    map[string]string{"instance": "localhost"},
		},
		Buckets: []float64{10, 20, 30, 40},
	}

	histogram := GetOrRegisterIHistogram(opts, nil)

	histogram.Observe(15)

	metrics, err := ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)

	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="10"} 0`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="20"} 1`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="30"} 1`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="40"} 1`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="+Inf"} 1`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_sum{instance="localhost"} 15`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_count{instance="localhost"} 1`)

	histogram.Observe(25)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)

	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="10"} 0`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="20"} 1`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="30"} 2`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="40"} 2`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="+Inf"} 2`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_sum{instance="localhost"} 40`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_count{instance="localhost"} 2`)

	histogram = GetOrRegisterIHistogram(opts, nil)

	histogram.Observe(80)
	histogram.Observe(100)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)

	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="10"} 0`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="20"} 1`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="30"} 2`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="40"} 2`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_bucket{instance="localhost",le="+Inf"} 4`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_sum{instance="localhost"} 220`)
	assert.Contains(t, metrics.String(), `aaa_test_histogram_count{instance="localhost"} 4`)

	histogram = GetOrRegisterIHistogram(nil, nil)
	assert.Nil(t, histogram)
}
