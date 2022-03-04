
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const pvs = kubernetes.persistent_volumes.list()
  console.log(`${pvs.length} PVs found:`)
  const info = pvs.map(function (pv) {
    return {
      name: pv.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
