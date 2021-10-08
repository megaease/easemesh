module github.com/megaease/easemesh/mesh-shadow

go 1.16

require (
	github.com/megaease/easemesh/mesh-operator v1.2.0
	github.com/megaease/easemeshctl v1.2.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.20.9
	k8s.io/apiextensions-apiserver v0.20.7 // indirect
	k8s.io/apimachinery v0.20.9
	k8s.io/client-go v0.20.9
	sigs.k8s.io/controller-runtime v0.7.2
)

replace github.com/megaease/easemeshctl => ../emctl/

replace github.com/megaease/easemesh/mesh-operator => ../operator/
