import { Kubernetes } from 'k6/x/kubernetes';

let pod = {
    apiVersion: "v1",
    kind:       "Pod",
    metadata: {
        name:      "busybox",
        namespace:  "default"
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

  const ns = "default"

  // create pod in test namespace
  pod.metadata.namespace = ns
  kubernetes.create(pod)

  // get helpers for test namespace
  const helpers = kubernetes.helpers(ns)

  // wait for pod to be running
  const timeout = 10
  if (!helpers.waitPodRunning(pod.metadata.name, timeout)) {
      console.log(`"pod ${pod.metadata.name} not ready after ${timeout} seconds`)
  }
}