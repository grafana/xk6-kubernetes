
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })

  const nodes = kubernetes.nodes.list();
  const nodeName = nodes[0].namespace;
  const namespace = "default"
  const jobName = "new-nodename-job"
  const image = "perl"
  const command = ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]

  kubernetes.jobs.create({
    namespace: namespace,
    name: jobName,
    node_name: nodeName,
    image: image,
    command: command
  })
  sleep(1)

  const job = kubernetes.jobs.list(namespace).find((job) => {
    return job.spec.template.spec.node_name === nodeName
  });

  if(job) {
    console.log(job.name + " job has been created successfully")
  } else {
    throw jobName + " job was not created"
  }
}
