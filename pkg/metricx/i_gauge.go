package metricx

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IGauge struct {
	key   string
	gauge prometheus.Gauge
}

func NewIGauge(key string, gauge prometheus.Gauge) *IGauge {
	return &IGauge{
		key:   key,
		gauge: gauge,
	}
}

// GetOrRegisterIGauge gets an existing IGauge or creates and registers a new one.
func GetOrRegisterIGauge(opts *IOpts, ir *IRegistry) *IGauge {
	if opts == nil {
		return nil
	}
	if ir == nil {
		ir = DefaultIRegistry
	}
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   opts.Namespace,
		Name:        opts.Name,
		Help:        opts.Help,
		ConstLabels: opts.Labels,
	})
	key := gauge.Desc().String()
	ic := NewIGauge(key, gauge)
	return ir.GetOrRegister(ic).(*IGauge)
}

func (i *IGauge) Inc() {
	i.gauge.Inc()
}

func (i *IGauge) Dec() {
	i.gauge.Dec()
}

func (i *IGauge) Add(v float64) {
	i.gauge.Add(v)
}

func (i *IGauge) Sub(v float64) {
	i.gauge.Sub(v)
}

func (i *IGauge) Set(v float64) {
	i.gauge.Set(v)
}

func (i *IGauge) GetKey() string {
	return i.key
}

func (i *IGauge) GetCollector() prometheus.Collector {
	return i.gauge
}
