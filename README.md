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
| apply         | manifest string| creates a Kubernetes resource given a YAML manifest |
| create         | spec object | creates a Kubernetes resource given its specification |
| delete         | kind  | removes the named resource |
|                | name  |
|                | namespace|
| get         | kind| returns the named resource |
|                | name  |
|                | namespace |
| list         | kind| returns a collection of resources of a given kind
|                | namespace |


The kinds of resources currently supported are:
* ConfigMap
* Deployment
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


## If things go wrong


### Are you using the custom binary?

  An easy mistake--which happens often--is to forget that `xk6` is generating a new executable. You may be accustomed to simply running `k6` from the command-line which probably isn't your new build. Make sure to use `./k6` after building your extended version otherwise you can expect to see an error similar to:


  ```bash
  ERRO[0000] The moduleSpecifier "k8s-test-script.js" couldn't be found on local disk. Make sure that you've specified the right path to the file. If you're running k6 using the Docker image make sure you have mounted the local directory (-v /local/path/:/inside/docker/path) containing your script and modules so that they're accessible by k6 from inside of the container, see https://k6.io/docs/using-k6/modules#using-local-modules-with-docker. Additionally it was tried to be loaded as remote module by prepending "https://" to it, which also didn't work. Remote resolution error: "Get "https://k8s-test-script.js": dial tcp: lookup k8s-test-script.js: no such host" 
  ```
