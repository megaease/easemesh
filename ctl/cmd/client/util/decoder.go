package util

import (
	"encoding/json"

	"github.com/megaease/easemeshctl/cmd/client/resource"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Decoder interface {
	Decode(data []byte) (resource.MeshObject, *resource.VersionKind, error)
}

type decoder struct {
	oc resource.ObjectCreator
}

func (d *decoder) Decode(data []byte) (resource.MeshObject, *resource.VersionKind, error) {
	vk := &resource.VersionKind{}
	err := yaml.Unmarshal(data, vk)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unmarshal data to resource.VersionKind failed")
	}

	meshObject, err := d.oc.New(vk)
	if err != nil {
		return nil, vk, err
	}

	err = json.Unmarshal(data, meshObject)
	if err != nil {
		return nil, vk, errors.Wrap(err, "unmarshal data to MeshObject error")
	}
	return meshObject, vk, nil
}

func newDefaultDecoder() Decoder {
	return &decoder{oc: resource.NewObjectCreator()}
}
