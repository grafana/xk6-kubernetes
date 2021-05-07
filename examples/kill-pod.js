
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

function getPodNames(nameSpace, kubernetes) {
  return kubernetes.pods.list(nameSpace).map(function(pod){
    return pod.name
  })
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  let podsList = getPodNames(nameSpace, kubernetes)
  const podName = podsList[0]
  kubernetes.pods.kill(podName, nameSpace)
  sleep(1)
  if (kubernetes.pods.isTerminating(podName, nameSpace)) {
    console.log(podName + " pod has been killed successfully")
  } else {
    throw `${podName} Pod was not killed ${podsList[podName]}`
  }
}