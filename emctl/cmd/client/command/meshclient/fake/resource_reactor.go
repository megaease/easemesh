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
	ReactorFunc func(action Action) (handled bool, ret meta.MeshObject, err error)

	// ResourceReactor is a reactor to react test request
	ResourceReactor interface {
		PrependReactor(verb, kind, resource string, f ReactorFunc)
		AddReactor(verb, kind, resource string, f ReactorFunc)
		DoRequest(string, string, string, meta.MeshObject) (meta.MeshObject, error)
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
		(resource == "*" || strings.EqualFold(resource, a.name))
}

func (w *writeActionImpl) GetObject() meta.MeshObject { return w.obj }

func (r *resourceReactor) PrependReactor(verb, kind, resource string, f ReactorFunc) {
	r.reactors = append([]reactor{{matchVerb: verb, matchKind: kind, matchResource: resource, f: f}}, r.reactors...)
}
func (r *resourceReactor) AddReactor(verb, kind, resource string, f ReactorFunc) {
	r.reactors = append(r.reactors, reactor{matchVerb: verb, matchKind: kind, matchResource: resource, f: f})

}
func (r *resourceReactor) DoRequest(verb, kind, resource string, obj meta.MeshObject) (meta.MeshObject, error) {
	var a Action
	action := &actionImpl{
		verb: verb,
		vk: meta.VersionKind{
			APIVersion: "mesh.megaease.com/v1alphla1", Kind: kind,
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
	case "*":
		a = &writeActionImpl{actionImpl: *action, obj: obj}
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
		panic("unknown reactor type {" + reactorType + "}")
	}
	return reactor
}
