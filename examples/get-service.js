
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const name = "svc_name"
  const svc = kubernetes.services.get(name, nameSpace)
  console.log(JSON.stringify(svc, null, 2))
}
