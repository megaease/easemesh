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

package main

import (
	"flag"
	"log"
	"time"

	"github.com/megaease/easemesh/mesh-shadow/pkg/controller"
	// load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	DefaultMeshServer = "easemesh-controlplane-svc.easemesh:2381"
)

var (
	meshServer = flag.String("mesh-server", DefaultMeshServer, "An address to access the EaseMesh control plane")
)

func easemeshOption(config *controller.Config) error {
	config.MeshServer = *meshServer
	config.RequestTimeout = 10 * time.Second
	config.PullInterval = 1 * time.Minute
	return nil
}

func main() {
	flag.Parse()
	controller, err := controller.NewShadowServiceController(easemeshOption)
	if err != nil {
		log.Fatalf("new collector service error: %s", err)
		return
	}

	<-controller.Do()
}
