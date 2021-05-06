
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
  const nameSpace = "default"
  const jobName = "new-job"
  kubernetes.jobs.kill(jobName, nameSpace)
  sleep(1)
  let jobsList = getJobNames(nameSpace, kubernetes)
  if(jobsList.indexOf(jobName) == -1) {
    console.log(jobName + " job has been killed successfully")
  } else {
    throw jobName + " job was not killed"
  }
}
