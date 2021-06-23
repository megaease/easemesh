package meshclient

type meshClient struct {
	server   string
	v1Alpha1 V1Alpha1Interface
}

func (m *meshClient) V1Alpha1() V1Alpha1Interface {
	return m.v1Alpha1
}

type v1alpha1Interface struct {
	loadbalanceGetter
	canaryGetter
	resilienceGetter
	serviceGetter
	tenantGetter
	observabilityGetter
	ingressGetter
}

var _ V1Alpha1Interface = &v1alpha1Interface{}

func New(server string) MeshClient {
	client := &meshClient{server: server}
	alpha1 := v1alpha1Interface{
		loadbalanceGetter:   loadbalanceGetter{client: client},
		canaryGetter:        canaryGetter{client: client},
		resilienceGetter:    resilienceGetter{client: client},
		tenantGetter:        tenantGetter{client: client},
		observabilityGetter: observabilityGetter{client: client},
		serviceGetter:       serviceGetter{client: client},
		ingressGetter:       ingressGetter{client: client},
	}
	client.v1Alpha1 = &alpha1
	return client
}
