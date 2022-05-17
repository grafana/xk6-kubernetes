
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
  const image = "perl"
  const command = ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]

  kubernetes.jobs.create({
    namespace: namespace,
    name: jobName,
    image: image,
    command: command,
    wait: "30s"
  })
  console.log(jobName + " job completed")  
}

