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

package base

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// Runtime carries base rutime for one controller.
	Runtime struct {
		Name             string
		Client           client.Client
		Scheme           *runtime.Scheme
		Recorder         record.EventRecorder
		Log              logr.Logger
		ImageRegistryURL string
		ImagePullPolicy  string
		SidecarImageName string
		// AgentInitializerImageName is the image name of the Agent initializer.
		AgentInitializerImageName string
		// Log4jConfigName is  the name of log4f config name.
		Log4jConfigName string

		APIAddr         string
		ClusterJoinURLs []string
		ClusterName     string
	}
)
