package metricx

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type IOpts struct {
	// Namespace and Name are components of the fully-qualified
	// name of the Metric (created by joining these components with
	// "_").
	Namespace string
	Name      string

	// Help provides information about this metric.
	//
	// Metrics with the same fully-qualified name must have the same Help
	// string.
	Help string

	// Labels are used to attach fixed labels to this metric. Metrics
	// with the same fully-qualified name must have the same label names in
	// their Labels.
	Labels map[string]string
}

type ISummaryOpts struct {
	*IOpts

	// Objectives defines the quantile rank estimates with their respective
	// absolute error. If Objectives[q] = e, then the value reported for q
	// will be the φ-quantile value for some φ between q-e and q+e.  The
	// default value is an empty map, resulting in a summary without
	// quantiles.
	Objectives map[float64]float64

	// MaxAge defines the duration for which an observation stays relevant
	// for the summary. Must be positive. The default value is 10 minutes.
	MaxAge time.Duration

	// AgeBuckets is the number of buckets used to exclude observations that
	// are older than MaxAge from the summary. A higher number has a
	// resource penalty, so only increase it if the higher resolution is
	// really required. For very high observation rates, you might want to
	// reduce the number of age buckets. With only one age bucket, you will
	// effectively see a complete reset of the summary each time MaxAge has
	// passed. The default value is 5.
	AgeBuckets uint32

	// BufCap defines the default sample stream buffer size.  The default
	// value of 500 should suffice for most uses. If there is a need
	// to increase the value, a multiple of 500 is recommended.
	BufCap uint32
}

type IHistogramOpts struct {
	*IOpts

	// Buckets defines the buckets into which observations are counted. Each
	// element in the slice is the upper inclusive bound of a bucket. The
	// values must be sorted in strictly increasing order. There is no need
	// to add the highest bucket with +Inf bound, it will be added
	// implicitly.
	Buckets []float64
}

type IMetric interface {
	GetKey() string
	GetCollector() prometheus.Collector
}
