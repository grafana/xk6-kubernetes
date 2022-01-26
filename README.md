> ### ⚠️ This is a proof of concept
>
> As this is a proof of concept,  it won't be supported by the k6 team.
> It may also break in the future as xk6 evolves. USE AT YOUR OWN RISK!
> Any issues with the tool should be raised [here](https://github.com/grafana/xk6-kubernetes/issues).

</br>
</br>

<div align="center">

# xk6-kubernetes
A k6 extension for interacting with Kubernetes clusters while testing. Built for [k6](https://github.com/loadimpact/k6) using [xk6](https://github.com/grafana/xk6).

</div>

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  $ xk6 build --with github.com/grafana/xk6-kubernetes
  ```

## Development

To make development a little smoother, you may run the build script provided in the root folder. It will create a k6 binary with your local code rather than from GitHub.

```bash
$ ./build.sh && ./k6 run my-test-script.js
```

## Example

```javascript
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  kubernetes = new Kubernetes()
  console.log(`${kubernetes.pods.list()} Pods found:`)
}
```

Result output:

```plain
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: ../xk6-kubernetes/script.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0001] 16 Pods found:                                source=console

running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=9.64ms min=9.64ms med=9.64ms max=9.64ms p(90)=9.64ms p(95)=9.64ms
     iterations...........: 1   25.017512/s

```

Inspect examples folder for more details.
