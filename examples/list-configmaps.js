
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const configMaps = kubernetes.config_maps.list()
  console.log(`${configMaps.length} ConfigMaps found:`)
  const info = configMaps.map(function (configMap) {
    return {
      namespace: configMap.namespace,
      name: configMap.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
