import { Kubernetes } from 'k6/x/kubernetes';
import {check, sleep} from 'k6';

function getNamespaceYaml(name) {
    return `kind: Namespace
apiVersion: v1
metadata:
  name: ` + name + `
  labels:
    name: ` + name + `
`;
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const name = "namespace-name";

  kubernetes.namespaces.apply(getNamespaceYaml(name));
  const ns_list = kubernetes.namespaces.list().map(function(ns){
      return ns.name;
  })
  sleep(1);
  check(ns_list, {'Namespace was created': (n) => n.indexOf(name) != -1});
}

