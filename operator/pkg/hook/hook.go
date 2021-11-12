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
