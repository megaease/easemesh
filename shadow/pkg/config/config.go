package config

import (
	"path"

	"k8s.io/client-go/util/homedir"
)

const (
	DefaultKubeDir = ".kube"
	DefaultKubernetesConfig = "config"
)

var (
	DefaultKubernetesConfigDir  = path.Join(homedir.HomeDir(), DefaultKubeDir)
	DefaultKubernetesConfigPath = path.Join(DefaultKubernetesConfigDir, DefaultKubernetesConfig)
)