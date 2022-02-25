import { Kubernetes } from 'k6/x/kubernetes';
import {check, sleep} from 'k6';

function getSecretYaml(name) {
    return `apiVersion: v1
kind: Secret
metadata:
  name: ` + name + `
type: Opaque
data:
  secretkey: c2VjcmV0dmFsdWUK
`
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default";
  const name = "secret-name";

  kubernetes.secrets.apply(getSecretYaml(name), nameSpace);
  const sc_list = kubernetes.secrets.list(nameSpace).map(function(sc){
      return sc.name;
  })
  sleep(1);
  check(sc_list, {'Secret was created': (s) => s.indexOf(name) != -1});
}

