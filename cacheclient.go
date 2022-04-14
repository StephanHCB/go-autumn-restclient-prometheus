package aurestclientprometheus

import (
	"context"
	"fmt"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	aurestcaching "github.com/StephanHCB/go-autumn-restclient/implementation/caching"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	cacheHitCounts  *prometheus.CounterVec
	cacheMissCounts *prometheus.CounterVec
	cacheErrCounts  *prometheus.CounterVec
)

func SetupCacheClientMetrics() {
	SetupCommon()

	if cacheHitCounts == nil {
		cacheHitCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_cache_hits_count",
				Help: "Number of cache hits by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(cacheHitCounts)
	}

	if cacheMissCounts == nil {
		cacheMissCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_cache_misses_count",
				Help: "Number of cache misses by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(cacheMissCounts)
	}

	if cacheErrCounts == nil {
		cacheErrCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_cache_errors_count",
				Help: "Number of cache errors (invalid entries) by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(cacheErrCounts)
	}

}

func InstrumentCacheClient(client aurestclientapi.Client) {
	SetupCacheClientMetrics()
	aurestcaching.Instrument(client, CacheHitMetricsCallback, CacheMissMetricsCallback, CacheInvalidMetricsCallback)
}

func CacheHitMetricsCallback(_ context.Context, method string, requestUrl string, status int, _ error, _ time.Duration, _ int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)
	cacheHitCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
}

func CacheMissMetricsCallback(_ context.Context, method string, requestUrl string, status int, _ error, _ time.Duration, _ int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)
	cacheMissCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
}

func CacheInvalidMetricsCallback(_ context.Context, method string, requestUrl string, status int, _ error, _ time.Duration, _ int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)
	cacheErrCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
}
