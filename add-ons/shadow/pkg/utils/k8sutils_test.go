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

package utils

import (
	"testing"

	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewKubernetesClient(t *testing.T) {
	NewKubernetesClient()
}

func prepareClientForTest() kubernetes.Interface {
	var result runtime.Object
	client := fake.NewSimpleClientset()
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

	client.PrependReactor("delete", "*", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, result, nil
	})

	return client
}

func TestDeleteDeployment(t *testing.T) {
	client := prepareClientForTest()
	err := DeleteDeployment("test", "test", client, metav1.DeleteOptions{})
	if err != nil {
		t.Fatalf("delete deploymet error: %s", err)
	}
}

func TestListDeployments(t *testing.T) {
	client := prepareClientForTest()
	_, err := ListDeployments("test", client, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("list namespace error: %s", err)
	}
}

func TestListNameSpaces(t *testing.T) {
	client := prepareClientForTest()
	_, err := ListNameSpaces(client)
	if err != nil {
		t.Fatalf("list namespace error: %s", err)
	}
}
