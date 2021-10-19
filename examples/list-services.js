
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const svcs = kubernetes.services.list()
  console.log(`${svcs.length} Services found:`)
  const info = svcs.map(function (svc) {
    return {
      namespace: svc.namespace,
      name: svc.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
