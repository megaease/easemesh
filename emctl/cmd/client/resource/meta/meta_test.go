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
package meta

import "testing"

func TestMeta(t *testing.T) {
	const (
		expectKind       = "Tenant"
		expectAPIVersion = "mesh.megaease.com/v1alpha1"
		expectName       = "pet"
	)

	mr := MeshResource{
		VersionKind{
			APIVersion: expectAPIVersion,
			Kind:       expectKind,
		},
		MetaData{
			Name:   expectName,
			Labels: nil,
		},
	}

	if mr.Kind() != expectKind {
		t.Fatalf("expect kind is %s but %s", expectKind, mr.Kind())
	}

	if mr.APIVersion() != expectAPIVersion {
		t.Fatalf("expect APIVersion is %s but %s", expectAPIVersion, mr.APIVersion())
	}

	if mr.Name() != expectName {
		t.Fatalf("expect name is %s but %s", expectName, mr.Name())
	}

	if mr.Labels() != nil {
		t.Fatalf("expect labels is nil but %+v", mr.Labels())
	}
}
