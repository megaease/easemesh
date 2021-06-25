package meshclient

const (
	apiURL = "/apis/v1"

	// MeshTenantsURL is the mesh tenant prefix.
	MeshTenantsURL = apiURL + "/mesh/tenants"

	// MeshTenantURL is the mesh tenant path.
	MeshTenantURL = apiURL + "/mesh/tenants/%s"

	// MeshServicesURL is mesh service prefix.
	MeshServicesURL = apiURL + "/mesh/services"

	// MeshServiceURL is the mesh service path.
	MeshServiceURL = apiURL + "/mesh/services/%s"

	// MeshServiceCanaryURL is the mesh service canary path.
	MeshServiceCanaryURL = apiURL + "/mesh/services/%s/canary"

	// MeshServiceResilienceURL is the mesh service resilience path.
	MeshServiceResilienceURL = apiURL + "/mesh/services/%s/resilience"

	// MeshServiceLoadBalanceURL is the mesh service load balance path.
	MeshServiceLoadBalanceURL = apiURL + "/mesh/services/%s/loadbalance"

	// MeshServiceOutputServerURL is the mesh service output server path.
	MeshServiceOutputServerURL = apiURL + "/mesh/services/%s/outputserver"

	// MeshServiceTracingsURL is the mesh service tracings path.
	MeshServiceTracingsURL = apiURL + "/mesh/services/%s/tracings"

	// MeshServiceMetricsURL is the mesh service metrics path.
	MeshServiceMetricsURL = apiURL + "/mesh/services/%s/metrics"

	// MeshServiceInstancesURL is the mesh service prefix.
	MeshServiceInstancesURL = apiURL + "/mesh/serviceinstances"

	// MeshServiceInstanceURL is the mesh service path.
	MeshServiceInstanceURL = apiURL + "/mesh/serviceinstances/%s/%s"

	// MeshIngressURL is the mesh ingress prefix.
	MeshIngressesURL = apiURL + "/mesh/ingresses"

	// MeshIngressURL is the mesh ingress path.
	MeshIngressURL = apiURL + "/mesh/ingresses/%s"
)
