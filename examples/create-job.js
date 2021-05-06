
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
    command: command
  })
  sleep(1)
  const jobsList = getJobNames(namespace, kubernetes)
  if(jobsList.indexOf(jobName) != -1) {
    console.log(jobName + " job has been created successfully")
  } else {
    throw jobName + " job was not created"
  }
}
