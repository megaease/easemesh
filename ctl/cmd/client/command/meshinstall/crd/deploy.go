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

package crd

import (
	_ "embed"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	"github.com/pkg/errors"
	apiExtensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientGoScheme "k8s.io/client-go/kubernetes/scheme"
)

//go:embed  crd.yaml
var easemeshDeploymentCRD []byte

// Deploy deploy resources of crd
func Deploy(context *installbase.StageContext) error {
	crd, err := getCRDSpec(easemeshDeploymentCRD)
	if err != nil {
		return err
	}

	err = installbase.DeployCustomResourceDefinition(crd, context.APIExtensionsClient)
	if err != nil {
		return errors.Wrapf(err, "can't deploy CRD %s", crd.Name)
	}
	return err
}

// PreCheck check prerequisite for installing CRD
func PreCheck(context *installbase.StageContext) error {
	return nil
}

// Clear will clear all installed resource about control panel
func Clear(context *installbase.StageContext) error {
	crd, err := getCRDSpec(easemeshDeploymentCRD)
	if err != nil {
		return err
	}
	return installbase.DeleteCRDResource(context.APIExtensionsClient, crd.Name)
}

// Describe leverage human-readable text to describe different phase
// in the process of the CRD installation
func Describe(context *installbase.StageContext, phase installbase.InstallPhase) string {
	switch phase {
	case installbase.BeginPhase:
		return "Begin to deploy CRD meshdeployment\n"
	case installbase.EndPhase:
		return "CustomeResourceDefine meshdeployment deployed successfully\n"
	}
	return ""
}

func getCRDSpec(yaml []byte) (*apiExtensionsV1.CustomResourceDefinition, error) {
	var err error
	sch := runtime.NewScheme()
	_ = clientGoScheme.AddToScheme(sch)
	_ = apiExtensionsV1.AddToScheme(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, kind, err := decode(yaml, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode Yaml string")
	}

	crd, ok := obj.(*apiExtensionsV1.CustomResourceDefinition)
	if !ok {
		return nil, errors.Errorf("can't convert custom resource to apiExtensionsV1.CustomResourceDefinition object, kind is %s", kind.Kind)
	}

	return crd, nil
}
