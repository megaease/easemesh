//go:generate go run github.com/megaease/easemeshctl/cmd/transformer Service Canary=canary

/*
 * Copyright (c) 2017, MegaEase
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

package meshclient

import (
	"context"

	"github.com/megaease/easemeshctl/cmd/client/resource"
)

// CanaryGetter represents a Canary resource accessor
type CanaryGetter interface {
	Canary() CanaryInterface
}

// CanaryInterface captures the set of operations for interacting with the EaseMesh REST apis of the canaray resource.
type CanaryInterface interface {
	Get(context.Context, string) (*resource.Canary, error)
	Patch(context.Context, *resource.Canary) error
	Create(context.Context, *resource.Canary) error
	Delete(context.Context, string) error
	List(context.Context) ([]*resource.Canary, error)
}
