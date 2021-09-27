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

package get

import (
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/megaease/easemeshctl/cmd/client/command/meshclient/fake"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"

	"github.com/spf13/cobra"
)

func TestDeleteRun(t *testing.T) {
	fakeExit := func(int) {
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()
	getFlag := meshtesting.PrepareGetFlags("__test_get_reactor", tenantSpec, t)
	fake.NewResourceReactorBuilder(getFlag.Server).
		AddReactor("*", "*", "*", func(action fake.Action) (handled bool, rets []meta.MeshObject, err error) {
			return true, nil, nil
		}).Added()
	cmd := &cobra.Command{}
	cmd.ParseFlags([]string{"tenant", "mesh-service"})
	Run(cmd, getFlag)
}

var tenantSpec = `
kind: Tenant
apiVersion: mesh.megaease.com/v1alpha1
metadata:
  name: mesh-service
spec:
  description: 'award tenant'
`
