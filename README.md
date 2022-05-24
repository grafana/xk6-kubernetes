[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/xk6-kubernetes.svg)](https://pkg.go.dev/github.com/grafana/xk6-kubernetes)
[![Version Badge](https://img.shields.io/github/v/release/grafana/xk6-kubernetes?style=flat-square)](https://github.com/grafana/xk6-kubernetes/releases)
![Build Status](https://img.shields.io/github/workflow/status/grafana/xk6-kubernetes/CI?style=flat-square)

# xk6-kubernetes
A k6 extension for interacting with Kubernetes clusters while testing.

## Build

To build a custom `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  xk6 build --with github.com/grafana/xk6-kubernetes
  ```

## Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```bash
make
```
Once built, you can run your newly extended `k6` using:
```shell
 ./k6 run my-test-script.js
 ```

## Example

```javascript
import { Kubernetes } from 'k6/x/kubernetes';
  
export default function () {
    const kubernetes = new Kubernetes({
        // config_path: "/path/to/kube/config", ~/.kube/config by default
      })
    const pods = kubernetes.pods.list()
    console.log(`${kubernetes.pods.list().length} Pods found:`)
    pods.map(function(pod) {
        console.log(`  ${pod.name}`)
    })
}
```

Result output:

```plain
$ ./k6 run my-test-script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: my-test-script.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0001] 16 Pods found:                                source=console
... snipped for brevity ...

running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=9.64ms min=9.64ms med=9.64ms max=9.64ms p(90)=9.64ms p(95)=9.64ms
     iterations...........: 1   25.017512/s

```

Inspect [examples](./examples) folder for more details.


## If things go wrong

### Are you using the custom binary?
An easy mistake--which happens often--is to forget that `xk6` is generating a new executable. You may be accustomed to simply running `k6` from the command-line which probably isn't your new build. Make sure to use `./k6` after building your extended version otherwise you can expect to see an error similar to:

```bash
ERRO[0000] The moduleSpecifier "my-test-script.js" couldn't be found on local disk. Make sure that you've specified the right path to the file. If you're running k6 using the Docker image make sure you have mounted the local directory (-v /local/path/:/inside/docker/path) containing your script and modules so that they're accessible by k6 from inside of the container, see https://k6.io/docs/using-k6/modules#using-local-modules-with-docker. Additionally it was tried to be loaded as remote module by prepending "https://" to it, which also didn't work. Remote resolution error: "Get "https://my-test-script.js": dial tcp: lookup my-test-script.js: no such host" 
```
