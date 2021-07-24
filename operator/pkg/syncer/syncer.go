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
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type (
	// Syncer represents a syncer. A syncer persists an object
	// (known as subject), into a store (kubernetes apiserver or generic stores)
	// and records kubernetes events.
	Syncer interface {
		// Object returns the object which is applies.
		Object() client.Object

		// ObjectOwner returns the owner who owns the object, nil means none.
		ObjectOwner() client.Object

		// Sync synchronizes the data.
		Sync(context.Context) (SyncResult, error)
	}

	// SyncResult is a result of an Sync call
	SyncResult struct {
		Operation    controllerutil.OperationResult
		EventType    string
		EventReason  string
		EventMessage string
	}
)

// SetEventData sets event data on an SyncResult
func (r *SyncResult) SetEventData(eventType, reason, message string) {
	r.EventType = eventType
	r.EventReason = reason
	r.EventMessage = message
}

// Sync mutates the subject of the syncer interface via controller-runtime
// CreateOrUpdate method.
func Sync(ctx context.Context, syncer Syncer, recorder record.EventRecorder) error {
	// NOTE: Export the result to caller if necessary in the future.
	result, err := syncer.Sync(ctx)

	owner := syncer.ObjectOwner()
	if recorder != nil && owner != nil && result.EventType != "" && result.EventReason != "" && result.EventMessage != "" {
		if err != nil || result.Operation != controllerutil.OperationResultNone {
			recorder.Eventf(owner, result.EventType, result.EventReason, result.EventMessage)
		}
	}

	return err
}

// New return a syncer.Interface object.
func New(baseRuntime *base.Runtime, owner, ownee client.Object, fn controllerutil.MutateFn) Syncer {
	return &k8sSyncer{
		Runtime:      baseRuntime,
		log:          baseRuntime.Log.WithName("syncer"),
		owner:        owner,
		ownee:        ownee,
		userMutateFn: fn,
	}
}

func idOfObject(obj client.Object) (types.NamespacedName, error) {
	key := types.NamespacedName{}
	objMeta, ok := obj.(metav1.Object)
	if !ok {
		return key, fmt.Errorf("%T is not a metav1.Object", obj)
	}

	key.Name = objMeta.GetName()
	key.Namespace = objMeta.GetNamespace()

	return key, nil
}

func basicEventReason(objKindName string, err error) string {
	if err != nil {
		return fmt.Sprintf("%sSyncFailed", strcase.ToCamel(objKindName))
	}

	return fmt.Sprintf("%sSyncSuccessfully", strcase.ToCamel(objKindName))
}

// stripSecrets returns a copy for the secret without secret data.
func stripSecrets(obj runtime.Object) runtime.Object {
	// if obj is secret, don't print secret data
	s, ok := obj.(*corev1.Secret)
	if ok {
		cObj := s.DeepCopyObject().(*corev1.Secret)
		cObj.Data = nil
		cObj.StringData = nil
		return cObj
	}

	return obj
}
