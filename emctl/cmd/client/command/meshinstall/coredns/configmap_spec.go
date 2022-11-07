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
	"context"
	"fmt"

	installbase "github.com/megaease/easemeshctl/cmd/client/command/meshinstall/base"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	// NOTE: https://github.com/kubernetes/kubernetes/blob/master/cmd/kubeadm/app/apis/kubeadm/types.go
	// Don't import kubeadm package, because it will cause a lot of dependencies.

	// ClusterConfiguration contains cluster-wide configuration for a kubeadm cluster
	ClusterConfiguration struct {
		// Networking holds configuration for the networking topology of the cluster.
		Networking Networking `yaml:"networking"`
	}

	// Networking contains elements describing cluster's networking configuration.
	Networking struct {
		// ServiceSubnet is the subnet used by k8s services. Defaults to "10.96.0.0/12".
		ServiceSubnet string `yaml:"serviceSubnet"`
		// PodSubnet is the subnet used by pods.
		PodSubnet string `yaml:"podSubnet"`
		// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
		DNSDomain string `yaml:"dnsDomain"`
	}
)

func configMapSpec(ctx *installbase.StageContext) installbase.InstallFunc {
	return func(ctx *installbase.StageContext) error {
		var dnsDomain string
		if ctx.CoreDNSFlags.DNSDomain != "" {
			dnsDomain = ctx.CoreDNSFlags.DNSDomain
		} else {
			// ClusterConfiguration.networking.dnsDomain in ConfigMap kube-system/kubeadm-config
			kubeConfigMap, err := ctx.Client.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "kubeadm-config", metav1.GetOptions{})
			if err != nil {
				return errors.Wrap(err, "failed to get kubeadm-config in order to set DNS domain, you can use --dns-domain to specify the DNS domain")
			}

			clusterBuff := kubeConfigMap.Data["ClusterConfiguration"]
			clusterConfig := ClusterConfiguration{}
			err = yaml.Unmarshal([]byte(clusterBuff), &clusterConfig)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal kubeadm-config in order to set DNS domain, you can use --dns-domain to specify the DNS domain")
			}

			dnsDomain = clusterConfig.Networking.DNSDomain
		}

		cfg := fmt.Sprintf(
			`.:53 {
	log
	errors
	health {
		lameduck 5s
	}
	ready
	easemesh %s {
		# endpoint could be retrieved by the plugin on its own.
		# endpoint http://{change_me_to_cluster_ip_of_easemesh_control_plane}:2379
	}
	kubernetes %s in-addr.arpa ip6.arpa {
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
`, dnsDomain, dnsDomain)

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

		err := installbase.DeployConfigMap(configMap, ctx.Client, coreDNSNamespace)
		if err != nil {
			return errors.Wrapf(err, "deploy ConfigMap %s failed", configMap.Name)
		}

		return nil
	}
}
