import { Kubernetes } from 'k6/x/kubernetes'
import { sleep } from 'k6'

const k8s = new Kubernetes()

const podSpec = {
    apiVersion: "v1",
    kind:       "Pod",
    metadata: {
        name:      "busybox",
        namespace: "default"
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

export default function(){
        var created = k8s.create(podSpec)
        console.log("pod '" + created.metadata.name +"' created")

        var pod = k8s.get(podSpec.kind, podSpec.metadata.name, podSpec.metadata.namespace)
        if (podSpec.metadata.name != pod.metadata.name) {
                throw new Error("Fetch by name did not return the Service. Expected: " + podSpec.metadata.name + " but got: " + fetched.name)
        }

        const pods = k8s.list(podSpec.kind, podSpec.metadata.namespace)
        if (pods === undefined || pods.length < 1) {
                throw new Error("Expected listing with 1 Pod")
        }

        k8s.delete(podSpec.kind, podSpec.metadata.name, podSpec.metadata.namespace)
        
        // give time for the pod to be deleted
        sleep(5)

        if (k8s.list(podSpec.kind, podSpec.metadata.namespace).length != 0) {
                throw new Error("Deletion failed to remove pod")
        }
}