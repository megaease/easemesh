## EaseMesh Deployment Operator

### Introduction


`meshdeployment-operator` is an operator who looks after services or applications in the mesh cluster of the MegaEase. We make use of [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) to deploy the mesh service or application, so `meshdeployement-operator` is a K8s's controller for our `meshdeployment` of custom resource


An `meshdeployment` resource is a K8s deployment resource which we enhance it with the name of Mesh Service and bind the metadata when the application register itself into the service registry of the EaseMesh

### How to build it

1. Before you build it, you need to generate the Custom Resource Definition via the following command:

```shell
make manifest
```

2. After generate manifest, some boilerplate codes need to be generated via the following command:

```shell
make generate
```

3. You can build the controller as a standalone binary, via the following command

```shell
make 
```

4. Except for standalone executable, in most production environments, the controller was run as a deployment resource of K8s, so you could build a docker image and push it to the local/public registry according to your image tag.

```bash
make docker-build docker-push IMG={image_name}:{image_version}
```

Replace `{image_name}` and `{image_version}` to your proper value


### How to deploy it

After you built a standalone binary or docker image, you can boot it up in two ways, one is `in-cluster` another is `out-cluster`

- `in-cluster` need build operation as a docker image, and deployed in K8s cluster
- `out-cluster` needs to build as an executable, boot it up as a standalone process.


You can use the following command to boot operator in the `in-cluster` way

```bash
make deploy IMG={image_name}:{image_version}
```

Replace `{image_name}` and `{image_version}` to your proper value.

> Attention: the k8s cluster has right 

If you want to run the operator in `out-cluster` way, you need to ensure that you can access a K8s cluster from your environment. Ensure you have the right to manage your K8s cluster.

An extra step needed to run the operator in the `out-cluster` way, which is to apply the custom resource definition before start the executable.

```bash
kustomize build config/default | kubectl apply -f - 
```

> `kustomize` can't be found, you can refer https://kubectl.docs.kubernetes.io/installation/kustomize to install it.

After the `CRD` was applied, run the executable via the following command:

```bash
bin/manage
```

### How to debug the operator

If you use the VsCode develop operator, you can leverage `dlv` and `out-cluster` deployment to debug our program. Edit a launch.json in .vscode directory

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "test",
            "type": "go",
            "request": "launch",
            "mode": "exec",
            "remotePath": "${workspaceFolder}/main.go",
            "port": 2345,
            "host": "127.0.0.1",
            "program": "${workspaceFolder}/bin/manager",
            "env": {},
            "args": []
        }
    ]
}
```

### Run the example

In the examples directory, we have an `MeshDeployment` spec file. In the file, we defined a `test-server-v1` meshdeployment resource, whose deployment contains a test-server container listening on `18080` port.  We could use the following command to deploy the test-server-v1 `MeshDeployment`:

```bash
kubectl apply -f examples
```

When we deployed the resource, an injected container named `easemesh-sidecar` will be injected into the pod of the MeshDeployment. A k8s deployment resource owned by the `test-server-v1` will be generated and deployed in the `test` namespace, it has two replicas of pods, assuming a pod named `test-server-v1-7d7bccf78f-2pps5`. The following command can help us to check whether the extra container has been injected

```bash
kubectl get pods -n test test-server-v1-7d7bccf78f-2pps5 -o jsonpath="{.spec.containers[*].name}"
```

The output should be:

```bash
test-server easemesh-sidecar
```
