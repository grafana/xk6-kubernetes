[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/xk6-kubernetes.svg)](https://pkg.go.dev/github.com/grafana/xk6-kubernetes)
[![Version Badge](https://img.shields.io/github/v/release/grafana/xk6-kubernetes?style=flat-square)](https://github.com/grafana/xk6-kubernetes/releases)
![Build Status](https://img.shields.io/github/workflow/status/grafana/xk6-kubernetes/CI?style=flat-square)

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


## Example

```javascript
import { Kubernetes } from 'k6/x/kubernetes';
  
export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  });

  const pods = kubernetes.pods.list();

  console.log(`${pods.length} Pods found:`);
  pods.map(function(pod) {
    console.log(`  ${pod.name}`)
  });
}
```

Using the `k6` binary with `xk6-kubernetes`, run the k6 test as usual:

```bash
./k6 run k8s-test-script.js

...
INFO[0001] 16 Pods found:     source=console
...

```


## APIs


### Create a client:  `new Kubernetes(config)`

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




### `Client.config_maps`

|  Method     |   Description |                               |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-configmap.js](./examples/apply-configmap.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named ConfigMap     |   |
| get         | returns the named ConfigMaps     | [get-configmap.js](./examples/get-configmap.js)  |
| list         | returns a collection of ConfigMaps     | [list-configmaps.js](./examples/list-configmaps.js)   |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

  const nameSpace = "default";
  const name = "config-map-name";
  kubernetesClient.config_maps.apply(getConfigMapYaml(name), nameSpace);
}
```

### `Client.deployments`

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


### `Client.ingresses`

|  Method     |   Description |                               |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-deployment-service-ingress.js](./examples/apply-deployment-service-ingress.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named Ingress     |   |
| get         | returns the named Ingress     | [get-ingress.js](./examples/get-ingress.js)  |
| list         | returns a collection of Ingresses     | [list-ingresses.js](./examples/list-ingresses.js)   |

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

### `Client.jobs`

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

### `Client.namespaces`

|  Method     |   Description |                               |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-namespace.js](./examples/apply-namespace.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named Namespaces     |   |
| get         | returns the named Namespace    | [get-namespace.js](./examples/get-namespace.js)  |
| list         | returns a collection of Namespaces     | [list-namespaces.js](./examples/list-namespaces.js)   |


```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});
  const name = "namespace-name";

  kubernetesClient.namespaces.apply(getNamespaceYaml(name));
}
```

### `Client.nodes`

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

### `Client.persistent_volumes`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-get-delete-pv.js](./examples/apply-get-delete-pv.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named persistent volume     | [apply-get-delete-pv.js](./examples/apply-get-delete-pv.js) |
| get         | returns the named persistent  volume instance     | [apply-get-delete-pv.js](./examples/apply-get-delete-pv.js)  |
| list         | returns a collection of persistent volumens     | [list-pv.js](./examples/list-pv.js)   |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

  const name = "example-pv";
  kubernetesClient.persistent_volumes.apply(getPVYaml(name, "1Gi", "local-storage"));
}
```

### `Client.persistent_volumes_claims`

|  Method     |   Description | Example                              |
| ------------ | ------ | ---------------------------------------- |
| apply         | creates the Kubernetes resource given a YAML configuration    | [apply-get-delete-pvc.js](./examples/apply-get-delete-pvc.js)  |
| create         | creates the Kubernetes resource given an object configuration    |   |
| delete         | removes the named persistent volume claim     | [apply-get-delete-pvc.js](./examples/apply-get-delete-pvc.js) |
| get         | returns the named persistent volume claim     | [apply-get-delete-pvc.js](./examples/apply-get-delete-pvc.js)  |
| list         | returns a collection of persistent volumen claims     | [list-pvc.js](./examples/list-pvc.js)   |

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetesClient = new Kubernetes({});

    const name = "example-pvc";
    const nameSpace = "default";

    kubernetes.persistent_volume_claims.apply(getPVCYaml(name, "1Gi", "nfs-csi"), nameSpace);
}
```

### `Client.pods`

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

### `Client.secrets`

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

### `Client.services`

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
