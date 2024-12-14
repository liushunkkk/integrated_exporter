package metricx

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IHistogram struct {
	key       string
	histogram prometheus.Histogram
}

func NewIHistogram(key string, histogram prometheus.Histogram) *IHistogram {
	return &IHistogram{
		key:       key,
		histogram: histogram,
	}
}

func NewRecommendIHistogramOpts(opts *IOpts) *IHistogramOpts {
	return &IHistogramOpts{
		opts,
		prometheus.DefBuckets,
	}
}

// GetOrRegisterIHistogram gets an existing IHistogram or creates and registers a new one.
func GetOrRegisterIHistogram(opts *IHistogramOpts, ir *IRegistry) *IHistogram {
	if opts == nil {
		return nil
	}
	if ir == nil {
		ir = DefaultIRegistry
	}
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   opts.Namespace,
		Name:        opts.Name,
		Help:        opts.Help,
		ConstLabels: opts.Labels,
		Buckets:     opts.Buckets,
	})
	key := histogram.Desc().String()
	ic := NewIHistogram(key, histogram)
	return ir.GetOrRegister(ic).(*IHistogram)
}

func (i *IHistogram) Observe(v float64) {
	i.histogram.Observe(v)
}

func (i *IHistogram) GetKey() string {
	return i.key
}

func (i *IHistogram) GetCollector() prometheus.Collector {
	return i.histogram
}
