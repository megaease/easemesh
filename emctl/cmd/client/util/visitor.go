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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"

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

// Visitor is visitor to visit all MeshObjects via VisitorFunc
type Visitor interface {
	Visit(VisitorFunc) error
}

// VisitorFunc executes visition logic
type VisitorFunc func(meta.MeshObject, error) error

// RawExtension is a raw struct that holds raw information of the spec
type RawExtension struct {
	Raw []byte `json:"-" protobuf:"bytes,1,opt,name=raw"`
}

// UnmarshalJSON unmarshal byte to RawExtension
func (re *RawExtension) UnmarshalJSON(in []byte) error {
	if re == nil {
		return errors.New("runtime.RawExtension: UnmarshalJSON on nil pointer")
	}
	if !bytes.Equal(in, []byte("null")) {
		re.Raw = append(re.Raw[0:0], in...)
	}
	return nil
}

// MarshalJSON marshal RawExtension to bytes
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
	return &fileVisitor{
		Path:          constSTDINstr,
		streamVisitor: newStreamVisitor(nil, decoder, constSTDINstr),
	}
}

func expandPathsToFileVisitors(decoder Decoder, paths string, recursive bool, extensions []string) ([]Visitor, error) {
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

		visitor := &fileVisitor{
			Path:          path,
			streamVisitor: newStreamVisitor(nil, decoder, path),
		}

		visitors = append(visitors, visitor)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return visitors, nil
}

// fileVisitor is wrapping around a StreamVisitor, to handle open/close files
type fileVisitor struct {
	Path string
	*streamVisitor
}

var _ Visitor = &fileVisitor{}

// Visit in a FileVisitor is just taking care of opening/closing files
func (v *fileVisitor) Visit(fn VisitorFunc) error {
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
	v.streamVisitor.Reader = transform.NewReader(f, utf16bom)

	return v.streamVisitor.Visit(fn)
}

type streamVisitor struct {
	io.Reader

	Decoder Decoder
	Source  string
}

var _ Visitor = &streamVisitor{}

// newStreamVisitor is a helper function that is useful when we want to change the fields of the struct but keep calls the same.
func newStreamVisitor(r io.Reader, decoder Decoder, source string) *streamVisitor {
	return &streamVisitor{
		Reader:  r,
		Decoder: decoder,
		Source:  source,
	}
}

// Visit implements Visitor over a stream. StreamVisitor is able to distinct multiple resources in one stream.
func (v *streamVisitor) Visit(fn VisitorFunc) error {
	var errs []error
	d := yaml.NewYAMLOrJSONDecoder(v.Reader, 4096)
	for {
		ext := RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err != io.EOF {
				errs = append(errs, errors.Errorf("error parsing %s: %v", v.Source, err))
			}
			break
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

		err1 := fn(info, err)
		if err1 != nil {
			errs = append(errs, err1)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	var finalErr error
	for _, err := range errs {
		if finalErr == nil {
			finalErr = fmt.Errorf("%v", err)
		} else {
			finalErr = fmt.Errorf("%v\n%v", finalErr, err)
		}
	}

	return finalErr
}

func (v *streamVisitor) decodeMeshObject(data []byte, source string) (meta.MeshObject, error) {
	meshObject, _, err := v.Decoder.Decode(data)
	if err != nil {
		return nil, err
	}
	return meshObject, nil
}

type urlVisitor struct {
	URL *url.URL
	*streamVisitor
	HTTPAttemptCount int
}

func (v *urlVisitor) Visit(fn VisitorFunc) error {
	body, err := readHTTPWithRetries(resty.New(), 5*time.Second, v.URL.String(), v.HTTPAttemptCount)
	if err != nil {
		return err
	}
	defer body.Close()
	v.streamVisitor.Reader = body
	return v.streamVisitor.Visit(fn)
}

// readHTTPWithRetries tries to http.Get the v.URL retries times before giving up.
func readHTTPWithRetries(client *resty.Client, duration time.Duration, u string, attempts int) (io.ReadCloser, error) {
	r, err := client.
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
		SetDoNotParseResponse(true).
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

type commandVisitor struct {
	Kind string
	Name string

	oc resource.ObjectCreator
}

var _ Visitor = &commandVisitor{}

func newCommandVisitor(kind, name string) *commandVisitor {
	return &commandVisitor{
		Kind: adaptCommandKind(kind),
		Name: name,
		oc:   resource.NewObjectCreator(),
	}
}

func adaptCommandKind(kind string) string {
	low := strings.ToLower
	switch low(kind) {
	case low(resource.KindMeshController):
		return resource.KindMeshController
	case low(resource.KindService):
		return resource.KindService
	case low(resource.KindServiceInstance):
		return resource.KindServiceInstance
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

func (v *commandVisitor) Visit(fn VisitorFunc) error {
	vk := meta.VersionKind{
		APIVersion: resource.DefaultAPIVersion,
		Kind:       v.Kind,
	}

	var mo meta.MeshObject
	var err error

	if v.Name == "" {
		mo, err = v.oc.NewFromKind(vk)
	} else {
		mo, err = v.oc.NewFromResource(meta.MeshResource{
			VersionKind: vk,
			MetaData: meta.MetaData{
				Name: v.Name,
			},
		})
	}

	return fn(mo, err)
}
