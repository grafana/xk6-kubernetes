
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const pods = kubernetes.pods.list()
  console.log(`${pods.length} Pods found:`)
  const names = pods.map(function(pod){
    return pod.object_meta.name
  })
  console.log(JSON.stringify(names, null, 2))
}