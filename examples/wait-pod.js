
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
  })
  const namespace = "default"
  const podName = "new-pod"
  const image = "busybox"
  const command = ["sh",  "-c", "sleep 5"]

  kubernetes.pods.create({
    namespace: namespace,
    name: podName,
    image: image,
    command: command
  })
  
 
  const options = {
    namespace: namespace,
    name: podName,
    status: "Succeeded",
    timeout: "10s"
  }
  if (kubernetes.pods.wait(options)) {
    console.log(podName + " pod completed successfully")
  } else {
    throw podName + " is not completed"
  }
}
