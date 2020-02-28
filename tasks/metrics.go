package tasks

import (
	"fmt"
	"net/http"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	counters = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "dgraph",
			Subsystem: "people",
			Name:      "query_count",
			Help:      "Total number of queries.",
		}, []string{"query", "status"})

	durations = prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: "dgraph",
			Subsystem: "people",
			Name:      "query_duration",
			Help:      "Histogram of query duration.",
			Buckets:   prom.DefBuckets,
		},
		[]string{"query", "status"},
	)

	throughput = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "dgraph",
			Subsystem: "people",
			Name:      "query_throughput",
			Help:      "Query throughput",
		},
		[]string{"query", "status"},
	)
)

func init() {
	prom.Register(counters)
	prom.Register(durations)
	prom.Register(throughput)
}

func StartPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
