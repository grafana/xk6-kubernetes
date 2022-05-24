
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

function getJobNames(nameSpace, kubernetes) {
  return kubernetes.jobs.list(nameSpace).map(function(job){
    return job.name
  })
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const namespace = "default"
  const jobName = "new-job"
  const image = "busybox"
  const command = ["sh",  "-c", "sleep 3"]

  kubernetes.jobs.create({
    namespace: namespace,
    name: jobName,
    image: image,
    command: command,
    wait: "10s",
    autodelete: true
  })
  console.log("job completed")
  sleep(3)  // wait for garbage collection
  const jobPod = kubernetes.pods.list().find(pod => pod.name.startsWith(jobName))
  console.log(jobPod)
  if (!jobPod) {
     console.log("pods deleted")
  }
}
