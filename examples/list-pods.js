
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const pods = kubernetes.pods.list()
  console.log(`${pods.length} Pods found:`)
  const info = pods.map(function(pod){
    return {
      namespace: pod.namespace,
      name: pod.name,
      status: pod.status.phase
    } 
  })
  console.log(JSON.stringify(info, null, 2))
}
