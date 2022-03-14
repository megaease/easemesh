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

package stdlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

const (
	// https://github.com/megaease/easemesh/blob/main/docs/sidecar-protocol.md#easemesh-traffic-hosting
	agentPort = 9900
)

var (
	agentAddr = fmt.Sprintf(":%d", agentPort)

	// DefaultAgent is the default global agent.
	DefaultAgent = NewAgent()
)

type (
	// Agent is the agent entry.
	Agent struct {
		headers atomic.Value // type: []string
	}

	// AgentConfig is the config pushed to agent.
	AgentConfig struct {
		Headers string `json:"easeagent.progress.forwarded.headers"`
	}

	// AgentHandler is the HTTP handler wrapper.
	AgentHandler struct {
		handlerFunc http.HandlerFunc
	}
)

// NewAgent returns the agent adapting part of EaseMesh sidecar protocol.
func NewAgent() *Agent {
	a := &Agent{}
	a.headers.Store([]string{})
	return a
}

// ServeDefault just runs global default agent in HTTP server,
// please notice it prints logs if the server failed listening.
// The caller must call it to activate default agent.
func ServeDefault() {
	go func() {
		err := http.ListenAndServe(agentAddr, DefaultAgent)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("easemesh agent listen %s failed: %v", agentAddr, err)
		}
	}()
}

// ServeHTTP serves function as agent such as health checking, config updating.
func (a *Agent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		return
	}

	if r.URL.Path == "/config" {
		a.handleConfig(w, r)
	}
}

func (a *Agent) handleConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read config body failed: %v", err)
		return
	}

	config := &AgentConfig{}
	err = json.Unmarshal(body, config)
	if err != nil {
		log.Printf("unmarshal config body failed: %v", err)
		return
	}

	headers := strings.Split(config.Headers, ",")
	a.headers.Store(headers)
}

// Headers returns HTTP header keys which need to be transmit along the chain.
func (a *Agent) Headers() []string {
	return a.headers.Load().([]string)
}

// Headers is the wrapper of Headers of default agent.
func Headers() []string {
	return DefaultAgent.Headers()
}

// WrapHandleFunc wraps http.HandleFunc to serve agent role.
// It copies canary headers to reponse writer, the function itself must not alter them.
func (a *Agent) WrapHandleFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys := a.headers.Load().([]string)
		for k, v := range r.Header.Clone() {
			if strInSlice(k, keys) {
				w.Header()[k] = v
			}
		}

		fn(w, r)

		// NOTE: Copying headers after fn it might not take effect,
		// in the case of fn invoking w.WriteHeader.
	}
}

// WrapHandleFunc is the wrapper of WrapHandleFunc of default agent.
func WrapHandleFunc(fn http.HandlerFunc) http.HandlerFunc {
	return DefaultAgent.WrapHandleFunc(fn)
}

// WrapHandler wraps http.Handler to serve agent role.
func (a *Agent) WrapHandler(handler http.Handler) http.Handler {
	return &AgentHandler{
		handlerFunc: a.WrapHandleFunc(handler.ServeHTTP),
	}
}

// WrapHandler is the wrapper of WrapHandler of default agent.
func WrapHandler(fn http.Handler) http.Handler {
	return DefaultAgent.WrapHandler(fn)
}

func (h *AgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

func strInSlice(s string, ss []string) bool {
	for _, str := range ss {
		if s == str {
			return true
		}
	}

	return false
}
