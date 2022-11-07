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
	"reflect"
	"sync"
	"testing"

	shadowfake "github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appsV1 "k8s.io/api/apps/v1"
)

func Test_isShadowDeployment(t *testing.T) {
	deployment1 := shadowfake.NewSourceDeployment()
	deployment2 := shadowfake.NewShadowDeployment()

	type args struct {
		spec appsV1.DeploymentSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				spec: deployment1.Spec,
			},
			want: false,
		},
		{
			name: "test2",
			args: args{
				spec: deployment2.Spec,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isShadowDeployment(tt.args.spec); got != tt.want {
				t.Errorf("isShadowDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceDeploySearcher_Search(t *testing.T) {
	searchChan := make(chan interface{})
	defer close(searchChan)

	searcher := &ShadowServiceDeploySearcher{
		KubeClient: prepareClientForTest(),
		CloneChan:  searchChan,
	}

	sourceDeployment := shadowfake.NewSourceDeployment()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case obj := <-searcher.CloneChan:
				if !reflect.DeepEqual(obj.(ShadowServiceBlock).deployment, sourceDeployment) {
					t.Errorf("Search Deployment Error, Searcher.Search() = %v, \n want %v", obj, sourceDeployment)
				}
				return
			}
		}
	}()

	shadowService := shadowfake.NewShadowService()
	objs := []object.ShadowService{shadowService}
	searcher.Search(objs)
	wg.Wait()
}
