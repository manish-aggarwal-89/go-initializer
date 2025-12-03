package util

import (
	"bytes"
	"context"
	"errors"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

const (
	ProtocolTag  = "protocol"
	RestProtocol = "REST"
)

// Request struct
type Request struct {
	client *httpclient.Client
	statsd metrics.MonitorStatsd
}

type Option func(r *Request)

func SetStatsd(statsd metrics.MonitorStatsd) Option {
	return func(r *Request) {
		r.statsd = statsd
	}
}

type HTTPRequest interface {
	Do(context context.Context, method, url string, reqBody []byte, headers map[string]string) (*http.Response, error)
	DoWithMetric(context context.Context, method, url string, reqBody []byte, headers map[string]string) (*http.Response, error)
}

// NewHTTPRequest function
// Request's Constructor
// Returns : *Request
func NewHTTPRequest(retries int, timeout time.Duration, opts ...Option) HTTPRequest {
	// define a maximum jitter interval
	maximumJitterInterval := 5 * time.Millisecond

	// create a backoff
	backoff := heimdall.NewConstantBackoff(2000, maximumJitterInterval)

	// create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	baseHttpClient := &http.Client{
		Timeout: timeout,
	}

	// set http client
	client := httpclient.NewClient(
		httpclient.WithHTTPClient(baseHttpClient),
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(retries),
	)

	request := &Request{
		client: client,
	}

	for _, opt := range opts {
		opt(request)
	}

	return request
}

func (r *Request) Do(_ context.Context, method, url string, requestBody []byte, headers map[string]string) (*http.Response, error) {

	// set request http
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// iterate optional data of headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// client request
	response, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	// close response body
	//defer response.Body.Close()

	return response, err
}

// DoWithMetric function, for http client call and send metric
func (r *Request) DoWithMetric(_ context.Context, method, url string, requestBody []byte, headers map[string]string) (*http.Response, error) {
	if r.statsd == nil {
		return nil, errors.New("statsd client is nil. you can inject the option ")
	}

	var (
		// respStatus string
		startTime = time.Now()
	)

	// set request http
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// iterate optional data of headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// client request
	response, err := r.client.Do(req)
	go r.sendMetric(req.URL.Path, time.Since(startTime), response, err, req.Method)
	if err != nil {
		return nil, err
	}
	// close response body
	// defer response.Body.Close()

	return response, err
}

func (r *Request) sendMetric(entity string, latency time.Duration, response *http.Response, httpError error, httpMethod string) {
	var (
		httpStatusCode int
		status         = metrics.Failed
		tags           = make(map[string]interface{})
	)

	tags[ProtocolTag] = RestProtocol

	if httpError == nil && response != nil {
		status = metrics.Success
		httpStatusCode = response.StatusCode
	}

	entity = httpMethod + ":" + entity

	err := r.statsd.CustomMonitorLatency(
		entity,
		metrics.API_OUT,
		status,
		httpStatusCode,
		tags,
		latency,
	)
	if err != nil {
		log.Errorf("Error sending metric to heimdall: %s", err.Error())
	}
}
