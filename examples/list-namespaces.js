
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nss = kubernetes.namespaces.list()
  console.log(`${nss.length} Namespaces found:`)
  const info = nss.map(function (ns) {
    return {
      name: ns.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
