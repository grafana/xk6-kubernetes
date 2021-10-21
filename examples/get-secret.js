
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const name = "secret_name"
  const secret = kubernetes.secrets.get(name, nameSpace)
  console.log(JSON.stringify(secret, null, 2))
}
