package handler

import (
	"log"

	"github.com/megaease/easemesh-api/v1alpha1"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/syncer"
	"github.com/megaease/easemeshctl/cmd/client/resource"
)

type ShadowServiceCanaryHandler struct {
	Server *syncer.Server
}

// CreateServiceCanary create ServiceCanary when ShadowService is created.
func (handler *ShadowServiceCanaryHandler) CreateServiceCanary(obj interface{}) {
	block := obj.(ShadowServiceBlock)
	shadowService := block.service
	err := handler.applyShadowServiceCanary(&shadowService)
	if err != nil {
		log.Printf("Apply ServiceCanary for ShadowService failed. ShadowService name: %s error: %s", shadowService.Name, err)
	}
}

// DeleteServiceCanary delete ServiceCanary when ShadowService is deleted.
func (handler *ShadowServiceCanaryHandler) DeleteServiceCanary(obj interface{}) {

	block := obj.(ShadowServiceBlock)
	shadowService := block.service
	err := handler.Server.DeleteServiceCanary(shadowService.Name)
	if err != nil {
		log.Printf("Delete ServiceCanary for ShadowService failed. ShadowService name: %s error: %s", shadowService.Name, err)
	}
}

func (handler *ShadowServiceCanaryHandler) applyShadowServiceCanary(shadowService *object.ShadowService) error {
	serviceCanary := createShadowServiceCanary(shadowService)
	canary, err := handler.Server.GetServiceCanary(shadowService.Name)
	if err != nil {
		return err
	}
	if canary != nil {
		err = handler.Server.PatchServiceCanary(serviceCanary)
	} else {
		err = handler.Server.CreateServiceCanary(serviceCanary)
	}
	return err
}

func createShadowServiceCanary(obj *object.ShadowService) *resource.ServiceCanary {
	return &resource.ServiceCanary{
		MeshResource: resource.NewServiceCanaryResource(
			resource.DefaultAPIVersion, obj.Name,
		),
		Spec: &resource.ServiceCanarySpec{
			Priority: shadowServiceCanaryDefaultPriority,
			Selector: &v1alpha1.ServiceSelector{
				MatchServices: []string{obj.ServiceName},
				MatchInstanceLabels: map[string]string{
					shadowServiceCanaryLabelKey: shadowServiceCanaryLabelValue,
				},
			},
			TrafficRules: &v1alpha1.TrafficRules{
				Headers: map[string]*v1alpha1.StringMatch{
					shadowServiceCanaryHeader: &v1alpha1.StringMatch{
						Exact: shadowServiceCanaryHeaderValue,
					},
				},
			},
		},
	}
}
