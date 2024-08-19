package aurestclientprometheus

import "net/http"

type PrometheusRoundTripper struct {
	wrapped http.RoundTripper
}

func NewPrometheusRoundTripper(wrapped http.RoundTripper) *PrometheusRoundTripper {
	SetupHttpClientMetrics()

	return &PrometheusRoundTripper{
		wrapped: wrapped,
	}
}

func (vrt *PrometheusRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	RequestMetricsCallback(req.Context(), req.Method, req.RequestURI, 0, e, 0, int(req.ContentLength))

	response, err := vrt.wrapped.RoundTrip(req)

	statusCode := 0
	if response != nil {
		statusCode = response.StatusCode
	}
	ResponseMetricsCallback(req.Context(), req.Method, req.RequestURI, statusCode, err, 0, 0)

	return response, err
}
