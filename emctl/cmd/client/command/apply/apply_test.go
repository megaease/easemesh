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
	"os"
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"

	"bou.ke/monkey"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func TestRun(t *testing.T) {
	flag := meshtesting.PrepareApplyFlags("__test_apply_reactor", tenantSpec, t)

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

	flag := meshtesting.PrepareApplyFlags("__test_apply_reactor", tenantSpec, t)

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
