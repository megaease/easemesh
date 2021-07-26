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

	"github.com/megaease/easemesh/mesh-operator/pkg/base"

	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	eventNormal  = "Normal"
	eventWarning = "Warning"
)

var (
	//ErrOwnerDeleted represent an error which syncer.Object's owner has been deleted
	ErrOwnerDeleted = errors.New("Owner has been deleted")
	//ErrIgnored represnets an ignoredable error for syncer
	ErrIgnored = errors.New("ignored error")
)

// k8sSyncer is a Syncer for syncing object to Kubernetes.
type k8sSyncer struct {
	*base.Runtime
	log logr.Logger

	owner         client.Object
	ownee         client.Object
	previousOwnee runtime.Object
	userMutateFn  controllerutil.MutateFn
}

var _ Syncer = &k8sSyncer{}

// Object returns the ObjectSyncer subject
func (s *k8sSyncer) Object() client.Object {
	return s.ownee
}

// ObjectOwner returns the owner of ObjectSyncer subject
func (s *k8sSyncer) ObjectOwner() client.Object {
	return s.owner
}

func (s *k8sSyncer) owneeType() string {
	return fmt.Sprintf("%T", s.ownee)
}

func (s *k8sSyncer) ownerType() string {
	return fmt.Sprintf("%T", s.owner)
}

// Sync does the actual syncing and implements the syncer.Inteface Sync method
func (s *k8sSyncer) Sync(ctx context.Context) (SyncResult, error) {
	result := SyncResult{}

	id, err := idOfObject(s.ownee)
	if err != nil {
		return result, err
	}

	result.Operation, err = controllerutil.CreateOrUpdate(ctx, s.Client, s.ownee, s.mutateFn())

	diff := deep.Equal(stripSecrets(s.previousOwnee), stripSecrets(s.ownee))

	// NOTE: Owner deletion is not an error to report.
	// nolint: gocritic
	if err == ErrOwnerDeleted {
		s.log.Info(string(result.Operation), "id", id, "ownerKind", s.ownerType(), "message", err)
		err = nil
	} else if err == ErrIgnored {
		s.log.V(1).Info("syncer skipped", "id", id, "owneeKind", s.owneeType())
		err = nil
	} else if err != nil {
		result.SetEventData(eventWarning, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s failed syncing: %s", s.owneeType(), id, err))
		s.log.V(1).Error(err, string(result.Operation), "id", id, "owneeKind", s.owneeType(), "diff", diff)
	} else {
		result.SetEventData(eventNormal, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s %s successfully", s.owneeType(), id, result.Operation))
		s.log.V(1).Info(string(result.Operation), "id", id, "owneeKind", s.owneeType(), "diff", diff)
	}

	return result, err
}

// mutateFn wraps user-defined mutateFn by setting owner reference if it has.
func (s *k8sSyncer) mutateFn() controllerutil.MutateFn {
	return func() error {
		s.previousOwnee = s.ownee.DeepCopyObject()

		err := s.userMutateFn()
		if err != nil {
			return err
		}

		if s.owner != nil {
			owneeMeta, ok := s.ownee.(metav1.Object)
			if !ok {
				return errors.Errorf("%s is not a metav1.Object", s.owneeType())
			}

			ownerMeta, ok := s.owner.(metav1.Object)
			if !ok {
				return errors.Errorf("%s is not a metav1.Object", s.ownerType())
			}

			// NOTE: Set owner reference only if owner resource is not being deleted
			// otherwise the owner reference will be reset in case of deleting with cascade=false.
			if ownerMeta.GetDeletionTimestamp().IsZero() {
				err := controllerutil.SetControllerReference(ownerMeta, owneeMeta, s.Scheme)
				if err != nil {
					return err
				}
			} else if ctime := owneeMeta.GetCreationTimestamp(); ctime.IsZero() {
				// NOTE: The owner is deleted, don't recreate the resource if it does not exist,
				// because gc will not delete it again since it's not set owner reference.
				return ErrOwnerDeleted
			}
		}

		return nil
	}
}
