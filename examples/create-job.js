
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
  const newJob = {
    TypeMeta:{
      ApiVersion: "batch/v1",
      Kind: "Job"
    },
    ObjectMeta: {
      Name: "pi"
    },
    Spec: {
      Template: {
        Spec: {
          Containers: [
            { 
              Name: "pi",
              Image: "perl",
              Command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
            }
          ],
          restartPolicy: "Never"
        }
      },
      BackoffLimit: 4
    }
  }

  kubernetes.jobs.create(nameSpace, newJob)
  console.log(getJobNames(nameSpace, kubernetes))
}