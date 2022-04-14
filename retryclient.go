package aurestclientprometheus

import (
	"context"
	"fmt"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	aurestretry "github.com/StephanHCB/go-autumn-restclient/implementation/retry"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	retryCounts  *prometheus.CounterVec
	giveUpCounts *prometheus.CounterVec
)

func SetupRetryClientMetrics() {
	SetupCommon()

	if retryCounts == nil {
		retryCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_retries_count",
				Help: "Number of attempted retries by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(retryCounts)
	}

	if giveUpCounts == nil {
		giveUpCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_retries_give_up_count",
				Help: "Number of times retry had to give up by target hostname, method, outcome and response status.",
			},
			[]string{"clientName", "method", "outcome", "status"},
		)
		prometheus.MustRegister(giveUpCounts)
	}
}

func InstrumentRetryClient(client aurestclientapi.Client) {
	SetupRetryClientMetrics()
	aurestretry.Instrument(client, RetryingMetricsCallback, GivingUpMetricsCallback)
}

func RetryingMetricsCallback(_ context.Context, method string, requestUrl string, status int, _ error, _ time.Duration, _ int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)
	retryCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
}

func GivingUpMetricsCallback(_ context.Context, method string, requestUrl string, status int, _ error, _ time.Duration, _ int) {
	clientName := ClientNameFromRequestUrl(requestUrl)
	outcome := OutcomeFromStatus(status)
	statusStr := fmt.Sprintf("%d", status)
	giveUpCounts.WithLabelValues(clientName, method, outcome, statusStr).Inc()
}
