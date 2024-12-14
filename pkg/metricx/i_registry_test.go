package metricx

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestIRegistry(t *testing.T) {
	DefaultIRegistry = NewIRegistry()
	counter := NewICounter("aaa", prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aaa",
		Help: "aaa",
	}))
	c := DefaultIRegistry.GetOrRegister(counter)
	assert.Equal(t, c, counter)

	err := DefaultIRegistry.Register(counter)
	assert.NotNil(t, err)

	c1 := DefaultIRegistry.GetOrRegister(counter)
	assert.Equal(t, c1, counter)
	assert.Equal(t, c1, c)

	metrics, err := DefaultIRegistry.ExportMetrics()
	assert.Nil(t, err)
	assert.Contains(t, metrics.String(), "aaa 0")

	c1.(*ICounter).Inc()
	metrics, err = DefaultIRegistry.ExportMetrics()
	assert.Nil(t, err)
	assert.Contains(t, metrics.String(), "aaa 1")

	counter2 := NewICounter("bbb", prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bbb",
		Help: "bbb",
	}))

	c2 := DefaultIRegistry.GetOrRegister(counter2)
	assert.Equal(t, c2, counter2)
	counter2.Inc()
	c2.(*ICounter).Inc()

	metrics, err = DefaultIRegistry.ExportMetrics()
	assert.Nil(t, err)
	assert.Contains(t, metrics.String(), "aaa 1")
	assert.Contains(t, metrics.String(), "bbb 2")

	DefaultIRegistry.Unregister(counter2.GetKey())

	metrics, err = DefaultIRegistry.ExportMetrics()
	assert.Nil(t, err)
	assert.Contains(t, metrics.String(), "aaa 1")
	assert.NotContains(t, metrics.String(), "bbb 2")

	DefaultIRegistry.Unregister(counter.GetKey())

	metrics, err = DefaultIRegistry.ExportMetrics()
	assert.Nil(t, err)
	assert.Empty(t, metrics)
}
