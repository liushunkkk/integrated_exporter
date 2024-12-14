package metricx

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

// DuplicateMetric is the error returned by IRegistry. Register when a metric
// already exists. If you mean to Register that metric you must first
// Unregister the existing metric.
type DuplicateMetric string

func (err DuplicateMetric) Error() string {
	return fmt.Sprintf("duplicate metric: %s", string(err))
}

type IRegistry struct {
	Registry *prometheus.Registry
	metrics  sync.Map
}

var DefaultIRegistry = NewIRegistry()

func NewIRegistry() *IRegistry {
	return &IRegistry{
		Registry: prometheus.NewRegistry(),
		metrics:  sync.Map{},
	}
}

// GetOrRegister gets an existing metric or creates and registers a new one.
func (ir *IRegistry) GetOrRegister(im IMetric) any {
	if ret, loaded := ir.metrics.LoadOrStore(im.GetKey(), im); loaded {
		return ret
	}
	if err := ir.Registry.Register(im.GetCollector()); err != nil {
		return nil
	}
	return im
}

// Register the given metric under the given key. Returns a DuplicateMetric
// if a metric by the given name is already registered.
func (ir *IRegistry) Register(im IMetric) error {
	if _, loaded := ir.metrics.Load(im.GetKey()); loaded {
		return DuplicateMetric(im.GetKey())
	}
	if err := ir.Registry.Register(im.GetCollector()); err != nil {
		return err
	}
	return nil
}

// Unregister the metric with the given key.
func (ir *IRegistry) Unregister(key string) {
	if v, loaded := ir.metrics.LoadAndDelete(key); loaded {
		ir.Registry.Unregister(v.(IMetric).GetCollector())
	}
}

// ExportMetrics export the metrics text to buffer.
func (ir *IRegistry) ExportMetrics() (*bytes.Buffer, error) {
	gather, err := ir.Registry.Gather()
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

// ExportDefaultIRegistryMetrics export the metrics text to buffer.
func ExportDefaultIRegistryMetrics() (*bytes.Buffer, error) {
	gather, err := DefaultIRegistry.Registry.Gather()
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
