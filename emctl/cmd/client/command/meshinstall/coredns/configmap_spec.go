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

package coredns

import (
	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func configMapSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	cfg := `.:53 {
	errors
	health {
		lameduck 5s
	}
	ready
	easemesh cluster.local {
		# endpoint could be retrieved by the plugin on its own.
		# endpoint http://{change_me_to_cluster_ip_of_easemesh_control_plane}:2379
	}
	kubernetes cluster.local in-addr.arpa ip6.arpa {
		pods insecure
		fallthrough in-addr.arpa ip6.arpa
	}
	prometheus :9153
	forward . /etc/resolv.conf {
		prefer_udp
	}
	cache 30
	loop
	reload
	loadbalance
}
	`
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      coreDNSConfigMap,
			Namespace: coreDNSNamespace,
			Annotations: map[string]string{
				"addonmanager.kubernetes.io/mode": "EnsureExists",
			},
		},
	}

	data := map[string]string{}
	data["Corefile"] = cfg
	configMap.Data = data

	return func(ctx *installbase.StageContext) error {
		err := installbase.DeployConfigMap(configMap, ctx.Client, coreDNSNamespace)
		if err != nil {
			return errors.Wrapf(err, "deploy ConfigMap %s failed", configMap.Name)
		}

		return nil
	}
}
