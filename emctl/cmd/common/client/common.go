/*
 * Copyright (c) 2021, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
)

// UnmarshalFunc is a function to unmarshal bytes into object
type UnmarshalFunc func([]byte, int) (interface{}, error)

// HTTPJSONResponseHandler is a handler to handle http response body
type HTTPJSONResponseHandler interface {
	HandleResponse(UnmarshalFunc) (interface{}, error)
}

// HTTPJSONClient wraps http json client
type HTTPJSONClient interface {
	Post(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	PostByContext(context.Context, string, interface{}, map[string]string) HTTPJSONResponseHandler
	Delete(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	DeleteByContext(context.Context, string, interface{}, map[string]string) HTTPJSONResponseHandler
	Patch(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	PatchByContext(context.Context, string, interface{}, map[string]string) HTTPJSONResponseHandler
	Put(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	PutByContext(context.Context, string, interface{}, map[string]string) HTTPJSONResponseHandler
	Get(string, interface{}, time.Duration, map[string]string) HTTPJSONResponseHandler
	GetByContext(context.Context, string, interface{}, map[string]string) HTTPJSONResponseHandler
}

// Option is option function
type Option func(*resty.Client)
type httpJSONClient struct {
	options []Option
}

type httpJSONResponseFunc func(UnmarshalFunc) (interface{}, error)

func (h httpJSONResponseFunc) HandleResponse(fn UnmarshalFunc) (interface{}, error) {
	return h(fn)
}

// NewHTTPJSON creates a HTTPJSONClient
func NewHTTPJSON(o ...Option) HTTPJSONClient {
	return &httpJSONClient{options: o}
}

// WrapRetryOptions wraps option to retryer
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

func (h *httpJSONClient) setupClient(timeout *time.Duration, extraHeaders map[string]string) *resty.Client {

	client := resty.New()
	client.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	if timeout != nil {
		client.SetTimeout(*timeout)

	}

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
	client := h.setupClient(&timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Post(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}
func (h *httpJSONClient) PostByContext(ctx context.Context, url string, reqBody interface{}, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(nil, extraHeaders)
	r, err := client.R().SetContext(ctx).SetBody(reqBody).Post(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Delete(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(&timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Delete(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) DeleteByContext(ctx context.Context, url string, reqBody interface{}, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(nil, extraHeaders)
	r, err := client.R().SetContext(ctx).SetBody(reqBody).Delete(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Patch(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(&timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Patch(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) PatchByContext(ctx context.Context, url string, reqBody interface{}, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(nil, extraHeaders)
	r, err := client.R().SetContext(ctx).SetBody(reqBody).Patch(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Put(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(&timeout, extraHeaders)
	r, err := client.R().SetBody(reqBody).Put(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) PutByContext(ctx context.Context, url string, reqBody interface{}, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(nil, extraHeaders)
	r, err := client.R().SetContext(ctx).SetBody(reqBody).Put(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) Get(url string, reqBody interface{}, timeout time.Duration, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(&timeout, extraHeaders)
	r, err := client.R().Get(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}

func (h *httpJSONClient) GetByContext(ctx context.Context, url string, reqBody interface{}, extraHeaders map[string]string) HTTPJSONResponseHandler {
	client := h.setupClient(nil, extraHeaders)
	r, err := client.R().SetContext(ctx).Get(url)
	return (httpJSONResponseFunc)(func(fn UnmarshalFunc) (interface{}, error) {
		defer closeRawBody(r)

		if err != nil {
			return nil, err
		}
		return fn(r.Body(), r.StatusCode())
	})
}
