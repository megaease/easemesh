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
package apply

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"

	"bou.ke/monkey"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	utiltesting "k8s.io/client-go/util/testing"
)

func prepareFlags() (flag *flags.Apply, err error) {

	specDir, err := utiltesting.MkTmpdir("specdir")
	if err != nil {
		return nil, err
	}

	specFile := filepath.Join(specDir, "01-tenant.yaml")
	err = ioutil.WriteFile(specFile, []byte(tenantSpec), 0600)
	if err != nil {
		return nil, err
	}

	return &flags.Apply{
		AdminGlobal: &flags.AdminGlobal{
			Server: "__apply_test_reactor",
		},
		AdminFileInput: &flags.AdminFileInput{
			YamlFile: specFile,
		},
	}, nil
}

func TestRun(t *testing.T) {
	flag, err := prepareFlags()
	if err != nil {
		t.Fatalf("prepare Flags error:%s.", err)
	}

	fake.NewResourceReactorBuilder(flag.Server).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			return true, nil, nil
		}).
		Added()

	cmd := &cobra.Command{}
	Run(cmd, flag)
}

func TestRunFail(t *testing.T) {
	fakeExit := func(int) {
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	flag, err := prepareFlags()
	if err != nil {
		t.Fatalf("prepare Flags error:%s.", err)
	}

	fake.NewResourceReactorBuilder(flag.Server).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			return true, nil, errors.Errorf("mock an error")
		}).
		Added()

	cmd := &cobra.Command{}
	Run(cmd, flag)

	flag.Server = ""
	Run(cmd, flag)

	flag.Server = "placehold"
	flag.YamlFile = ""
	Run(cmd, flag)
}

var tenantSpec = `
kind: Tenant
apiVersion: mesh.megaease.com/v1alpha1
metadata:
  name: mesh-service
spec:
  description: 'award tenant'
`
