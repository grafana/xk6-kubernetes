
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  const name = "deploymentname"
  const deployment = kubernetes.deployments.get(name, nameSpace)
  console.log(JSON.stringify(deployment.object_meta, null, 2))
}
