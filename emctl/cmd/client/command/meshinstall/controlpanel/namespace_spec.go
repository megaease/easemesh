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

package controlpanel

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func namespaceSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:   ctx.Flags.MeshNamespace,
		Labels: map[string]string{},
	}}
	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployNamespace(ns, ctx.Client)
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
		return nil
	}
}
