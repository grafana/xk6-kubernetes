
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const pvcs = kubernetes.persistent_volume_claims.list(nameSpace)
  console.log(`${pvcs.length} PVCs found:`)
  const info = pvcs.map(function (pvc) {
    return {
      name: pvc.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
