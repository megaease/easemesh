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

package printer

import (
	"fmt"
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"
)

func TestPrinter(t *testing.T) {

	yamlPrinter := New("yaml")
	jsonPrinter := New("json")
	tablePrinter := New("table")

	for _, rtk := range meshtesting.GetAllResourceKinds() {
		fmt.Printf("%+v", rtk)
		obj := meshtesting.CreateMeshObjectFromType(rtk.Type, rtk.Kind, "obj")

		yamlPrinter.PrintObjects([]meta.MeshObject{obj})
		jsonPrinter.PrintObjects([]meta.MeshObject{obj})
		tablePrinter.PrintObjects([]meta.MeshObject{obj})

	}
}
