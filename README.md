[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/xk6-kubernetes.svg)](https://pkg.go.dev/github.com/grafana/xk6-kubernetes)
[![Version Badge](https://img.shields.io/github/v/release/grafana/xk6-kubernetes?style=flat-square)](https://github.com/grafana/xk6-kubernetes/releases)
![Build Status](https://img.shields.io/github/actions/workflow/status/grafana/xk6-kubernetes/ci.yml?style=flat-square)

# xk6-kubernetes
A k6 extension for interacting with Kubernetes clusters while testing.

## Build

To build a custom `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

1. Download [xk6](https://github.com/grafana/xk6):
  
    ```bash
    go install go.k6.io/xk6/cmd/xk6@latest
    ```

2. [Build the k6 binary](https://github.com/grafana/xk6#command-usage):
  
    ```bash
    xk6 build --with github.com/grafana/xk6-kubernetes
    ```

    The `xk6 build` command creates a k6 binary that includes the xk6-kubernetes extension in your local folder. This k6 binary can now run a k6 test using [xk6-kubernetes APIs](#apis).


### Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```shell
git clone git@github.com:grafana/xk6-kubernetes.git
cd xk6-kubernetes
make
```

Using the `k6` binary with `xk6-kubernetes`, run the k6 test as usual:

```bash
./k6 run k8s-test-script.js

```
# Usage

The API assumes a `kubeconfig` configuration is available at any of the following default locations:
* at the location pointed by the `KUBECONFIG` environment variable
* at `$HOME/.kube`


# API

## Generic API

This API offers methods for creating, retrieving, listing and deleting resources of any of the supported kinds.

|  Method     | Parameters|   Description |
| ------------ | ---| ------ |
| apply         | manifest string| creates a Kubernetes resource given a YAML manifest or updates it if already exists |
| create         | spec object | creates a Kubernetes resource given its specification |
| delete         | kind  | removes the named resource |
|                | name  |
|                | namespace|
| get         | kind| returns the named resource |
|                | name  |
|                | namespace |
| list         | kind| returns a collection of resources of a given kind
|                | namespace |
| update         | spec object | updates an existing resource


The kinds of resources currently supported are:
* ConfigMap
* Deployment
* Ingress
* Job
* Namespace
* Node
* PersistentVolume
* PersistentVolumeClaim
* Pod
* Secret
* Service
* StatefulSet

### Examples

#### Creating a pod using a specification 
```javascript
import { Kubernetes } from 'k6/x/kubernetes';

const podSpec = {
    apiVersion: "v1",
    kind:       "Pod",
    metadata: {
        name:      "busybox",
        namespace: "testns"
    },
    spec: {
        containers: [
            {
                name:    "busybox",
                image:   "busybox",
                command: ["sh", "-c", "sleep 30"]
            }
        ]
    }
}

export default function () {
  const kubernetes = new Kubernetes();

  kubernetes.create(pod)

  const pods = kubernetes.list("Pod", "testns");

  console.log(`${pods.length} Pods found:`);
  pods.map(function(pod) {
    console.log(`  ${pod.metadata.name}`)
  });
}
```

#### Creating a job using a YAML manifest
```javascript
import { Kubernetes } from 'k6/x/kubernetes';

const manifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: busybox
  namespace: testns
spec:
  template:
    spec:
      containers:
      - name: busybox
        image: busybox
        command: ["sleep", "300"]
    restartPolicy: Never
`

export default function () {
  const kubernetes = new Kubernetes();

  kubernetes.apply(manifest)

  const jobs = kubernetes.list("Job", "testns");

  console.log(`${jobs.length} Jobs found:`);
  pods.map(function(job) {
    console.log(`  ${job.metadata.name}`)
  });
}
```

## Helpers

The `xk6-kubernetes` extension offers helpers to facilitate common tasks when setting up a tests. All helper functions work in a namespace to facilitate the development of tests segregated by namespace. The helpers are accessed using the following method:

|  Method      | Parameters|   Description |
| -------------| ---| ------ |
| helpers      | namespace | returns helpers that operate in the given namespace. If none is specified, "default" is used |

The methods above return an object that implements the following helper functions:

|  Method     | Parameters|   Description |
| ------------ | --------| ------ |
| getExternalIP        | service        | returns the external IP of a service if any is assigned before timeout expires|
|                      | timeout in seconds | |
| waitPodRunning | pod name | waits until the pod is in 'Running' state or the timeout expires. Returns a boolean indicating of the pod was ready or not. Throws an error if the pod is Failed. |
|                | timeout in seconds | |
| waitServiceReady         | service name | waits until the given service has at least one endpoint ready or the timeout expires |
|                | timeout in seconds | |



### Examples

### Creating a pod in a random namespace and wait until it is running

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

let podSpec = {
    apiVersion: "v1",
    kind:       "Pod",
    metadata: {
        name:      "busybox",
        namespace:  "default"
    },
    spec: {
        containers: [
            {
                name:    "busybox",
                image:   "busybox",
                command: ["sh", "-c", "sleep 30"]
            }
        ]
    }
}

export default function () {
  const kubernetes = new Kubernetes();

  // create pod
  kubernetes.create(pod)

  // get helpers for test namespace
  const helpers = kubernetes.helpers()

  // wait for pod to be running
  const timeout = 10
  if (!helpers.waitPodRunning(pod.metadata.name, timeout)) {
      console.log(`"pod ${pod.metadata.name} not ready after ${timeout} seconds`)
  }
}
```

## Resource kind helpers

This API offers a helper for each kind of Kubernetes resources supported (Pods, Deployments, Secrets, et cetera). For each one, an interface for creating, getting, listing and deleting objects is offered. 

>⚠️ This interface is deprecated and will be removed soon
> -
Migrate to the usage of the generic resources API.
</br>


### (Deprecated) Create a client:  `new Kubernetes(config)`

Creates a Kubernetes client to interact with the Kubernetes cluster.

|  Config options |   Type | Description                              |  Default |
| ------------ | ------ | ---------------------------------------- | ------ |
| config_path         | String    | The path to the kubeconfig file         |  ~/.kube/config |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({
    // config_path: "/path/to/kube/config"
  })
}
```

### (Deprecated) `Client.config_maps`

|  Method     |   Description |
| ------------ | ------ |
| apply         | creates the Kubernetes resource given a YAML configuration |
| create         | creates the Kubernetes resource given an object configuration |
| delete         | removes the named ConfigMap     |
| get         | returns the named ConfigMaps     |
| list         | returns a collection of ConfigMaps     |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

  const nameSpace = "default";
  const name = "config-map-name";
  kubernetesClient.config_maps.apply(getConfigMapYaml(name), nameSpace);
}
```

### (Deprecated) `Client.deployments`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-deployment-service-ingress.js](./examples/apply-deployment-service-ingress.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named Deployment     |   |
| get         | returns the named Deployment     | [get-configmap.js](./examples/get-deployment.js)  |
| list         | returns a collection of Deployments     | [list-configmaps.js](./examples/list-deployments.js)   |



```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

  const nameSpace = "default";
  const name = "deployment-name";
  const app = 'app-label';

  kubernetesClient.deployments.apply(getDeploymentYaml(name, app), nameSpace);
}
```


### (Deprecated) `Client.ingresses`

|  Method     |   Description |
| ------------ | ------ |
| apply         | creates the Kubernetes resource given a YAML configuration    |
| create         | creates the Kubernetes resource given an object configuration    |
| delete         | removes the named Ingress     |
| get         | returns the named Ingress     |
| list         | returns a collection of Ingresses     |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

  const nameSpace = "default";
  const name = "deployment-name";
  const url = 'ingress-url.com';

  kubernetesClient.ingresses.apply(getIngressYaml(name, url), nameSpace);
}
```

### (Deprecated) `Client.jobs`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-job.js](./examples/apply-job.js)  |
| create         | creates the Kubernetes resource given an object configuration    | [create-job.js](./examples/create-job.js), [create-job-wait.js](./examples/create-job-wait.js), [create-job-by-nodename.js](./examples/create-job-by-nodename.js), [create-job-autodelete.js](./examples/create-job-autodelete.js)   |
| delete         | removes the named Job     |   |
| get         | returns the named Jobs     | [get-job.js](./examples/get-job.js)  |
| list         | returns a collection of Jobs     | [list-jobs.js](./examples/list-jobs.js)   |
| wait         | wait for all Jobs to complete    | [wait-job.js](./examples/wait-job.js)   |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const namespace = "default"
  const jobName = "new-job"
  const image = "perl"
  const command = ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]

  kubernetesClient.jobs.create({
    namespace: namespace,
    name: jobName,
    image: image,
    command: command
  })

  const completed = kubernetesClient.jobs.wait({
    namespace: namespace,
    name: jobName,
    timeout: "30s"
  })
  const jobStatus = completed? "completed": "not completed"
}
```

### (Deprecated) `Client.namespaces`

|  Method     |   Description |
| ------------ | ------ |
| apply         | creates the Kubernetes resource given a YAML configuration    |
| create         | creates the Kubernetes resource given an object configuration    |
| delete         | removes the named Namespaces     |
| get         | returns the named Namespace    |
| list         | returns a collection of Namespaces     |


```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const name = "namespace-name";

  kubernetesClient.namespaces.apply(getNamespaceYaml(name));
}
```

### (Deprecated) `Client.nodes`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| list         | returns a collection of Nodes comprising the cluster    | [list-nodes.js](./examples/list-nodes.js) |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const nodes = kubernetesClient.nodes.list()
}
```

