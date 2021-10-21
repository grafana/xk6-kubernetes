
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const name = "default"
  const ns = kubernetes.namespaces.get(name)
  console.log(JSON.stringify(ns, null, 2))
}
