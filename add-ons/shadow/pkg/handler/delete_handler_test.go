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
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func prepareClientForTest() kubernetes.Interface {
	var result runtime.Object
	namespace := shadowfake.NewNamespace()
	deployment := shadowfake.NewSourceDeployment()
	shadowDeployment := shadowfake.NewShadowDeployment()

	client := fake.NewSimpleClientset(
		namespace,
		deployment,
		shadowDeployment,
	)
	client.PrependReactor("create", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		result = action.(k8stesting.CreateAction).GetObject()

		return true, action.(k8stesting.CreateAction).GetObject(), k8serr.NewAlreadyExists(schema.GroupResource{
			Resource: "Namespace",
			Group:    "v1",
		}, "na")
	})

	client.PrependReactor("update", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, action.(k8stesting.UpdateAction).GetObject(), nil
	})

	client.PrependReactor("get", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, result, nil
	})

	return client
}

func Test_namespacedName(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				name:      "test1",
				namespace: "testns",
			},
			want: "testns" + "/" + "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := namespacedName(tt.args.namespace, tt.args.name); got != tt.want {
				t.Errorf("namespacedName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shadowServiceExists(t *testing.T) {
	type args struct {
		namespacedName       string
		shadowServiceNameMap map[string]object.ShadowService
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				namespacedName: "testns/test1",
				shadowServiceNameMap: map[string]object.ShadowService{
					"testns/test1": {
						Name:      "test1",
						Namespace: "testns",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shadowServiceExists(tt.args.namespacedName, tt.args.shadowServiceNameMap); got != tt.want {
				t.Errorf("shadowServiceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceDeleter_Delete(t *testing.T) {
	deleteChan := make(chan interface{})
	defer close(deleteChan)

	deleter := &ShadowServiceDeleter{
		KubeClient: prepareClientForTest(),
		DeleteChan: deleteChan,
	}

	clonedDeployment := shadowfake.NewShadowDeployment()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case obj := <-deleter.DeleteChan:
				block := obj.(ShadowServiceBlock)
				if !reflect.DeepEqual(block.deployObj, *clonedDeployment) {
					t.Errorf("FindDeletableObjs() = %v, \n want %v", obj, clonedDeployment)
				}
				deleter.Delete(obj)
				return
			}
		}
	}()

	var objs []object.ShadowService
	deleter.FindDeletableObjs(objs)
	wg.Wait()
}
