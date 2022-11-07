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
	"fmt"
	"log"

	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/syncer"
	"github.com/megaease/easemeshctl/cmd/client/resource"
)

// ShadowServiceCanaryHandler  added or deleted according to the creation and deletion of ShadowService.
type ShadowServiceCanaryHandler struct {
	Server syncer.MeshControlPlane
}

// GenerateServiceCanary create ServiceCanary for all ShadowService.
func (handler *ShadowServiceCanaryHandler) GenerateServiceCanary(objs interface{}) {
	shadowServices := objs.([]object.ShadowService)
	if len(shadowServices) == 0 {
		return
	}
	serviceCanaries := createShadowServiceCanaries(shadowServices)
	err := handler.applyShadowServiceCanaries(serviceCanaries)
	if err != nil {
		log.Printf("apply service canaries for shadow service failed: %s", err)
		return
	}

	log.Printf("apply service canaries for shadow service succeed")
}

// DeleteShadowService delete service from ServiceCanary's selector when ShadowService is deleted.
func (handler *ShadowServiceCanaryHandler) DeleteShadowService(obj interface{}) {
	shadowService := obj.(ShadowServiceBlock).shadowService

	serviceCanary, err := handler.deleteShadowService(shadowService)
	if err != nil {
		log.Printf("delete shadow service failed: %v", err)
		return
	}

	if len(serviceCanary.Spec.Selector.MatchServices) == 0 {
		err = handler.Server.DeleteServiceCanary(serviceCanary.Name())
		if err != nil {
			log.Printf("delete service canary %s failed: %v", serviceCanary.Name(), err)
			return
		}

		log.Printf("delete service canary %s succeed", serviceCanary.Name())
		return
	}

	err = handler.Server.PatchServiceCanary(serviceCanary)
	if err != nil {
		log.Printf("update service canary %s failed: %v", serviceCanary.Name(), err)
		return
	}

	log.Printf("update service canary %s succeed", serviceCanary.Name())
}

func (handler *ShadowServiceCanaryHandler) applyShadowServiceCanaries(serviceCanaries map[string]*resource.ServiceCanary) error {
	for _, canary := range serviceCanaries {
		oldCanary, err := handler.Server.GetServiceCanary(canary.Name())

		if oldCanary != nil {
			err = handler.Server.PatchServiceCanary(canary)
			if err != nil {
				return fmt.Errorf("update service canary %s failed: %v", canary.Name(), err)
			}
		} else {
			err = handler.Server.CreateServiceCanary(canary)
			if err != nil {
				return fmt.Errorf("create service canary %s failed: %v", canary.Name(), err)
			}
		}
	}

	return nil
}

func (handler *ShadowServiceCanaryHandler) deleteShadowService(shadowService object.ShadowService) (*resource.ServiceCanary, error) {
	canaryname := shadowService.CanaryName()

	canary, err := handler.Server.GetServiceCanary(canaryname)
	if canary == nil {
		return nil, err
	}

	matchServices := canary.Spec.Selector.MatchServices
	var newMatchServices []string
	for _, serviceName := range matchServices {
		if serviceName == shadowService.ServiceName {
			continue
		}
		newMatchServices = append(newMatchServices, serviceName)
	}
	canary.Spec.Selector.MatchServices = newMatchServices

	return canary, err
}

func createShadowServiceCanaries(shadowServices []object.ShadowService) map[string]*resource.ServiceCanary {
	// The key is shadow header value.
	canaries := map[string]*resource.ServiceCanary{}
	for _, ss := range shadowServices {
		canaryName := ss.CanaryName()

		canary, exists := canaries[canaryName]
		if !exists {

			canary = &resource.ServiceCanary{
				MeshResource: resource.NewServiceCanaryResource(resource.DefaultAPIVersion, canaryName),

				Spec: &resource.ServiceCanarySpec{
					Priority: shadowServiceCanaryDefaultPriority,
					Selector: &v2alpha1.ServiceSelector{
						MatchServices: []string{ss.ServiceName},
						MatchInstanceLabels: map[string]string{
							shadowServiceCanaryLabelKey: canaryName,
						},
					},
					TrafficRules: &v2alpha1.TrafficRules{
						Headers: map[string]*v2alpha1.StringMatch{
							shadowServiceCanaryHeader: {
								Exact: ss.TrafficHeaderValue,
							},
						},
					},
				},
			}

			canaries[canaryName] = canary
		}

		appended := false
		for _, serviceName := range canary.Spec.Selector.MatchServices {
			if serviceName == ss.ServiceName {
				appended = true
				break
			}
		}

		if !appended {
			canary.Spec.Selector.MatchServices = append(canary.Spec.Selector.MatchServices, ss.ServiceName)
		}
	}

	return canaries
}
