
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const secrets = kubernetes.secrets.list()
  console.log(`${secrets.length} Secrets found:`)
  const info = secrets.map(function (secret) {
    return {
      namespace: secret.namespace,
      name: secret.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
