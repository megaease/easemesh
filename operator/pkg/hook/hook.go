package hook

import (
	"context"
	"fmt"

	"github.com/megaease/easemesh/mesh-operator/pkg/base"

	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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
	if !h.needInject(&req) {
		return ignoreResp(&req)
	}

	h.Log.Info("mutate", "id", fmt.Sprintf("%s %s/%s", req.Kind.Kind, req.Namespace, req.Name))
	currentRaw, err := h.injectSidecar(&req)
	if err != nil {
		h.Log.Error(err, "")
		return errorResp(err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, currentRaw)
}

func ignoreResp(req *admission.Request) admission.Response {
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
