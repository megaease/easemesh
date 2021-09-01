package hook

import (
	"encoding/json"
	"strconv"

	"github.com/megaease/easemesh/mesh-operator/pkg/sidecarinjector"
	"github.com/megaease/easemesh/mesh-operator/pkg/util/labelstool"
	"github.com/pkg/errors"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	annotationPrefix              = "mesh.megaease.com/"
	annotationServiceNameKey      = annotationPrefix + "service-name"
	annotationServiceLabels       = annotationPrefix + "service-labels"
	annotationAppContainerNameKey = annotationPrefix + "app-container-name"
	annotationApplicationPortKey  = annotationPrefix + "application-port"
	annotationAliveProbeURLKey    = annotationPrefix + "alive-probe-url"

	defaultAliveProbeURL = "http://localhost:9900/health"
)

type (
	// BaseObject contains base fields of k8s object.
	BaseObject struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata"`
	}
)

func (h *MutateHook) needInject(req *admission.Request) bool {
	switch req.Operation {
	case admissionv1.Connect, admissionv1.Delete:
		return false
	}

	switch req.Kind.Kind {
	case "Pod", "ReplicaSet", "Deployment", "StatefulSet", "DaemonSet":
	default:
		return false
	}

	baseObject := &BaseObject{}
	err := json.Unmarshal(req.Object.Raw, baseObject)
	if err != nil {
		h.Log.Error(err, "unmarshal json to base object", "raw", req.String())
		return false
	}

	if baseObject.Annotations[annotationServiceNameKey] == "" {
		return false
	}

	return true
}

func (h *MutateHook) extractMeshService(baseObject *BaseObject) (*sidecarinjector.MeshService, error) {
	name := baseObject.Annotations[annotationServiceNameKey]
	if name == "" {
		return nil, errors.New("no service name")
	}

	applicationPortValue := baseObject.Annotations[annotationApplicationPortKey]
	var applicationPort uint16
	if applicationPortValue != "" {
		port, err := strconv.ParseUint(applicationPortValue, 10, 16)
		if err != nil {
			return nil, errors.Wrapf(err, "parse application port %s", applicationPortValue)
		}
		applicationPort = uint16(port)
	}

	labels, err := labelstool.Unmarshal(baseObject.Annotations[annotationServiceLabels])
	if err != nil {
		return nil, err
	}

	aliveProbeURL := baseObject.Annotations[annotationAliveProbeURLKey]
	if aliveProbeURL == "" {
		aliveProbeURL = defaultAliveProbeURL
	}

	return &sidecarinjector.MeshService{
		Name:             name,
		Labels:           labels,
		AppContainerName: baseObject.Annotations[annotationAppContainerNameKey],
		AliveProbeURL:    aliveProbeURL,
		ApplicationPort:  applicationPort,
	}, nil
}

func (h *MutateHook) injectSidecar(req *admission.Request) ([]byte, error) {
	baseObject := &BaseObject{}
	err := json.Unmarshal(req.Object.Raw, baseObject)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal json %s to base object", req.String())
	}

	meshService, err := h.extractMeshService(baseObject)
	if err != nil {
		return nil, err
	}

	object := h.newObject(req.Kind.Kind)
	err = json.Unmarshal(req.Object.Raw, object)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal json %s", req.String())
	}

	podSpec := h.getPodSpec(object)

	injector := sidecarinjector.New(h.Runtime, meshService, podSpec)
	err = injector.Inject()
	if err != nil {
		return nil, errors.Wrapf(err, "inject sidecar")
	}

	currentRaw, err := json.Marshal(object)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal %+v to json", object)
	}

	return currentRaw, nil
}

func (h *MutateHook) newObject(kind string) interface{} {
	switch kind {
	case "Pod":
		return &corev1.Pod{}
	case "ReplicaSet":
		return &v1.ReplicaSet{}
	case "Deployment":
		return &v1.Deployment{}
	case "StatefulSet":
		return &v1.StatefulSet{}
	case "DaemonSet":
		return &v1.DaemonSet{}
	}

	return nil
}

func (h *MutateHook) getPodSpec(object interface{}) *corev1.PodSpec {
	switch obj := object.(type) {
	case *corev1.Pod:
		return &obj.Spec
	case *v1.ReplicaSet:
		return &obj.Spec.Template.Spec
	case *v1.Deployment:
		return &obj.Spec.Template.Spec
	case *v1.StatefulSet:
		return &obj.Spec.Template.Spec
	case *v1.DaemonSet:
		return &obj.Spec.Template.Spec
	}

	return nil
}