### (Deprecated) `Client.persistent_volumes`

|  Method     |   Description |
| ------------ | ------ |
| apply         | creates the Kubernetes resource given a YAML configuration    |
| create         | creates the Kubernetes resource given an object configuration    |
| delete         | removes the named persistent volume     |
| get         | returns the named persistent  volume instance     |
| list         | returns a collection of persistent volumens     |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

  const name = "example-pv";
  kubernetesClient.persistent_volumes.apply(getPVYaml(name, "1Gi", "local-storage"));
}
```

### (Deprecated) `Client.persistent_volumes_claims`

|  Method     |   Description |
| ------------ | ------ |
| apply         | creates the Kubernetes resource given a YAML configuration    |
| create         | creates the Kubernetes resource given an object configuration    |
| delete         | removes the named persistent volume claim     |
| get         | returns the named persistent volume claim     |
| list         | returns a collection of persistent volumen claims     |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

    const name = "example-pvc";
    const nameSpace = "default";

    kubernetes.persistent_volume_claims.apply(getPVCYaml(name, "1Gi", "nfs-csi"), nameSpace);
}
```

### (Deprecated) `Client.pods`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| create         | runs a pod    | [create-pod.js](./examples/create-pod.js), [create-pod-wait.js](./examples/create-pod-wait.js)  |
| delete         | removes the named Pod     |   |
| get         | returns the named Pod     | [get-pod.js](./examples/get-pod.js)  |
| list         | returns a collection of Pods     | [list-pods.js](./examples/list-pods.js)   |
| wait         | wait for the Pod to be in a given status    | [wait-pod.js](./examples/wait-pod.js)   |
| exec         | executes a non-interactive command    | [exec-command.js](./examples/exec-command.js)   |
| addEphemeralContainer         | adds an ephemeral container to a running pod    | [add-ephemeral.js](./examples/add-ephemeral.js)   |



