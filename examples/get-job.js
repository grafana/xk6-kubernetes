
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "somenamespace"
  const name = "jobname"
  const pod = kubernetes.jobs.get(name, nameSpace)
  console.log(JSON.stringify(pod.object_meta, null, 2))
}
