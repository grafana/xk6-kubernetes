import { Kubernetes } from 'k6/x/kubernetes';
import { sleep } from 'k6';

function getPVCYaml(name, size, storageClass) {
    return `apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ` + name + `
spec:
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: ` + size + `
  storageClassName: ` + storageClass
}

export default function () {
    const kubernetes = new Kubernetes({
        // config_path: "/path/to/kube/config", ~/.kube/config by default
    })
    const name = "example-pvc";
    const nameSpace = "default";

    kubernetes.persistent_volume_claims.apply(getPVCYaml(name, "1Gi", "nfs-csi"), nameSpace);

    sleep(5)

    const pvc_get = kubernetes.persistent_volume_claims.get(name, nameSpace)
    console.log(JSON.stringify(pvc_get, null, 2))

    kubernetes.persistent_volume_claims.delete(name, nameSpace, {});
}
