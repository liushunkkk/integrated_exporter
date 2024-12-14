package metricx

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ISummary struct {
	key     string
	summary prometheus.Summary
}

func NewISummary(key string, summary prometheus.Summary) *ISummary {
	return &ISummary{
		key:     key,
		summary: summary,
	}
}

func NewRecommendISummaryOpts(opts *IOpts) *ISummaryOpts {
	return &ISummaryOpts{
		IOpts: opts,
		Objectives: map[float64]float64{
			0.5:  0.05,  // 中位数，误差 ±5%
			0.75: 0.02,  // 75th 分位数，误差 ±2%
			0.9:  0.01,  // 90th 分位数，误差 ±1%
			0.95: 0.005, // 95th 分位数，误差 ±0.5%
			0.99: 0.001, // 99th 分位数，误差 ±0.1%
		},
		MaxAge:     prometheus.DefMaxAge,
		AgeBuckets: prometheus.DefAgeBuckets,
		BufCap:     prometheus.DefBufCap,
	}
}

// GetOrRegisterISummary gets an existing ISummary or creates and registers a new one.
func GetOrRegisterISummary(opts *ISummaryOpts, ir *IRegistry) *ISummary {
	if opts == nil {
		return nil
	}
	if ir == nil {
		ir = DefaultIRegistry
	}
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   opts.Namespace,
		Name:        opts.Name,
		Help:        opts.Help,
		ConstLabels: opts.Labels,
		Objectives:  opts.Objectives,
		MaxAge:      opts.MaxAge,
		AgeBuckets:  opts.AgeBuckets,
		BufCap:      opts.BufCap,
	})
	key := summary.Desc().String()
	ic := NewISummary(key, summary)
	return ir.GetOrRegister(ic).(*ISummary)
}

func (i *ISummary) Observe(v float64) {
	i.summary.Observe(v)
}

func (i *ISummary) GetKey() string {
	return i.key
}

func (i *ISummary) GetCollector() prometheus.Collector {
	return i.summary
}
