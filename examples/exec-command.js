
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
  })
  const namespace = "default"
  const podName = "new-pod"
  const image = "busybox"
  const command = ["sh",  "-c", "sleep 300"]

  kubernetes.pods.create({
    namespace: namespace,
    name: podName,
    image: image,
    command: command
  })
  sleep(3)

  const newPod = kubernetes.pods.list(namespace).find(function(pod) { return pod.name == podName}) 
  if (!newPod) {
    throw podName + " pod was not created"
  }

  const container = newPod.spec.containers[0].name
  const result = kubernetes.pods.exec({
    namespace: namespace,
    pod: podName,
    container: container.name,
    command: ["echo", "'hello xk6'"],
    stadin:  []
  })

  const stdout = String.fromCharCode(...result.stdout)

  if (stdout.includes("xk6")) {
    console.log("command executed")
  } else {
    throw "command not executed correctly"
  }
}
