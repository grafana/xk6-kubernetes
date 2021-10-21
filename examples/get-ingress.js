
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const name = "ing_name"
  const pod = kubernetes.ingresses.get(name, nameSpace)
  console.log(JSON.stringify(pod, null, 2))
}
