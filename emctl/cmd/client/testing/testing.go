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

package testing

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"

	"github.com/spf13/cobra"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	utiltesting "k8s.io/client-go/util/testing"
)

// ResourceTypeKind is pair of the reflecting type and kind of the resource.
type ResourceTypeKind struct {
	Type reflect.Type
	Kind string
}

// GetAllResourceKinds initials a default resource kind pair dimmention.
func GetAllResourceKinds() []ResourceTypeKind {

	return []ResourceTypeKind{
		{Type: reflect.TypeOf(resource.Tenant{}), Kind: resource.KindTenant},
		{Type: reflect.TypeOf(resource.MeshController{}), Kind: resource.KindMeshController},
		{Type: reflect.TypeOf(resource.Ingress{}), Kind: resource.KindIngress},
		{Type: reflect.TypeOf(resource.CustomResourceKind{}), Kind: resource.KindCustomResourceKind},
		{Type: reflect.TypeOf(resource.CustomResource{}), Kind: "-"},
		{Type: reflect.TypeOf(resource.LoadBalance{}), Kind: resource.KindLoadBalance},
		{Type: reflect.TypeOf(resource.ObservabilityMetrics{}), Kind: resource.KindObservabilityMetrics},
		{Type: reflect.TypeOf(resource.ObservabilityOutputServer{}), Kind: resource.KindObservabilityOutputServer},
		{Type: reflect.TypeOf(resource.ObservabilityTracings{}), Kind: resource.KindObservabilityTracings},
		{Type: reflect.TypeOf(resource.Canary{}), Kind: resource.KindCanary},
		{Type: reflect.TypeOf(resource.Service{}), Kind: resource.KindService},
		{Type: reflect.TypeOf(resource.Resilience{}), Kind: resource.KindResilience},
	}
}

// CreateMeshObjectFromType constructs a mesh object via reflecting with type and base information.
func CreateMeshObjectFromType(t reflect.Type, kind, nm string) meta.MeshObject {
	meshObject := reflect.New(t).
		Elem() // reflect.Value

	versionKind := meshObject.FieldByName("VersionKind")
	version := versionKind.FieldByName("APIVersion")
	knd := versionKind.FieldByName("Kind")

	knd.SetString(kind)
	version.SetString("v1alpha1")

	metaData := meshObject.FieldByName("MetaData")
	name := metaData.FieldByName("Name")
	name.SetString(nm)
	return meshObject.Addr().Interface().(meta.MeshObject)
}

func prepareAdminGlobal(server string) *flags.AdminGlobal {
	return &flags.AdminGlobal{
		Server: server,
	}
}

// PrepareYamlFile prepare a temporary spec file for using
func PrepareYamlFile(spec string, t *testing.T) (specFile string) {
	specDir, err := utiltesting.MkTmpdir("specdir")
	if err != nil {
		t.Fatalf("mkdir tmpdir %s error:%s", specDir, err)
	}

	specFile = filepath.Join(specDir, "01-tenant.yaml")
	err = ioutil.WriteFile(specFile, []byte(spec), 0600)
	if err != nil {
		t.Fatalf("write %s file error:%s", specFile, err)
	}
	return specFile

}

func prepareFileInput(spec string, t *testing.T) *flags.AdminFileInput {
	specFile := PrepareYamlFile(spec, t)
	return &flags.AdminFileInput{
		YamlFile: specFile,
	}
}

// PrepareApplyFlags return a mock Apply flag
func PrepareApplyFlags(server, spec string, t *testing.T) *flags.Apply {
	return &flags.Apply{AdminGlobal: prepareAdminGlobal(server), AdminFileInput: prepareFileInput(spec, t)}
}

// PrepareDeleteFlags return a mock Apply flag
func PrepareDeleteFlags(server, spec string, t *testing.T) *flags.Delete {
	return &flags.Delete{AdminGlobal: prepareAdminGlobal(server), AdminFileInput: prepareFileInput(spec, t)}
}

// PrepareGetFlags return a mock Get flag
func PrepareGetFlags(server, spec string, t *testing.T) *flags.Get {
	return &flags.Get{AdminGlobal: prepareAdminGlobal(server), OutputFormat: "yaml"}
}

// PrepareInstallContext return a StageContext of install
func PrepareInstallContext(cmd *cobra.Command,
	client kubernetes.Interface,
	extensionClient apiextensions.Interface,
	installFlags *flags.Install) *installbase.StageContext {
	clearFunc := func(*installbase.StageContext) error { return nil }
	return &installbase.StageContext{
		Cmd:                 cmd,
		Client:              client,
		APIExtensionsClient: extensionClient,
		ClearFuncs:          []func(*installbase.StageContext) error{clearFunc},
		Flags:               installFlags,
	}
}
