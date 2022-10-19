module github.com/megaease/easemesh/mesh-shadow

go 1.16

require (
	github.com/megaease/easemesh-api v1.4.4
	github.com/megaease/easemeshctl v1.2.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.3
	sigs.k8s.io/controller-runtime v0.10.3
	sigs.k8s.io/yaml v1.2.0
)

// Fix security problems: https://www.oscs1024.com/cd/1530049682318094336?sign=9164b062

replace github.com/dgrijalva/jwt-go v3.2.0+incompatible => github.com/golang-jwt/jwt/v4 v4.4.2

replace github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.5.0

replace github.com/miekg/dns v1.0.14 => github.com/miekg/dns v1.1.50

//

replace github.com/megaease/easemeshctl => ../../emctl/

replace github.com/megaease/easemesh/mesh-operator => ../../operator/
