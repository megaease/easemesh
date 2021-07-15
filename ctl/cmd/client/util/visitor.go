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

package util

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/resource"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Here, we massively refer the implementation of kubernetes client-runtime
// https://github.com/kubernetes/cli-runtime/blob/master/pkg/resource/visitor.go.

const (
	constSTDINstr = "STDIN"
)

type Visitor interface {
	Visit(VisitorFunc) error
}

type VisitorFunc func(resource.MeshObject, error) error

type RawExtension struct {
	Raw []byte `json:"-" protobuf:"bytes,1,opt,name=raw"`
}

func (re *RawExtension) UnmarshalJSON(in []byte) error {
	if re == nil {
		return errors.New("runtime.RawExtension: UnmarshalJSON on nil pointer")
	}
	if !bytes.Equal(in, []byte("null")) {
		re.Raw = append(re.Raw[0:0], in...)
	}
	return nil
}

func (re RawExtension) MarshalJSON() ([]byte, error) {
	return re.Raw, nil
}

func ignoreFile(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return false
	}
	ext := filepath.Ext(path)
	for _, s := range extensions {
		if s == ext {
			return false
		}
	}
	return true
}

// FileVisitorForSTDIN return a special FileVisitor just for STDIN
func FileVisitorForSTDIN(decoder Decoder) Visitor {
	return &FileVisitor{
		Path:          constSTDINstr,
		StreamVisitor: NewStreamVisitor(nil, decoder, constSTDINstr),
	}
}

func ExpandPathsToFileVisitors(decoder Decoder, paths string, recursive bool, extensions []string) ([]Visitor, error) {
	var visitors []Visitor
	err := filepath.Walk(paths, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if path != paths && !recursive {
				return filepath.SkipDir
			}
			return nil
		}
		// Don't check extension if the filepath was passed explicitly
		if path != paths && ignoreFile(path, extensions) {
			return nil
		}

		visitor := &FileVisitor{
			Path:          path,
			StreamVisitor: NewStreamVisitor(nil, decoder, path),
		}

		visitors = append(visitors, visitor)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return visitors, nil
}

// FileVisitor is wrapping around a StreamVisitor, to handle open/close files
type FileVisitor struct {
	Path string
	*StreamVisitor
}

// Visit in a FileVisitor is just taking care of opening/closing files
func (v *FileVisitor) Visit(fn VisitorFunc) error {
	var f *os.File
	if v.Path == constSTDINstr {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(v.Path)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	utf16bom := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	v.StreamVisitor.Reader = transform.NewReader(f, utf16bom)

	return v.StreamVisitor.Visit(fn)
}

type StreamVisitor struct {
	io.Reader

	Decoder Decoder
	Source  string
}

// NewStreamVisitor is a helper function that is useful when we want to change the fields of the struct but keep calls the same.
func NewStreamVisitor(r io.Reader, decoder Decoder, source string) *StreamVisitor {
	return &StreamVisitor{
		Reader:  r,
		Decoder: decoder,
		Source:  source,
	}
}

// Visit implements Visitor over a stream. StreamVisitor is able to distinct multiple resources in one stream.
func (v *StreamVisitor) Visit(fn VisitorFunc) error {
	d := yaml.NewYAMLOrJSONDecoder(v.Reader, 4096)
	for {
		ext := RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				return nil
			}
			return errors.Errorf("error parsing %s: %v", v.Source, err)
		}
		jsonBuff, err := ext.MarshalJSON()
		if err != nil {
			return err
		}

		ext.Raw = bytes.TrimSpace(jsonBuff)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		info, err := v.decodeMeshObject(jsonBuff, v.Source)
		if err != nil {
			if fnErr := fn(info, err); fnErr != nil {
				return fnErr
			}
			continue
		}
		if err := fn(info, nil); err != nil {
			return err
		}
	}
}

func (v *StreamVisitor) decodeMeshObject(data []byte, source string) (resource.MeshObject, error) {
	meshObject, _, err := v.Decoder.Decode(data)
	if err != nil {
		return nil, err
	}
	return meshObject, nil
}

type URLVisitor struct {
	URL *url.URL
	*StreamVisitor
	HttpAttemptCount int
}

func (v *URLVisitor) Visit(fn VisitorFunc) error {
	body, err := readHttpWithRetries(resty.New(), time.Second, v.URL.String(), v.HttpAttemptCount)
	if err != nil {
		return err
	}
	defer body.Close()
	v.StreamVisitor.Reader = body
	return v.StreamVisitor.Visit(fn)
}

// readHttpWithRetries tries to http.Get the v.URL retries times before giving up.
func readHttpWithRetries(client *resty.Client, duration time.Duration, u string, attempts int) (io.ReadCloser, error) {

	r, err := client.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTimeout(duration).
		SetRetryCount(attempts).
		SetRetryWaitTime(duration).
		AddRetryCondition(func(r *resty.Response, e error) bool {
			if e != nil {
				return true
			}

			if r.StatusCode() >= 500 && r.StatusCode() < 600 {
				// Retry 500's
				return true
			}
			return false
		}).
		R().
		Get(u)
	if err != nil {
		return nil, err
	}

	if r.StatusCode() != http.StatusOK {
		defer r.RawBody().Close()
		return nil, errors.Errorf("unable to read URL %q, status code=%d", u, r.StatusCode())
	}
	return r.RawBody(), nil
}

type CommandVisitor struct {
	Kind string
	Name string

	oc resource.ObjectCreator
}

func NewCommandVisitor(kind, name string) *CommandVisitor {
	return &CommandVisitor{
		Kind: adaptCommandKind(kind),
		Name: name,
		oc:   resource.NewObjectCreator(),
	}
}

func adaptCommandKind(kind string) string {
	low := strings.ToLower
	switch low(kind) {
	case low(resource.KindService):
		return resource.KindService
	case low(resource.KindTenant):
		return resource.KindTenant
	case low(resource.KindLoadBalance):
		return resource.KindLoadBalance
	case low(resource.KindCanary):
		return resource.KindCanary
	case low(resource.KindObservabilityTracings):
		return resource.KindObservabilityTracings
	case low(resource.KindObservabilityOutputServer):
		return resource.KindObservabilityOutputServer
	case low(resource.KindObservabilityMetrics):
		return resource.KindObservabilityMetrics
	case low(resource.KindResilience):
		return resource.KindResilience
	case low(resource.KindIngress):
		return resource.KindIngress
	default:
		return kind
	}
}

func (v *CommandVisitor) Visit(fn VisitorFunc) error {
	vk := resource.VersionKind{
		APIVersion: resource.DefaultAPIVersion,
		Kind:       v.Kind,
	}

	var mo resource.MeshObject
	var err error

	if v.Name == "" {
		mo, err = v.oc.NewFromKind(vk)
	} else {
		mo, err = v.oc.NewFromResource(resource.MeshResource{
			VersionKind: vk,
			MetaData: resource.MetaData{
				Name: v.Name,
			},
		})
	}

	return fn(mo, err)
}
