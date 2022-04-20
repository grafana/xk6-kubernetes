import { sleep } from 'k6';
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const namespace = "default"
  const podName = "new-pod"
  const image = "busybox"
  const command = ["sh",  "-c", "sleep 300"]
  const containerName = "ephemeral"
  const containerImage = "busybox" 
  const containerCommand = ["sh", "-c", "sleep 300"]

  kubernetes.pods.create({
    namespace: namespace,
    name: podName,
    image: image,
    command: command
  })
  sleep(1)

  kubernetes.pods.addEphemeralContainer(
    podName,
    namespace,
    {
      name: containerName,
      image: containerImage,
      command: containerCommand,
    }   
  )
  sleep(1)

  let pod = kubernetes.pods.get(podName, namespace)
  if (pod.spec.ephemeral_containers[0].name == containerName) {
    console.log(containerName + " container successfully created")
  } else {
    throw containername + " container not created"
  }
}
