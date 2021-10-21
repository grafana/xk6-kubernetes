
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const ingresses = kubernetes.pods.list()
  console.log(`${ingresses.length} Ingresses found:`)
  const info = ingresses.map(function (ingress) {
    return {
      namespace: ingress.namespace,
      name: ingress.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
