package aurestclientprometheus

import (
	"context"
	"fmt"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	httpClientCounts      *prometheus.CounterVec
	httpClientErrCounts   *prometheus.CounterVec
	httpClientReqBytes    *prometheus.SummaryVec
	httpClientResBytes    *prometheus.SummaryVec
	httpClientLatencySums *prometheus.SummaryVec
)

func SetupHttpClientMetrics() {
	SetupCommon()

	if httpClientCounts == nil {
		httpClientCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_requests_seconds_count",
				Help: "Number of downstream http requests by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(httpClientCounts)
	}

	if httpClientErrCounts == nil {
		httpClientErrCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_requests_errors_count",
				Help: "Number of downstream http requests that raised a technical error by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(httpClientErrCounts)
	}

	if httpClientReqBytes == nil {
		httpClientReqBytes = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_client_requests_request_bytes_sum",
				Help: "Size of the request by target hostname and method.",
			},
			[]string{"clientName", "method"},
		)
		prometheus.MustRegister(httpClientReqBytes)
	}

	if httpClientResBytes == nil {
		httpClientResBytes = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_client_requests_response_bytes_sum",
				Help: "Size of the response by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(httpClientResBytes)
	}

	if httpClientLatencySums == nil {
		httpClientLatencySums = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_client_requests_seconds_sum",
				Help: "How long it took to process downstream http requests by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(httpClientLatencySums)
	}
}

func InstrumentHttpClient(client aurestclientapi.Client) {
	SetupHttpClientMetrics()
	auresthttpclient.Instrument(client, RequestMetricsCallback, ResponseMetricsCallback)
}

func RequestMetricsCallback(ctx context.Context, method string, requestUrl string, status int, err error, latency time.Duration, size int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)

	// if no error, do not count the request at this point, or we may double count it
	if err != nil {
		httpClientErrCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
	}
	if size > 0 {
		httpClientReqBytes.WithLabelValues(clientName, method).Observe(float64(size))
	}
}

func ResponseMetricsCallback(ctx context.Context, method string, requestUrl string, status int, err error, latency time.Duration, size int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)

	httpClientCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()

	if size > 0 {
		httpClientResBytes.WithLabelValues(clientName, method, outcome, statusStr).Observe(float64(size))
	}
	if latency > 0 {
		httpClientLatencySums.WithLabelValues(clientName, method, outcome, statusStr).Observe(float64(latency.Microseconds()) / 1000000)
	}
	if err != nil {
		httpClientErrCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
	}
}
