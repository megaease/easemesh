package installbase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/megaease/easemeshctl/cmd/client/command/flags"
	"github.com/megaease/easemeshctl/cmd/common"

	"github.com/spf13/cobra"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

func NewKubernetesApiExtensionsClient() (*apiextensions.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", DefaultKubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := apiextensions.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func CreateNameSpace(namespace *v1.Namespace, clientSet *kubernetes.Clientset) error {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), namespace.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		_, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil && errors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func applyResource(createFunc func() error, updateFunc func() error) error {
	err := createFunc()
	if err != nil && errors.IsAlreadyExists(err) {
		err = updateFunc()
	}
	return err
}

func DeployDeployment(deployment *appsV1.Deployment, clientSet *kubernetes.Clientset, namespaces string) error {
	return applyResource(
		func() error {
			_, err := clientSet.AppsV1().Deployments(namespaces).Create(context.TODO(), deployment, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.AppsV1().Deployments(namespaces).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return err
		})
}

func DeployStatefulset(statefulSet *appsV1.StatefulSet, clientSet *kubernetes.Clientset, namespaces string) error {
	return applyResource(
		func() error {
			_, err := clientSet.AppsV1().StatefulSets(namespaces).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.AppsV1().StatefulSets(namespaces).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
			return err
		},
	)
}

func DeployService(service *v1.Service, clientSet *kubernetes.Clientset, namespaces string) error {
	return applyResource(
		func() error {
			_, err := clientSet.CoreV1().Services(namespaces).Create(context.TODO(), service, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.CoreV1().Services(namespaces).Update(context.TODO(), service, metav1.UpdateOptions{})
			return err
		},
	)
}

func DeployConfigMap(configMap *v1.ConfigMap, clientSet *kubernetes.Clientset, namespaces string) error {
	return applyResource(
		func() error {
			_, err := clientSet.CoreV1().ConfigMaps(namespaces).Create(context.TODO(), configMap, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.CoreV1().ConfigMaps(namespaces).Update(context.TODO(), configMap, metav1.UpdateOptions{})
			return err
		})
}

func ListPersistentVolume(clientSet *kubernetes.Clientset) (*v1.PersistentVolumeList, error) {
	return clientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
}

func DeployRole(role *rbacv1.Role, clientSet *kubernetes.Clientset, namespaces string) error {
	return applyResource(
		func() error {
			_, err := clientSet.RbacV1().Roles(namespaces).Create(context.TODO(), role, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.RbacV1().Roles(namespaces).Update(context.TODO(), role, metav1.UpdateOptions{})
			return err
		})
}

func DeployRoleBinding(roleBinding *rbacv1.RoleBinding, clientSet *kubernetes.Clientset, namespaces string) error {
	_, err := clientSet.RbacV1().RoleBindings(namespaces).Get(context.TODO(), roleBinding.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		_, err = clientSet.RbacV1().RoleBindings(namespaces).Create(context.TODO(), roleBinding, metav1.CreateOptions{})
	} else {
		_, err = clientSet.RbacV1().RoleBindings(namespaces).Update(context.TODO(), roleBinding, metav1.UpdateOptions{})
	}
	return err
}

func DeployClusterRole(clusterRole *rbacv1.ClusterRole, clientSet *kubernetes.Clientset) error {
	return applyResource(
		func() error {
			_, err := clientSet.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.RbacV1().ClusterRoles().Update(context.TODO(), clusterRole, metav1.UpdateOptions{})
			return err
		})
}

func DeployClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding, clientSet *kubernetes.Clientset) error {
	return applyResource(
		func() error {
			_, err := clientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.RbacV1().ClusterRoleBindings().Update(context.TODO(), clusterRoleBinding, metav1.UpdateOptions{})
			return err
		})
}

func DeployCustomResourceDefinition(crd *apiextensionsv1.CustomResourceDefinition, clientSet *apiextensions.Clientset) error {
	return applyResource(
		func() error {
			_, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})
			return err
		},
		func() error {
			_, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Update(context.TODO(), crd, metav1.UpdateOptions{})
			return err
		})
}

type PredictFunc func(interface{}) bool

func StatefulsetReadyPredict(object interface{}) (ready bool) {
	statefulset, ok := object.(*appsV1.StatefulSet)
	if !ok {
		return
	}
	return statefulset.Status.ReadyReplicas == *statefulset.Spec.Replicas
}
func DeploymentReadyPredict(object interface{}) (ready bool) {
	deploy, ok := object.(*appsV1.Deployment)
	if !ok {
		return
	}
	return deploy.Status.ReadyReplicas == *deploy.Spec.Replicas
}
func CheckStatefulsetResourceStatus(client *kubernetes.Clientset, namespace, resourceName string, predict PredictFunc) (bool, error) {
	statefulset, err := client.AppsV1().StatefulSets(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return predict(statefulset), nil
}

func CheckDeploymentResourceStatus(client *kubernetes.Clientset, namespace, name string, predict PredictFunc) (bool, error) {

	deploy, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return predict(deploy), nil
}

func GetMeshControlPanelEntryPoints(client *kubernetes.Clientset, namespace, resourceName, portName string) ([]string, error) {
	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodePort int32
	for _, p := range service.Spec.Ports {
		if p.Name == portName {
			nodePort = p.NodePort
			break
		}
	}
	entrypoints := []string{}
	for _, n := range nodes.Items {
		for _, i := range n.Status.Addresses {
			if i.Type == v1.NodeInternalIP {
				entrypoints = append(entrypoints, "http://"+i.Address+":"+strconv.Itoa(int(nodePort)))
			}
		}
	}
	return entrypoints, nil
}

func BatchDeployResources(cmd *cobra.Command, client *kubernetes.Clientset, flags *flags.Install, installFuncs []InstallFunc) error {
	for _, installer := range installFuncs {
		err := installer.Deploy(cmd, client, flags)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteStatefulsetResource(client *kubernetes.Clientset, resource, namespace, name string) error {
	err := client.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func DeleteAppsV1Resource(client *kubernetes.Clientset, resource, namespace, name string) error {
	err := client.AppsV1().RESTClient().Delete().Resource(resource).Namespace(namespace).Name(name).Do(context.Background()).Error()
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func DeleteCoreV1Resource(client *kubernetes.Clientset, resource, namespace, name string) error {
	err := client.CoreV1().RESTClient().Delete().Resource(resource).Namespace(namespace).Name(name).Do(context.Background()).Error()
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func DeleteRbacV1Resources(client *kubernetes.Clientset, resources, namespace, name string) error {
	err := client.RbacV1().RESTClient().Delete().Resource(resources).Namespace(namespace).Name(name).Do(context.Background()).Error()
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func DeleteCRDResource(client *apiextensions.Clientset, name string) error {
	err := client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

type deleteResourceFunc func(*kubernetes.Clientset, string, string, string) error

func DeleteResources(client *kubernetes.Clientset, resourceAndName [][]string, namespace string, deletefunc deleteResourceFunc) {
	for _, s := range resourceAndName {
		err := deletefunc(client, s[0], namespace, s[1])
		if err != nil {
			common.OutputErrorf("clear resource %s of %s in %s error: %s\n", s[1], s[0], namespace, err)
		}
	}
}

type PodStatus struct {
	Name            string
	ReadyContainer  int
	ExpectContainer int
	Status          string
	Restarts        int
}

func listPodByLabels(client *kubernetes.Clientset, labels map[string]string, namespace string) []PodStatus {
	i := 0
	labelSelector := ""
	for k, v := range labels {

		labelSelector += fmt.Sprintf("%s=%s", k, v)
		i++
		if i < len(labels) {
			labelSelector += ","
		}
	}

	podList, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		common.OutputErrorf("Ignore error %s", err)
		return nil
	}

	podStatus := []PodStatus{}
	for _, p := range podList.Items {
		ready, restart := readyAndRestartCount(p.Status.ContainerStatuses)
		podStatus = append(podStatus, PodStatus{
			Name:            p.Name,
			ReadyContainer:  ready,
			ExpectContainer: len(p.Spec.Containers),
			Status:          string(p.Status.Phase),
			Restarts:        restart,
		})
	}

	return podStatus
}

func readyAndRestartCount(statuses []v1.ContainerStatus) (readyCount, restartCount int) {
	for _, s := range statuses {
		if s.Ready {
			readyCount++
		}
		readyCount += int(s.RestartCount)
	}
	return
}

type ListPodFunc func(*kubernetes.Clientset, string) []PodStatus

func AdaptListPodFunc(labels map[string]string) ListPodFunc {
	return func(client *kubernetes.Clientset, namespace string) []PodStatus {
		return listPodByLabels(client, labels, namespace)
	}
}

func FormatPodStatus(client *kubernetes.Clientset, namespace string, fn ListPodFunc) (format string) {
	pods := fn(client, namespace)
	format += fmt.Sprintf("%-45s%-8s%-15s%-10s\n", "Name", "Ready", "Status", "Restarts")
	for _, p := range pods {
		format += fmt.Sprintf("%-45s%-8s%-15s%-10d\n", p.Name, fmt.Sprintf("%d/%d", p.ReadyContainer, p.ExpectContainer), p.Status, p.Restarts)
	}
	return
}
