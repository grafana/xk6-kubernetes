import { Kubernetes } from 'k6/x/kubernetes';
import { sleep } from 'k6';

function getPVYaml(name, size, storageClass) {
    return `apiVersion: v1
kind: PersistentVolume
metadata:
  name: ` + name + `
spec:
  capacity:
    storage: ` + size + `
  volumeMode: Filesystem
  accessModes:
  - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: ` + storageClass
}

export default function () {
    const kubernetes = new Kubernetes({
        // config_path: "/path/to/kube/config", ~/.kube/config by default
    })
    const name = "example-pv";

    kubernetes.persistent_volumes.apply(getPVYaml(name, "1Gi", "local-storage"));

    sleep(5)

    const pv_get = kubernetes.persistent_volumes.get(name)
    console.log(JSON.stringify(pv_get, null, 2))

    kubernetes.persistent_volumes.delete(name, {});
}
