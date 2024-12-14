package metricx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestISummary(t *testing.T) {
	opts := NewRecommendISummaryOpts(&IOpts{
		Namespace: "aaa",
		Name:      "test_summary",
		Labels:    map[string]string{"instance": "localhost"},
	})
	summary := GetOrRegisterISummary(opts, nil)

	summary.Observe(100)

	metrics, err := ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)

	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.5"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.75"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.9"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.95"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.99"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary_sum{instance="localhost"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary_count{instance="localhost"} 1`)

	summary.Observe(200)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)

	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.5"} 100`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.75"} 200`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.9"} 200`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.95"} 200`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.99"} 200`)
	assert.Contains(t, metrics.String(), `aaa_test_summary_sum{instance="localhost"} 300`)
	assert.Contains(t, metrics.String(), `aaa_test_summary_count{instance="localhost"} 2`)

	summary = GetOrRegisterISummary(opts, nil)

	summary.Observe(300)
	summary.Observe(400)

	metrics, err = ExportDefaultIRegistryMetrics()
	assert.NoError(t, err)

	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.5"} 200`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.75"} 300`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.9"} 400`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.95"} 400`)
	assert.Contains(t, metrics.String(), `aaa_test_summary{instance="localhost",quantile="0.99"} 400`)
	assert.Contains(t, metrics.String(), `aaa_test_summary_sum{instance="localhost"} 1000`)
	assert.Contains(t, metrics.String(), `aaa_test_summary_count{instance="localhost"} 4`)

	summary = GetOrRegisterISummary(nil, nil)
	assert.Nil(t, summary)
}
