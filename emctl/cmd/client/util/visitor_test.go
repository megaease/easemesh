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
package util

import (
	"testing"

	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	meshtesting "github.com/megaease/easemeshctl/cmd/client/testing"
)

func TestFileVisitor(t *testing.T) {

	types := meshtesting.GetAllResourceKinds()
	types = append(types, meshtesting.ResourceTypeKind{Type: nil, Kind: resource.KindServiceInstance})
	for i, tp := range types {
		name := "resource"
		if i == 0 {
			name = ""

		}
		newCommandVisitor(tp.Kind, name).
			Visit(func(mo meta.MeshObject, e error) error { return nil })
	}

}

func TestVisitorForSTDIN(t *testing.T) {
	FileVisitorForSTDIN(newDefaultDecoder()).Visit(func(mo meta.MeshObject, e error) error { return nil })
}
