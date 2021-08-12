package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"
	"github.com/megaease/easemesh/mesh-operator/pkg/deploymentmodifier"
	"github.com/megaease/easemesh/mesh-operator/pkg/util/labelstool"
	"github.com/pkg/errors"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
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
	// MutateHook handle requests from the injector MutatingWebhookConfiguration.
	MutateHook struct {
		*base.Runtime
		Admission *webhook.Admission
	}
)

// NewMutateHook creates a mutate hook.
func NewMutateHook(baseRuntime *base.Runtime) *MutateHook {
	h := &MutateHook{
		Runtime: baseRuntime,
	}
	h.Admission = &webhook.Admission{
		Handler: admission.HandlerFunc(h.mutateHandler),
	}
	h.Admission.InjectLogger(h.Log)

	return h
}

func (h *MutateHook) mutateHandler(cxt context.Context, req admission.Request) admission.Response {
	switch req.Operation {
	case admissionv1.Connect, admissionv1.Delete:
		return ignoreResp(req)
	}

	if req.Kind.Kind != "Deployment" {
		return ignoreResp(req)
	}

	deploy := &v1.Deployment{}
	err := json.Unmarshal(req.Object.Raw, &deploy)
	if err != nil {
		err := errors.Wrapf(err, "unmarshal json to Deployment: %s", req.String())
		h.Log.Error(err, "")
		return errorResp(err)
	}

	name := deploy.Annotations[annotationServiceNameKey]
	if name == "" {
		return ignoreResp(req)
	}

	h.Log.Info("mutate Deployment", "id", fmt.Sprintf("%s/%s", req.Namespace, req.Name))

	applicationPortValue := deploy.Annotations[annotationApplicationPortKey]
	var applicationPort uint16
	if applicationPortValue != "" {
		port, err := strconv.ParseUint(applicationPortValue, 10, 16)
		if err != nil {
			err := errors.Wrapf(err, "parse application port %s", applicationPortValue)
			h.Log.Error(err, "")
			return errorResp(err)
		}
		applicationPort = uint16(port)
	}

	labels, err := labelstool.Unmarshal(deploy.Annotations[annotationServiceLabels])
	if err != nil {
		h.Log.Error(err, "")
		return errorResp(err)
	}

	aliveProbeURL := deploy.Annotations[annotationAliveProbeURLKey]
	if aliveProbeURL == "" {
		aliveProbeURL = defaultAliveProbeURL
	}

	service := &deploymentmodifier.MeshService{
		Name:             name,
		Labels:           labels,
		AppContainerName: deploy.Annotations[annotationAppContainerNameKey],
		AliveProbeURL:    aliveProbeURL,
		ApplicationPort:  applicationPort,
	}
	modifier := deploymentmodifier.New(h.Runtime, service, deploy)

	err = modifier.Modify()
	if err != nil {
		err := errors.Wrapf(err, "modify deployment")
		h.Log.Error(err, "")
		return errorResp(err)
	}

	currentRaw, err := json.Marshal(deploy)
	if err != nil {
		err := errors.Wrapf(err, "marshal %#v to json failed", deploy)
		h.Log.Error(err, "")
		return errorResp(err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, currentRaw)
}

func ignoreResp(req admission.Request) admission.Response {
	return admission.Response{
		AdmissionResponse: admissionv1.AdmissionResponse{
			UID:     req.UID,
			Allowed: true,
		},
	}
}

func errorResp(err error) admission.Response {
	return admission.Errored(400, err)
}
