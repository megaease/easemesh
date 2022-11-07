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

package handler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func TestShadowServiceCloner_cloneDeploymentSpec(t *testing.T) {
	type fields struct {
		KubeClient kubernetes.Interface
	}

	deployment := fake.NewSourceDeployment()
	shadowService := fake.NewShadowService()
	clonedDeployment := fake.NewShadowDeployment()
	type args struct {
		sourceDeployment *appsV1.Deployment
		shadowService    *object.ShadowService
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsV1.Deployment
	}{
		{
			name:   "test",
			fields: fields{},
			args: args{
				sourceDeployment: deployment,
				shadowService:    &shadowService,
			},
			want: clonedDeployment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloner := &ShadowServiceCloner{
				KubeClient: tt.fields.KubeClient,
			}
			got := cloner.cloneDeploymentSpec(tt.args.sourceDeployment, tt.args.shadowService)

			buff, _ := yaml.Marshal(got)
			fmt.Printf("%s\n\n", buff)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cloneDeploymentSpec() = %v, \n want %v", got, tt.want)
			}
		})
	}
}
