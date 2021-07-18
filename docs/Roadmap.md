
# EaseMesh Roadmap

- [EaseMesh Roadmap](#easemesh-roadmap)
  - [Roadmap 2021](#roadmap-2021)
    - [Ease to Integrate](#ease-to-integrate)
    - [Java-Compatible Ecosystem](#java-compatible-ecosystem)
    - [High SLA](#high-sla)
    - [Backlogs](#backlogs)

As we said in [Purpose](../README.md#1-purposes) and [Principles](../README.md#2-principles), so we will focus on enhancing these features in 2021:

1. Easy to integrate
2. Java-compatible ecosystem
3. High SLA
4. Backlogs

## Roadmap 2021

### Ease to Integrate


| Description                                            | Priority | Related Issues |
| ------------------------------------------------------ | -------- | -------------- |
| Eliminate CRD MeshDeployment by k8s native deployment  | **High**     |                |
| Adapt to external service registry in spring ecosystem | **High**     |                |

### Java-Compatible Ecosystem

| Description                             | Priority | Related Issues |
| --------------------------------------- | -------- | -------------- |
| Support Spring Cloud 3.X or higher version| **High**     |                |
| Support Service discovery based on DNS mechanism | **High**     |                |
| Support Eclipse Vert.x                  | Middle   |                |

### High SLA


| Description                                 | Priority | Related Issues |
| ------------------------------------------- | -------- | -------------- |
| Tracing on and report metrics Elasticsearch | **High**     |                |
| Maturer production-ready canary deployment  | **High**     |                |
| Tracing and report metrics for MongoDB      | Middle   |                |
| Tracing and report metrics for ActiveMQ     | Low      |                |
| Tracing and report metrics for Amazon S3    | Low      |                |


### Backlogs

|Description | Priority | Related Issues |
|-|-|-|
|SMI support, support [ServiceMesh Interface](https://smi-spec.io)|Low||
|External registry, or co-exists with registries in spring ecosystem|Low||
|Multi-Cluster MeshControl and observe multiple clusters|Low||
|Fault Injection|Low||
|Delay Injection|Low||
|Access control|Low||
