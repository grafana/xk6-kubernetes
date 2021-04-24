
import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

function getPodNames(nameSpace, kubernetes) {
  const result = {}
  kubernetes.pods.list(nameSpace).forEach(function(pod){
    result[pod.name] = pod.status.phase
  })
  return result
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default"
  let podsList = getPodNames(nameSpace, kubernetes)
  const podName = Object.keys(podsList)[0]
  kubernetes.pods.kill(podName, nameSpace)
  // TODO: for some reason we can't feel terminating status
  // so this example just checks if pod was removed 1min after
  sleep(60) 
  podsList = getPodNames(nameSpace, kubernetes)

  if (!podsList[podName]) {
    console.log `${podName} pod has been killed successfully`
  } else {
    throw `${podName} POD WAS NOT KILLED! ${podsList[podName]}`
  }
}