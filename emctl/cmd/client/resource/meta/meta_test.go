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
