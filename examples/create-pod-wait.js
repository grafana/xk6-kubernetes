
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
  })
  const namespace = "default"
  const podName = "new-pod"
  const image = "busybox"
  const command = ["sh",  "-c", "/bin/false"]

try {  
  kubernetes.pods.create({
    namespace: namespace,
    name: podName,
    image: image,
    command: command,
    wait: "5s"
  })
  console.log(podName + " has been created")
} catch (err) {
  console.log("error creating pod " + podName + ": " + err)
}
}
