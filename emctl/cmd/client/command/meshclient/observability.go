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

//go:generate go run github.com/megaease/easemeshctl/cmd/transformer Service.Observability ObservabilityTracings=services/%s/tracings ObservabilityMetrics=services/%s/metrics ObservabilityOutputServer=services/%s/outputserver

package meshclient

import (
	"context"

	"github.com/megaease/easemeshctl/cmd/client/resource"
)

// ObservabilityGetter represents an Observability resource accessor
type ObservabilityGetter interface {
	ObservabilityTracings() ObservabilityTracingsInterface
	ObservabilityMetrics() ObservabilityMetricsInterface
	ObservabilityOutputServer() ObservabilityOutputServerInterface
}

// ObservabilityOutputServerInterface captures the set of operations for interacting with the EaseMesh REST apis of the observability output server resource.
type ObservabilityOutputServerInterface interface {
	Get(context.Context, string) (*resource.ObservabilityOutputServer, error)
	Patch(context.Context, *resource.ObservabilityOutputServer) error
	Create(context.Context, *resource.ObservabilityOutputServer) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.ObservabilityOutputServer, error)
}

// ObservabilityMetricsInterface captures the set of operations for interacting with the EaseMesh REST apis of the observability metric resource.
type ObservabilityMetricsInterface interface {
	Get(context.Context, string) (*resource.ObservabilityMetrics, error)
	Patch(context.Context, *resource.ObservabilityMetrics) error
	Create(context.Context, *resource.ObservabilityMetrics) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.ObservabilityMetrics, error)
}

// ObservabilityTracingsInterface captures the set of operations for interacting with the EaseMesh REST apis of the observability tracings resource.
type ObservabilityTracingsInterface interface {
	Get(context.Context, string) (*resource.ObservabilityTracings, error)
	Patch(context.Context, *resource.ObservabilityTracings) error
	Create(context.Context, *resource.ObservabilityTracings) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.ObservabilityTracings, error)
}
