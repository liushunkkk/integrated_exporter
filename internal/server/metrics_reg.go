package server

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"sync"
)

var Reg *prometheus.Registry

var metricsMap sync.Map

func NewMetricsRegistry() *prometheus.Registry {
	if Reg == nil {
		Reg = prometheus.NewRegistry()
	}
	metricsMap = sync.Map{}
	return Reg
}

func GetOrRegisterGauge(opts *prometheus.GaugeOpts) prometheus.Gauge {
	g := prometheus.NewGauge(*opts)
	key := g.Desc().String()
	if ret, loaded := metricsMap.LoadOrStore(key, g); loaded {
		return ret.(prometheus.Gauge)
	}
	if err := Reg.Register(g); err != nil {
		return nil
	}
	return g
}

func GetOrRegisterCounter(opts *prometheus.CounterOpts) prometheus.Counter {
	g := prometheus.NewCounter(*opts)
	key := g.Desc().String()
	if ret, loaded := metricsMap.LoadOrStore(key, g); loaded {
		return ret.(prometheus.Counter)
	}
	if err := Reg.Register(g); err != nil {
		return nil
	}
	return g
}

func GetOrRegisterHistogram(opts *prometheus.HistogramOpts) prometheus.Histogram {
	g := prometheus.NewHistogram(*opts)
	key := g.Desc().String()
	if ret, loaded := metricsMap.LoadOrStore(key, g); loaded {
		return ret.(prometheus.Histogram)
	}
	if err := Reg.Register(g); err != nil {
		return nil
	}
	return g
}

func GetOrRegisterSummary(opts *prometheus.SummaryOpts) prometheus.Summary {
	g := prometheus.NewSummary(*opts)
	key := g.Desc().String()
	if ret, loaded := metricsMap.LoadOrStore(key, g); loaded {
		return ret.(prometheus.Summary)
	}
	if err := Reg.Register(g); err != nil {
		return nil
	}
	return g
}

func GetMetricsText() (*bytes.Buffer, error) {
	gather, err := Reg.Gather()
	if err != nil {
		return nil, err
	}
	b := new(bytes.Buffer)
	for _, mf := range gather {
		if _, err = expfmt.MetricFamilyToText(b, mf); err != nil {
			return nil, err
		}
	}
	return b, nil
}
