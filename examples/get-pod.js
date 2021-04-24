
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const name = "pod_name"
  const pod = kubernetes.pods.get(name, nameSpace)
  console.log(JSON.stringify(pod, null, 2))
}