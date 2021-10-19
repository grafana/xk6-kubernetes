
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const deployments = kubernetes.deployments.list()
  console.log(`${deployments.length} Deployments found:`)
  const info = deployments.map(function (deployment) {
    return {
      namespace: deployment.namespace,
      name: deployment.name,
    }
  })
  console.log(JSON.stringify(info, null, 2))
}
