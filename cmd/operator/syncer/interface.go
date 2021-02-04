/*
Copyright 2021 MegaEase Ltd

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
