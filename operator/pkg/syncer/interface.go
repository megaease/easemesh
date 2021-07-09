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

package syncer

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// SyncResult is a result of an Sync call
type SyncResult struct {
	Operation    controllerutil.OperationResult
	EventType    string
	EventReason  string
	EventMessage string
}

// SetEventData sets event data on an SyncResult
func (r *SyncResult) SetEventData(eventType, reason, message string) {
	r.EventType = eventType
	r.EventReason = reason
	r.EventMessage = message
}

// Interface represents a syncer. A syncer persists an object
// (known as subject), into a store (kubernetes apiserver or generic stores)
// and records kubernetes events
type Interface interface {
	// Object returns the object for which sync applies
	Object() client.Object

	// ObjectOwner returns the object owner or nil if object does not have one
	ObjectOwner() client.Object

	//Sync persists data into the kube-apiserver
	Sync(context.Context) (SyncResult, error)
}
