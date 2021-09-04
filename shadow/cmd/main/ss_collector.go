package main

import (
	"flag"
	"log"

	"github.com/megaease/easemesh/mesh-shadow/cmd/main/rcfile"
	"github.com/megaease/easemesh/mesh-shadow/pkg/common"
	"github.com/megaease/easemesh/mesh-shadow/pkg/flow"

	// load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	meshServer = flag.String("mesh-server", "", "An address to access the EaseMesh control plane")
)

func easemeshOption(config *flow.ServiceConfig) error {
	config.MeshServer = *meshServer
	if config.MeshServer == "" {
		config.MeshServer = GetServerAddress()
	}
	return nil
}

func GetServerAddress() string {
	rc, err := rcfile.New()
	if err != nil {
		return ""
	}

	err = rc.Unmarshal()
	if err != nil {
		common.OutputErrorf("unmarshal rcfile failed: %v", err)
		return ""
	}
	return rc.Server
}

func main() {

	flag.Parse()

	service, err := flow.New(easemeshOption)
	if err != nil {
		log.Fatalf("new collector service error: %s", err)
		return
	}

	<-service.Do()
}
