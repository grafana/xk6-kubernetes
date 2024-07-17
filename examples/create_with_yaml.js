import { Kubernetes } from 'k6/x/kubernetes';

const manifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: busybox
  namespace: testns
spec:
  template:
    spec:
      containers:
      - name: busybox
        image: busybox
        command: ["sleep", "300"]
    restartPolicy: Never
`

export default function () {
  const kubernetes = new Kubernetes();

  kubernetes.apply(manifest)

  const jobs = kubernetes.list("Job", "testns");

  console.log(`${jobs.length} Jobs found:`);
  pods.map(function(job) {
    console.log(`  ${job.metadata.name}`)
  });
}