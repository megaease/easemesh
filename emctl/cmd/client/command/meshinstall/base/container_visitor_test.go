package installbase

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

type fakeContainerVisitor struct{}

func (v *fakeContainerVisitor) VisitorCommandAndArgs(c *v1.Container) (command []string, args []string) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorEnvs(c *v1.Container) ([]v1.EnvVar, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorEnvFrom(c *v1.Container) ([]v1.EnvFromSource, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorResourceRequirements(c *v1.Container) (*v1.ResourceRequirements, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorVolumeMounts(c *v1.Container) ([]v1.VolumeMount, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorVolumeDevices(c *v1.Container) ([]v1.VolumeDevice, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorLivenessProbe(c *v1.Container) (*v1.Probe, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorReadinessProbe(c *v1.Container) (*v1.Probe, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorLifeCycle(c *v1.Container) (*v1.Lifecycle, error) {
	return nil, nil
}
func (v *fakeContainerVisitor) VisitorSecurityContext(c *v1.Container) (*v1.SecurityContext, error) {
	return nil, nil
}

func TestAcceptContainerVisitor(t *testing.T) {

	_, err := AcceptContainerVisitor("a", "b", v1.PullAlways, &fakeContainerVisitor{})
	if err != nil {
		t.Fatalf("visits fakeContainer error: %s", err)
	}

}
