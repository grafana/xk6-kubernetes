import { Kubernetes } from 'k6/x/kubernetes';
import {check, sleep} from 'k6';

function getDeploymentYaml(name, app) {
    return `kind: Deployment
apiVersion: apps/v1
metadata:
  name: ` + name + `
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ` + app +`
  
  template:
    metadata:
      labels:
        app: ` + app + `
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
`
}

function getServiceYaml(name, app) {
    return `apiVersion: v1
kind: Service
metadata:
  name: ` + name + `
spec:
  selector:
    app: ` + app + `
  type: ClusterIP
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
      name: http
`
}

function getIngressYaml(name, url) {
    return `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ` + name + `
  annotations:
    kuberentes.io/ingress.class: "nginx"
spec:
  rules:
    - host: ` + url + `
      http:
        paths:
        - path: /
          pathType: ImplementationSpecific
          backend:
            service:
              name: ` + name + `
              port:
                name: http
`
}

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nameSpace = "default";
  const name = "deployment-name";
  const app = 'app-label';
  const url = 'ingress-url.com';

  kubernetes.deployments.apply(getDeploymentYaml(name, app), nameSpace);
  kubernetes.services.apply(getServiceYaml(name, app), nameSpace);
  kubernetes.ingresses.apply(getIngressYaml(name, url), nameSpace);
  sleep(1);

  const depl_list = kubernetes.deployments.list(nameSpace).map(function(depl){
      return depl.name;
  })
  check(depl_list, {'Deployment was created': (d) => d.indexOf(name) != -1});
  const serv_list = kubernetes.services.list(nameSpace).map(function(serv){
      return serv.name;
  })
  check(serv_list, {'Service was created': (s) => s.indexOf(name) != -1});
  const ing_list = kubernetes.ingresses.list(nameSpace).map(function(ing){
      return ing.name;
  })
  check(ing_list, {'Ingress was created': (i) => i.indexOf(name) != -1});
}

