package meta

type (

	// VersionKind holds version and kind information for APIs
	VersionKind struct {
		APIVersion string `yaml:"apiVersion" yaml:"apiVersion" jsonschema:"omitempty"`
		Kind       string `yaml:"kind" yaml:"kind" jsonschema:"required"`
	}

	// MetaData is meta data for resources of the EaseMesh
	MetaData struct {
		Name   string            `yaml:"name" yaml:"name" jsonschema:"required"`
		Labels map[string]string `yaml:"labels,omitempty" yaml:"labels,omitempty" jsonschema:"omitempty"`
	}

	// MeshResource holds common information for a resource of the EaseMesh
	MeshResource struct {
		VersionKind `yaml:",inline" yaml:",inline"`
		MetaData    MetaData `yaml:"metadata" yaml:"metadata" jsonschema:"required"`
	}

	// MeshObject describes what's feature of a comman EaseMesh object
	MeshObject interface {
		Name() string
		Kind() string
		APIVersion() string
		Labels() map[string]string
	}
)

// Name returns name of the EaseMesh resource
func (m *MeshResource) Name() string {
	return m.MetaData.Name
}

// Kind returns kind of the EaseMesh resource
func (m *MeshResource) Kind() string {
	return m.VersionKind.Kind
}

// APIVersion returns api version of the EaseMesh resource
func (m *MeshResource) APIVersion() string {
	return m.VersionKind.APIVersion
}

// Labels returns labels of the EaseMesh resource
func (m *MeshResource) Labels() map[string]string {
	return m.MetaData.Labels
}
