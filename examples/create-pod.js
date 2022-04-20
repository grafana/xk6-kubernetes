
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

function getPodNames(nameSpace, kubernetes) {
  return kubernetes.pods.list(nameSpace).map(function(pod){
    return pod.name
  })
}

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
  sleep(1)
  const podsList = getPodNames(namespace, kubernetes)
  if(podsList.indexOf(podName) != -1) {
    console.log(podName + " pod has been created successfully")
  } else {
    throw podName + " pod was not created"
  }
}
