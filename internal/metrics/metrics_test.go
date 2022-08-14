package metrics

import (
	"time"

	"github.com/kuptan/terraform-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prometheus/client_golang/prometheus"
)

var _ = Describe("Metrics Recorder", func() {
	rec := NewRecorder()

	reg := prometheus.NewRegistry()
	reg.MustRegister(rec.Collectors()...)

	const (
		name      = "terraform-workflow"
		namespace = "default"
	)

	Context("Recording Total", func() {
		It("should record the total count metric", func() {
			rec.RecordTotal(name, namespace)

			var (
				value      float64 = 1.0
				metricName string  = "tfo_workflow_total"
			)

			metricFamilies, err := reg.Gather()

			Expect(err).ToNot(HaveOccurred())
			Expect(metricFamilies).To(HaveLen(1))
			Expect(metricFamilies[0].Name).To(Equal(&metricName))
			Expect(metricFamilies[0].Metric).To(HaveLen(1))
			Expect(metricFamilies[0].Metric[0].Counter).ToNot(BeNil())
			Expect(metricFamilies[0].Metric[0].Counter.Value).To(Equal(&value))
		})
	})

	Context("Recording Status", func() {
		It("should record the waitingForDependency status", func() {
			rec.RecordStatus(name, namespace, v1alpha1.RunWaitingForDependency)

			var (
				value      float64 = -1.0
				metricName string  = "tfo_workflow_status"
			)

			metricFamilies, err := reg.Gather()

			Expect(err).ToNot(HaveOccurred())
			Expect(metricFamilies).To(HaveLen(2))
			Expect(metricFamilies[0].Name).To(Equal(&metricName))
			Expect(metricFamilies[0].Metric).To(HaveLen(1))
			Expect(metricFamilies[0].Metric[0].Gauge).ToNot(BeNil())
			Expect(metricFamilies[0].Metric[0].Gauge.Value).To(Equal(&value))
		})

		It("should record the failed status", func() {
			rec.RecordStatus(name, namespace, v1alpha1.RunFailed)

			var (
				value      float64 = 1.0
				metricName string  = "tfo_workflow_status"
			)

			metricFamilies, err := reg.Gather()

			Expect(err).ToNot(HaveOccurred())
			Expect(metricFamilies).To(HaveLen(2))
			Expect(metricFamilies[0].Name).To(Equal(&metricName))
			Expect(metricFamilies[0].Metric).To(HaveLen(1))
			Expect(metricFamilies[0].Metric[0].Gauge).ToNot(BeNil())
			Expect(metricFamilies[0].Metric[0].Gauge.Value).To(Equal(&value))
		})

		It("should record the completed status", func() {
			rec.RecordStatus(name, namespace, v1alpha1.RunCompleted)

			var (
				value      float64 = 0.0
				metricName string  = "tfo_workflow_status"
			)

			metricFamilies, err := reg.Gather()

			Expect(err).ToNot(HaveOccurred())
			Expect(metricFamilies).To(HaveLen(2))
			Expect(metricFamilies[0].Name).To(Equal(&metricName))
			Expect(metricFamilies[0].Metric).To(HaveLen(1))
			Expect(metricFamilies[0].Metric[0].Gauge).ToNot(BeNil())
			Expect(metricFamilies[0].Metric[0].Gauge.Value).To(Equal(&value))
		})

	})

	Context("Recording Duration", func() {
		It("should record the duration metric", func() {
			rec.RecordDuration(name, namespace, time.Now())

			var (
				metricName string = "tfo_workflow_duration_seconds"
			)

			metricFamilies, err := reg.Gather()

			Expect(err).ToNot(HaveOccurred())
			Expect(metricFamilies).To(HaveLen(3))
			Expect(metricFamilies[0].Name).To(Equal(&metricName))
			Expect(metricFamilies[0].Metric).To(HaveLen(1))
			Expect(metricFamilies[0].Metric[0].Histogram).ToNot(BeNil())
		})
	})
})
