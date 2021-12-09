module github.com/megaease/easemesh/mesh-shadow

go 1.16

require (
	github.com/megaease/easemesh-api v1.3.4
	github.com/megaease/easemeshctl v1.2.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.3
	sigs.k8s.io/controller-runtime v0.10.3
)

replace github.com/megaease/easemeshctl => ../../emctl/

replace github.com/megaease/easemesh/mesh-operator => ../../operator/
