
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const pods = kubernetes.jobs.list('ingress-nginx')
  console.log(`${pods.length} Jobs found:`)
  const info = pods.map(function(job){
    return job.name
  })
  console.log(JSON.stringify(info, null, 2))
}