```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const namespace = "default"
  const podName = "new-pod"
  const image = "busybox"
  const command = ["sh",  "-c", "sleep 5"]

  kubernetesClient.pods.create({
    namespace: namespace,
    name: podName,
    image: image,
    command: command
  });
 
  const options = {
    namespace: namespace,
    name: podName,
    status: "Succeeded",
    timeout: "10s"
  }
  if (kubernetesClient.pods.wait(options)) {
    console.log(podName + " pod completed successfully")
  } else {
    throw podName + " is not completed"
  }
}
```

### (Deprecated) `Client.secrets`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-secret.js](./examples/apply-secret.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named secret     |   |
| get         | returns the named secret    | [get-secret.js](./examples/get-secret.js)  |
| list         | returns a collection of secrets     | [list-secrets.js](./examples/list-secrets.js)   |


```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const secrets = kubernetesClient.secrets.list()
}
```

### (Deprecated) `Client.services`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-deployment-service-ingress.js](./examples/apply-deployment-service-ingress.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named service     |   |
| get         | returns the named service    | [get-service.js](./examples/get-service.js)  |
| list         | returns a collection of services     | [list-services.js](./examples/list-services.js)   |


```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const svcs = kubernetesClient.services.list()
}
```


## If things go wrong


### Are you using the custom binary?

  An easy mistake--which happens often--is to forget that `xk6` is generating a new executable. You may be accustomed to simply running `k6` from the command-line which probably isn't your new build. Make sure to use `./k6` after building your extended version otherwise you can expect to see an error similar to:


  ```bash
  ERRO[0000] The moduleSpecifier "k8s-test-script.js" couldn't be found on local disk. Make sure that you've specified the right path to the file. If you're running k6 using the Docker image make sure you have mounted the local directory (-v /local/path/:/inside/docker/path) containing your script and modules so that they're accessible by k6 from inside of the container, see https://k6.io/docs/using-k6/modules#using-local-modules-with-docker. Additionally it was tried to be loaded as remote module by prepending "https://" to it, which also didn't work. Remote resolution error: "Get "https://k8s-test-script.js": dial tcp: lookup k8s-test-script.js: no such host" 
  ```
