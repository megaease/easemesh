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

package fake

import (
	"strings"

	"github.com/megaease/easemeshctl/cmd/client/resource/meta"

	"github.com/pkg/errors"
)

type (
	// Action represents the generic action of the request
	Action interface {
		GetVerb() string
		GetVersionKind() meta.VersionKind
		GetName() string
		Matches(verb, kind, resource string) bool
	}

	// WriteAction represents the action will modify action
	WriteAction interface {
		Action
		GetObject() meta.MeshObject
	}

	actionImpl struct {
		verb string
		vk   meta.VersionKind
		name string
	}

	writeActionImpl struct {
		actionImpl
		obj meta.MeshObject
	}

	// ReactorFunc represent a mock behavior
	ReactorFunc func(action Action) (handled bool, rets []meta.MeshObject, err error)

	// ResourceReactor is a reactor to react test request
	ResourceReactor interface {
		DoRequest(verb, kind, resource string, obj meta.MeshObject) ([]meta.MeshObject, error)
	}

	// ResourceReactorBuilder is a builder to build a ResourceReactor
	ResourceReactorBuilder struct {
		reactorType string
		reactors    []reactor
	}

	reactor struct {
		matchVerb     string
		matchKind     string
		matchResource string
		f             ReactorFunc
	}
	resourceReactor struct {
		reactors []reactor
	}
)

func (a *actionImpl) GetVerb() string                  { return a.verb }
func (a *actionImpl) GetVersionKind() meta.VersionKind { return a.vk }
func (a *actionImpl) GetName() string                  { return a.name }

func (a *actionImpl) Matches(verb, kind, resource string) bool {
	return (verb == "*" || strings.EqualFold(verb, a.verb)) &&
		(kind == "*" || strings.EqualFold(kind, a.vk.Kind)) &&
		(resource == "" || resource == "*" || strings.EqualFold(resource, a.name))
}

func (w *writeActionImpl) GetObject() meta.MeshObject { return w.obj }

func (r *resourceReactor) DoRequest(verb, kind, resource string, obj meta.MeshObject) ([]meta.MeshObject, error) {
	var a Action
	action := &actionImpl{
		verb: verb,
		vk: meta.VersionKind{
			APIVersion: "mesh.megaease.com/v2alpha1", Kind: kind,
		},
	}
	switch verb {
	case "get":
		a = action
	case "create":
		fallthrough
	case "update":
		fallthrough
	case "delete":
		fallthrough
	case "list":
		fallthrough
	case "*":
		a = &writeActionImpl{actionImpl: *action, obj: obj}
	default:
		return nil, errors.Errorf("unknown operation %s", verb)
	}
	for _, reactor := range r.reactors {
		if a.Matches(reactor.matchVerb, reactor.matchKind, reactor.matchResource) {
			h, o, e := reactor.f(a)
			if h {
				return o, e
			}
		}
	}
	return nil, errors.Errorf("no any reactor process this request")
}

var globalReactor = map[string]ResourceReactor{
	"test": &resourceReactor{},
}

// ResourceReactorForType return a specific global ResourceReactor
func ResourceReactorForType(reactorType string) ResourceReactor {
	reactor := globalReactor[reactorType]
	if reactor == nil {
		return nil
	}
	return reactor
}

// PrependReactor prepend a reactor header
func (r *ResourceReactorBuilder) PrependReactor(verb, kind, resource string, f ReactorFunc) *ResourceReactorBuilder {
	r.reactors = append([]reactor{{matchVerb: verb, matchKind: kind, matchResource: resource, f: f}}, r.reactors...)
	return r
}

// AddReactor add a reactor to tail
func (r *ResourceReactorBuilder) AddReactor(verb, kind, resource string, f ReactorFunc) *ResourceReactorBuilder {
	r.reactors = append(r.reactors, reactor{matchVerb: verb, matchKind: kind, matchResource: resource, f: f})
	return r
}

// Added add construct a resourceReactor and insert it to globalReactor
func (r *ResourceReactorBuilder) Added() {
	globalReactor[r.reactorType] = &resourceReactor{r.reactors}
}

// NewResourceReactorBuilder return a ResourceReactorBuilder
func NewResourceReactorBuilder(t string) *ResourceReactorBuilder {
	return &ResourceReactorBuilder{reactorType: t}
}
