import { Kubernetes } from 'k6/x/kubernetes';

const podSpec = {
    apiVersion: "v1",
    kind:       "Pod",
    metadata: {
        name:      "busybox",
        namespace: "testns"
    },
    spec: {
        containers: [
            {
                name:    "busybox",
                image:   "busybox",
                command: ["sh", "-c", "sleep 30"]
            }
        ]
    }
}

export default function () {
  const kubernetes = new Kubernetes();

  kubernetes.create(podSpec)

  const pods = kubernetes.list("Pod", "testns");

  console.log(`${pods.length} Pods found:`);
  pods.map(function(pod) {
    console.log(`  ${pod.metadata.name}`)
  });
}