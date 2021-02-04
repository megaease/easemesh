package syncer

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
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

	log = logf.Log.WithName("syncer")
)

// ObjectSyncer is a syncer.Interface for syncing object to K8s
type objectSyncer struct {
	Owner          client.Object
	Self           client.Object
	SyncFn         controllerutil.MutateFn
	Name           string
	Client         client.Client
	Scheme         *runtime.Scheme
	previousObject runtime.Object
}

var _ Interface = &objectSyncer{}

// Object returns the ObjectSyncer subject
func (s *objectSyncer) Object() client.Object {
	return s.Self
}

// ObjectOwner returns the owner of ObjectSyncer subject
func (s *objectSyncer) ObjectOwner() client.Object {
	return s.Owner
}

func (s *objectSyncer) objectType() string {
	return fmt.Sprintf("%T", s.Self)
}

func (s *objectSyncer) ownerType() string {
	return fmt.Sprintf("%T", s.Owner)
}

// Sync does the actual syncing and implements the syncer.Inteface Sync method
func (s *objectSyncer) Sync(ctx context.Context) (SyncResult, error) {
	result := SyncResult{}

	key, err := getKey(s.Self)
	if err != nil {
		return result, err
	}

	result.Operation, err = controllerutil.CreateOrUpdate(ctx, s.Client, s.Self, s.mutateFn())

	diff := deep.Equal(stripSecrets(s.previousObject), stripSecrets(s.Self))

	// don't pass to user error for owner deletion, just don't create the object
	// nolint: gocritic
	if err == ErrOwnerDeleted {
		log.Info(string(result.Operation), "key", key, "kind", s.objectType(), "error", err)
		err = nil
	} else if err == ErrIgnored {
		log.V(1).Info("syncer skipped", "key", key, "kind", s.objectType())
		err = nil
	} else if err != nil {
		result.SetEventData(eventWarning, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s failed syncing: %s", s.objectType(), key, err))
		log.Error(err, string(result.Operation), "key", key, "kind", s.objectType(), "diff", diff)
	} else {
		result.SetEventData(eventNormal, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s %s successfully", s.objectType(), key, result.Operation))
		log.V(1).Info(string(result.Operation), "key", key, "kind", s.objectType(), "diff", diff)
	}
	return result, err
}

// Given an ObjectSyncer, returns a controllerutil.MutateFn which also sets the
// owner reference if the subject has one
func (s *objectSyncer) mutateFn() controllerutil.MutateFn {
	return func() error {
		s.previousObject = s.Self.DeepCopyObject()

		err := s.SyncFn()
		if err != nil {
			return err
		}

		if s.Owner != nil {
			existingMeta, ok := s.Self.(metav1.Object)
			if !ok {
				return errors.Errorf("%s is not a metav1.Object", s.objectType())
			}

			ownerMeta, ok := s.Owner.(metav1.Object)
			if !ok {
				return errors.Errorf("%s is not a metav1.Object", s.ownerType())
			}

			// set owner reference only if owner resource is not being deleted, otherwise the owner
			// reference will be reset in case of deleting with cascade=false.
			if ownerMeta.GetDeletionTimestamp().IsZero() {
				err := controllerutil.SetControllerReference(ownerMeta, existingMeta, s.Scheme)
				if err != nil {
					return err
				}
			} else if ctime := existingMeta.GetCreationTimestamp(); ctime.IsZero() {
				// the owner is deleted, don't recreate the resource if does not exist, because gc
				// will not delete it again because has no owner reference set
				return ErrOwnerDeleted
			}
		}

		return nil
	}
}
