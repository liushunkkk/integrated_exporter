package metricx

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ICounter struct {
	key     string
	counter prometheus.Counter
}

func NewICounter(key string, counter prometheus.Counter) *ICounter {
	return &ICounter{
		key:     key,
		counter: counter,
	}
}

// GetOrRegisterICounter gets an existing ICounter or creates and registers a new one.
func GetOrRegisterICounter(opts *IOpts, ir *IRegistry) *ICounter {
	if opts == nil {
		return nil
	}
	if ir == nil {
		ir = DefaultIRegistry
	}
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   opts.Namespace,
		Name:        opts.Name,
		Help:        opts.Help,
		ConstLabels: opts.Labels,
	})
	key := counter.Desc().String()
	ic := NewICounter(key, counter)
	return ir.GetOrRegister(ic).(*ICounter)
}

func (i *ICounter) Inc() {
	i.counter.Inc()
}

func (i *ICounter) Add(v float64) {
	i.counter.Add(v)
}

func (i *ICounter) GetKey() string {
	return i.key
}

func (i *ICounter) GetCollector() prometheus.Collector {
	return i.counter
}
