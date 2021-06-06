package installbase

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
)

type ContainerVisitor interface {
	VisitorCommandAndArgs(c *v1.Container) (command []string, args []string)
	VisitorContainerPorts(c *v1.Container) ([]v1.ContainerPort, error)
	VisitorEnvs(c *v1.Container) ([]v1.EnvVar, error)
	VisitorEnvFrom(c *v1.Container) ([]v1.EnvFromSource, error)
	VisitorResourceRequirements(c *v1.Container) (*v1.ResourceRequirements, error)
	VisitorVolumeMounts(c *v1.Container) ([]v1.VolumeMount, error)
	VisitorVolumeDevices(c *v1.Container) ([]v1.VolumeDevice, error)
	VisitorLivenessProbe(c *v1.Container) (*v1.Probe, error)
	VisitorReadinessProbe(c *v1.Container) (*v1.Probe, error)
	VisitorLifeCycle(c *v1.Container) (*v1.Lifecycle, error)
	VisitorSecurityContext(c *v1.Container) (*v1.SecurityContext, error)
}

func AcceptContainerVisistor(name, image string, imagePullPolicy v1.PullPolicy, visitor ContainerVisitor) (*v1.Container, error) {
	container := &v1.Container{Name: name, Image: image, ImagePullPolicy: imagePullPolicy}
	command, args := visitor.VisitorCommandAndArgs(container)
	if command != nil {
		container.Command = command
	}
	if args != nil {
		container.Args = args
	}

	ports, err := visitor.VisitorContainerPorts(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(ports, func() { container.Ports = ports })

	envs, err := visitor.VisitorEnvs(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(envs, func() { container.Env = envs })

	envFromSource, err := visitor.VisitorEnvFrom(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(envFromSource, func() { container.EnvFrom = envFromSource })

	resources, err := visitor.VisitorResourceRequirements(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(resources, func() { container.Resources = *resources })

	volumeMounts, err := visitor.VisitorVolumeMounts(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(volumeMounts, func() { container.VolumeMounts = volumeMounts })

	volumeDevices, err := visitor.VisitorVolumeDevices(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(volumeDevices, func() { container.VolumeDevices = volumeDevices })

	livenessProbe, err := visitor.VisitorLivenessProbe(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(livenessProbe, func() { container.LivenessProbe = livenessProbe })

	readinessProbe, err := visitor.VisitorReadinessProbe(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(readinessProbe, func() { container.ReadinessProbe = readinessProbe })

	lifecycle, err := visitor.VisitorLifeCycle(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(lifecycle, func() { container.Lifecycle = lifecycle })

	securityContext, err := visitor.VisitorSecurityContext(container)
	if err != nil {
		return nil, err
	}
	setIfNotNull(securityContext, func() { container.SecurityContext = securityContext })
	return container, nil
}

func setIfNotNull(ele interface{}, fn func()) {
	if ele == nil {
		return
	}
	switch reflect.TypeOf(ele).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		if reflect.ValueOf(ele).IsNil() {
			return
		}
	}
	fn()
}
