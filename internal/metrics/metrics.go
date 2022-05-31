package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionDeleted = "Deleted"
)

// Recorder is a struct for recording GitOps Toolkit metrics for a controller.
//
// Use NewRecorder to initialise it with properly configured metric names.
type Recorder struct {
	conditionGauge    *prometheus.GaugeVec
	suspendGauge      *prometheus.GaugeVec
	durationHistogram *prometheus.HistogramVec
}

// NewRecorder returns a new Recorder with all metric names configured confirm GitOps Toolkit standards.
func NewRecorder() *Recorder {
	return &Recorder{
		conditionGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "tfo_reconcile_condition",
				Help: "The current condition status of a Terraform workflow/run resource reconciliation.",
			},
			[]string{"kind", "name", "namespace", "type", "status"},
		),
		suspendGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "tfo_suspend_status",
				Help: "The current suspend status of a GitOps Toolkit resource.",
			},
			[]string{"kind", "name", "namespace"},
		),
		durationHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "tfo_workflow_duration_seconds",
				Help:    "The duration in seconds of a Terraform workflow/run resource reconciliation.",
				Buckets: prometheus.ExponentialBuckets(10e-9, 10, 10),
			},
			[]string{"kind", "name", "namespace"},
		),
	}
}

// Collectors returns a slice of Prometheus collectors, which can be used to register them in a metrics registry.
func (r *Recorder) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		r.conditionGauge,
		r.suspendGauge,
		r.durationHistogram,
	}
}

// RecordCondition records the condition as given for the ref.
func (r *Recorder) RecordCondition(ref corev1.ObjectReference, condition metav1.Condition, deleted bool) {
	statuses := []string{
		string(metav1.ConditionTrue),
		string(metav1.ConditionFalse),
		string(metav1.ConditionUnknown),
		ConditionDeleted,
	}

	for _, status := range statuses {
		var value float64
		if deleted {
			if status == ConditionDeleted {
				value = 1
			}
		} else {
			if status == string(condition.Status) {
				value = 1
			}
		}
		r.conditionGauge.WithLabelValues(ref.Kind, ref.Name, ref.Namespace, condition.Type, status).Set(value)
	}
}

// RecordSuspend records the suspend status as given for the ref.
func (r *Recorder) RecordSuspend(ref corev1.ObjectReference, suspend bool) {
	var value float64
	if suspend {
		value = 1
	}
	r.suspendGauge.WithLabelValues(ref.Kind, ref.Name, ref.Namespace).Set(value)
}

// RecordDuration records the duration since start for the given ref.
func (r *Recorder) RecordDuration(ref corev1.ObjectReference, start time.Time) {
	r.durationHistogram.WithLabelValues(ref.Kind, ref.Name, ref.Namespace).Observe(time.Since(start).Seconds())
}
