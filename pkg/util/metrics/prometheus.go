package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Counter   = "counter"
	Gauge     = "gauge"
	Summary   = "summary"
	Histogram = "histogram"
)

type Metric struct {
	prometheus.Opts
	Collector prometheus.Collector
	Type      string
	Labels    []string
}

//RegisterMetric use vec when reporting metric with dimension
func RegisterMetric(metric *Metric) {
	switch metric.Type {
	case Counter:
		metric.Collector = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      metric.Name,
			Namespace: metric.Namespace,
			Subsystem: metric.Subsystem,
			Help:      metric.Help,
		}, metric.Labels)
	case Gauge:
		metric.Collector = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      metric.Name,
			Namespace: metric.Namespace,
			Subsystem: metric.Subsystem,
			Help:      metric.Help,
		}, metric.Labels)
	case Summary:
		metric.Collector = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:      metric.Name,
			Namespace: metric.Namespace,
			Subsystem: metric.Subsystem,
			Help:      metric.Help,
		}, metric.Labels)
	case Histogram:
		metric.Collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:      metric.Name,
			Namespace: metric.Namespace,
			Subsystem: metric.Subsystem,
			Help:      metric.Help,
		}, metric.Labels)
	}
	prometheus.MustRegister(metric.Collector)
	return
}
