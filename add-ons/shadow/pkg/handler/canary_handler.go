package handler

import (
	"log"

	"github.com/megaease/easemesh-api/v2alpha1"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/syncer"
	"github.com/megaease/easemeshctl/cmd/client/resource"
)

const (
	shadowServiceCanaryName = "shadow-service-canary"
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
	serviceCanary := createShadowServiceCanary(shadowServices)
	err := handler.applyShadowServiceCanary(serviceCanary)
	if err != nil {
		log.Printf("Create ServiceCanary for ShadowService failed. error: %s", err)
	}
}

// DeleteShadowService delete service from ServiceCanary's selector when ShadowService is deleted.
func (handler *ShadowServiceCanaryHandler) DeleteShadowService(obj interface{}) {
	shadowService := obj.(ShadowServiceBlock).service
	serviceCanary, err := handler.deleteShadowService(shadowService)
	err = handler.applyShadowServiceCanary(serviceCanary)
	if err != nil {
		log.Printf("Update ServiceCanary for ShadowService failed. ShadowService name: %s error: %s", shadowService.Name, err)
	}
}

func (handler *ShadowServiceCanaryHandler) applyShadowServiceCanary(serviceCanary *resource.ServiceCanary) error {
	canary, err := handler.Server.GetServiceCanary(serviceCanary.Name())
	if canary != nil {
		err = handler.Server.PatchServiceCanary(serviceCanary)
	} else {
		err = handler.Server.CreateServiceCanary(serviceCanary)
	}
	return err
}

func (handler *ShadowServiceCanaryHandler) deleteShadowService(shadowService object.ShadowService) (*resource.ServiceCanary, error) {
	canary, err := handler.Server.GetServiceCanary(shadowServiceCanaryName)
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

func createShadowServiceCanary(services []object.ShadowService) *resource.ServiceCanary {
	var matchServices []string
	for _, service := range services {
		matchServices = append(matchServices, service.ServiceName)
	}

	serviceCanary := &resource.ServiceCanary{
		MeshResource: resource.NewServiceCanaryResource(
			resource.DefaultAPIVersion, shadowServiceCanaryName,
		),
		Spec: &resource.ServiceCanarySpec{
			Priority: shadowServiceCanaryDefaultPriority,
			Selector: &v2alpha1.ServiceSelector{
				MatchServices: matchServices,
				MatchInstanceLabels: map[string]string{
					shadowServiceCanaryLabelKey: shadowServiceCanaryLabelValue,
				},
			},
			TrafficRules: &v2alpha1.TrafficRules{
				Headers: map[string]*v2alpha1.StringMatch{
					shadowServiceCanaryHeader: {
						Exact: shadowServiceCanaryHeaderValue,
					},
				},
			},
		},
	}
	return serviceCanary
}
