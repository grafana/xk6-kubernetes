import { Kubernetes } from 'k6/x/kubernetes';
import {check, sleep} from 'k6';

function getConfigMapYaml(name) {
    return `kind: ConfigMap 
apiVersion: v1 
metadata:
  name: ` + name + `
data:
  configkey: configvalue

  configfile: | 
    configproperty1=42
    configproperty2=foo
`
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default";
  const name = "config-map-name";

  kubernetes.config_maps.apply(getConfigMapYaml(name), nameSpace);
  const cm_list = kubernetes.config_maps.list(nameSpace).map(function(cm){
      return cm.name;
  })
  sleep(1);
  check(cm_list, {'ConfigMap was created': (c) => c.indexOf(name) != -1});
}

