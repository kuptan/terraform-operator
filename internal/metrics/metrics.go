package metrics

import (
	"time"

	"github.com/kuptan/terraform-operator/api/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
)

// RecorderInterface is an interface that holds the functions used by the recorder struct
type RecorderInterface interface {
	RecordTotal(name string, namespace string)
	RecordStatus(name string, namespace string, status v1alpha1.TerraformRunStatus)
	RecordDuration(name string, namespace string, start time.Time)
	Collectors() []prometheus.Collector
}

// Recorder is a struct for recording GitOps Toolkit metrics for a controller.
//
// Use NewRecorder to initialise it with properly configured metric names.
type Recorder struct {
	totalCount        *prometheus.CounterVec
	statusGauge       *prometheus.GaugeVec
	durationHistogram *prometheus.HistogramVec
}

// NewRecorder returns a new Recorder with all metric names configured confirm GitOps Toolkit standards.
func NewRecorder() RecorderInterface {
	return &Recorder{
		totalCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "tfo_workflow_total",
				Help: "The total number of submitted workflows",
			},
			[]string{"name", "namespace"},
		),
		statusGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "tfo_workflow_status",
				Help: "The current status of a Terraform workflow/run resource reconciliation.",
			},
			[]string{"name", "namespace"},
		),
		durationHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "tfo_workflow_duration_seconds",
				Help:    "The duration in seconds of a Terraform workflow/run.",
				Buckets: prometheus.ExponentialBuckets(10e-9, 10, 10),
			},
			[]string{"name", "namespace"},
		),
	}
}

// Collectors returns a slice of Prometheus collectors, which can be used to register them in a metrics registry.
func (r *Recorder) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		r.totalCount,
		r.statusGauge,
		r.durationHistogram,
	}
}

// RecordTotal records the total number of submitted workflows
func (r *Recorder) RecordTotal(name string, namespace string) {
	r.totalCount.WithLabelValues(name, namespace).Inc()
}

// RecordStatus records the status for a given terraform workflow/run
func (r *Recorder) RecordStatus(name string, namespace string, status v1alpha1.TerraformRunStatus) {
	var value float64

	if status == v1alpha1.RunWaitingForDependency {
		value = -1
	}

	if status == v1alpha1.RunFailed {
		value = 1
	}

	r.statusGauge.WithLabelValues(name, namespace).Set(value)
}

// RecordDuration records the duration since start for the given ref.
func (r *Recorder) RecordDuration(name string, namespace string, start time.Time) {
	r.durationHistogram.WithLabelValues(name, namespace).Observe(time.Since(start).Seconds())
}
