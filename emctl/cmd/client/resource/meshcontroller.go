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

package resource

import (
	"fmt"

	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
)

type (

	// MeshController is the spec of MeshController on Easegress.
	MeshController struct {
		meta.MeshResource   `yaml:",inline"`
		MeshControllerAdmin `yaml:",inline"`
	}

	// MeshControllerV1Alpha1 is the v1alphv1 version of mesh controller.
	MeshControllerV1Alpha1 struct {
		Kind                string `yaml:"kind"`
		Name                string `yaml:"name"`
		MeshControllerAdmin `yaml:",inline"`
	}

	// MeshControllerAdmin is the admin config of mesh controller.
	MeshControllerAdmin struct {
		// HeartbeatInterval is the interval for one service instance reporting its heartbeat.
		HeartbeatInterval string `yaml:"heartbeatInterval" jsonschema:"required,format=duration"`

		// RegistryTime indicates which protocol the registry center accepts.
		RegistryType string `yaml:"registryType" jsonschema:"required"`

		// APIPort is the port for worker's API server
		APIPort int `yaml:"apiPort" jsonschema:"required"`

		// IngressPort is the port for http server in mesh ingress
		IngressPort int `yaml:"ingressPort" jsonschema:"required"`

		// ExternalServiceRegistry is the external service registry name.
		ExternalServiceRegistry string `yaml:"externalServiceRegistry" jsonschema:"omitempty"`
		CleanExternalRegistry   bool   `yaml:"cleanExternalRegistry"`

		Security *Security `yaml:"security" jsonschema:"omitempty"`

		// Sidecar injection relevant config.
		ImageRegistryURL          string `yaml:"imageRegistryURL" jsonschema:"omitempty"`
		ImagePullPolicy           string `yaml:"imagePullPolicy" jsonschema:"omitempty"`
		SidecarImageName          string `yaml:"sidecarImageName" jsonschema:"omitempty"`
		AgentInitializerImageName string `yaml:"agentInitializerImageName" jsonschema:"omitempty"`
		Log4jConfigName           string `yaml:"log4jConfigName" jsonschema:"omitempty"`

		MonitorMTLS *MonitorMTLS `yaml:"monitorMTLS" jsonschema:"omitempty"`
	}

	// Security is the spec for mesh-wide security.
	Security struct {
		MTLSMode     string `yaml:"mtlsMode" jsonschema:"required"`
		CertProvider string `yaml:"certProvider" jsonschema:"required"`

		RootCertTTL string `yaml:"rootCertTTL" jsonschema:"required,format=duration"`
		AppCertTTL  string `yaml:"appCertTTL" jsonschema:"required,format=duration"`
	}

	// MonitorMTLS is the spec of mTLS specification of monitor.
	MonitorMTLS struct {
		Enabled  bool   `yaml:"enabled" jsonschema:"required"`
		URL      string `yaml:"url" jsonschema:"required"`
		Username string `yaml:"username" jsonschema:"required"`
		Password string `yaml:"password" jsonschema:"required"`

		ReporterAppendType string         `yaml:"reporterAppendType"`
		CaCertBase64       string         `yaml:"caCertBase64" jsonschema:"required,format=base64"`
		Certs              []*MonitorCert `yaml:"certs" jsonschema:"required"`
	}

	// MonitorCert is the spec for single pack of mTLS.
	MonitorCert struct {
		CertBase64 string   `yaml:"certBase64" jsonschema:"required,format=base64"`
		KeyBase64  string   `yaml:"keyBase64" jsonschema:"required,format=base64"`
		Services   []string `yaml:"services" jsonschema:"required"`
	}
)

var _ meta.TableObject = &MeshController{}

// Columns returns the columns of MeshController.
func (mc *MeshController) Columns() []*meta.TableColumn {
	ports := fmt.Sprintf("%d/API,%d/Ingress", mc.APIPort, mc.IngressPort)

	return []*meta.TableColumn{
		{
			Name:  "Heartbeat",
			Value: mc.HeartbeatInterval,
		},
		{
			Name:  "Registry",
			Value: mc.RegistryType,
		},
		{
			Name:  "Ports",
			Value: ports,
		},
		{
			Name:  "ExternalRegistry",
			Value: mc.ExternalServiceRegistry,
		},
	}
}

// ToV1Alpha1 converts MeshController resouce to v1alpha1.
func (mc *MeshController) ToV1Alpha1() *MeshControllerV1Alpha1 {
	return &MeshControllerV1Alpha1{
		Kind:                mc.Kind(),
		Name:                mc.Name(),
		MeshControllerAdmin: mc.MeshControllerAdmin,
	}
}

// ToMeshController converts a MeshControllerV1Alpha1 resouce to a MeshController resource.
func ToMeshController(meshController *MeshControllerV1Alpha1) *MeshController {
	return &MeshController{
		MeshResource:        NewMeshResource(DefaultAPIVersion, meshController.Kind, meshController.Name),
		MeshControllerAdmin: meshController.MeshControllerAdmin,
	}
}
