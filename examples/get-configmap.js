
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const name = "configMap_name"
  const configMap = kubernetes.config_maps.get(name, nameSpace)
  console.log(JSON.stringify(configMap, null, 2))
}
