package handler

import (
	"reflect"
	"testing"

	"github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	"github.com/megaease/easemesh/mesh-shadow/pkg/syncer"
)

func Test_createShadowServiceCanary(t *testing.T) {
	shadowService1 := fake.NewShadowService()
	shadowService1.ServiceName = "service1"

	shadowService2 := fake.NewShadowService()
	shadowService2.ServiceName = "service2"

	shadowService3 := fake.NewShadowService()
	shadowService3.ServiceName = "service3"

	fakeCanary := fake.NewServiceCanary()
	ss := []object.ShadowService{
		shadowService1,
		shadowService2,
		shadowService3,
	}
	newCanary := createShadowServiceCanaries(ss)
	if !reflect.DeepEqual(newCanary, fakeCanary) {
		t.Errorf("createShadowServiceCanary() = %v, want %v", newCanary, fakeCanary)
	}
}

func TestShadowServiceCanaryHandler_deleteShadowService(t *testing.T) {
	handler := ShadowServiceCanaryHandler{
		Server: syncer.NewMockServer(),
	}
	shadowService1 := fake.NewShadowService()
	shadowService1.ServiceName = "service3"

	newServiceCanary, _ := handler.deleteShadowService(shadowService1)
	fakeServiceCanary := fake.NewDeletedServiceCanary()
	if !reflect.DeepEqual(newServiceCanary, fakeServiceCanary) {
		t.Errorf("handler.deleteShadowService() = %v, want %v", newServiceCanary, fakeServiceCanary)
	}
}
