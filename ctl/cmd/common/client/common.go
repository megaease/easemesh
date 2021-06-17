package client

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type UnmarshalFunc func([]byte, int) (interface{}, error)

type HTTPJSONResponseHandler interface {
	HandleResponse(UnmarshalFunc) (interface{}, error)
}

type HTTPJSONClient interface {
	Post(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	Delete(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	Patch(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	Put(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	Get(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
}

type Option func(*resty.Client)
type httpJSONClient struct {
	options []Option
}

type httpJSONResponseFunc func(UnmarshalFunc) (interface{}, error)

func (h httpJSONResponseFunc) HandleResponse(fn UnmarshalFunc) (interface{}, error) {
	return h(fn)
}

func NewHTTPJSON(o ...Option) HTTPJSONClient {
	return &httpJSONClient{options: o}
}

func WrapRetryOptions(retryCount int, retryWaitTime time.Duration, conditionFunc func(b []byte, err error) bool) []Option {
	return []Option{
		func(client *resty.Client) {
			client.SetRetryCount(retryCount)
		},
		func(client *resty.Client) {
			client.SetRetryWaitTime(retryWaitTime)
		},
		func(client *resty.Client) {
			client.AddRetryCondition(func(r *resty.Response, e error) bool {
				return conditionFunc(r.Body(), e)
			})
		},
	}
}

func (h *httpJSONClient) setupClient(timeout time.Duration, extraHeaders map[string]string) *resty.Client {

	client := resty.New()
	client.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTimeout(timeout)

	for _, o := range h.options {
		o(client)
	}

	if len(extraHeaders) != 0 {
		for k, v := range extraHeaders {
			client.SetHeader(k, v)
		}
	}
	return client
}
func closeRawBody(r *resty.Response) {
	if r != nil && r.RawBody() != nil {
		defer r.RawBody().Close()
	}
}

func (h *httpJSONClient) Post(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Post(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, errors.Wrapf(err, "Post to url %s error", url)
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Delete(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Delete(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, errors.Wrapf(err, "Post to url %s error", url)
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Patch(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Patch(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, errors.Wrapf(err, "Post to url %s error", url)
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Put(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Put(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, errors.Wrapf(err, "Post to url %s error", url)
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Get(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {

	client := h.setupClient(timeout, extraHeaders)
	r, err := client.R().Get(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, errors.Wrapf(err, "Post to url %s error", url)
		}
		return fn(r.Body(), r.StatusCode())
	})
}
