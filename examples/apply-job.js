import { Kubernetes } from 'k6/x/kubernetes';
import {check, sleep} from 'k6';

function getJobYaml(name) {
    return `apiVersion: batch/v1
kind: Job
metadata:
  name: ` + name + `
spec:
  template:
    spec:
      containers:
      - name: ` + name + `
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
`
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default";
  const name = "job-name";

  kubernetes.jobs.apply(getJobYaml(name), nameSpace);
  const job_list = kubernetes.jobs.list(nameSpace).map(function(job){
      return job.name;
  })
  sleep(1);
  check(job_list, {'Job was created': (job) => job.indexOf(name) != -1});
}

