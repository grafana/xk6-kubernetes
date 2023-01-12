import { Kubernetes } from 'k6/x/kubernetes';
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.2/index.js';

let json = {
    apiVersion: "networking.k8s.io/v1",
    kind: "Ingress",
    metadata: {
        name:      "json-ingress",
        namespace: "default",
    },
    spec: {
        ingressClassName: "nginx",
        rules: [
            {
                http: {
                    paths: [
                        {
                            path: "/my-service-path",
                            pathType: "Prefix",
                            backend: {
                                service: {
                                    name: "my-service",
                                    port: {
                                        number: 80
                                    }
                                }
                            }
                        }
                    ]
                }
            }
        ]
    }
}

let yaml = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: yaml-ingress
  namespace: default
spec:
  ingressClassName: nginx
  rules:
  - http:
      paths:
      - path: /my-service-path
        pathType: Prefix
        backend:
          service:
            name: my-service
            port:
              number: 80
`

export default function () {
    const kubernetes = new Kubernetes();

    describe('JSON-based resources', () => {
        const name = json.metadata.name
        const ns = json.metadata.namespace

        let ingress

        describe('Create our ingress using the JSON definition', () => {
            ingress = kubernetes.create(json)
            expect(ingress.metadata).to.have.property('uid')
        })

        describe('Retrieve all available ingresses', () => {
            expect(kubernetes.list("Ingress.networking.k8s.io", ns).length).to.be.at.least(1)
        })

        describe('Retrieve our ingress by name and namespace', () => {
            let fetched = kubernetes.get("Ingress.networking.k8s.io", name, ns)
            expect(ingress.metadata.uid).to.equal(fetched.metadata.uid)
        })

        describe('Update our ingress with a modified JSON definition', () => {
            const newValue = json.spec.rules[0].http.paths[0].path + '-updated'
            json.spec.rules[0].http.paths[0].path = newValue

            kubernetes.update(json)
            let updated = kubernetes.get("Ingress.networking.k8s.io", name, ns)
            expect(updated.spec.rules[0].http.paths[0].path).to.be.equal(newValue)
        })

        describe('Remove our ingresses to cleanup', () => {
            kubernetes.delete("Ingress.networking.k8s.io", name, ns)
        })
    })

    describe('YAML-based resources', () => {
        const name = 'yaml-ingress'
        const ns = 'default'

        let uid

        describe('Create our ingress using the YAML definition', () => {
            kubernetes.apply(yaml)
            let created = kubernetes.get("Ingress.networking.k8s.io", name, ns)
            expect(created.metadata).to.have.property('uid')
            uid = created.metadata.uid
        })

        describe('Update our ingress with a modified YAML definition', () => {
            kubernetes.apply(yaml)
            let updated = kubernetes.get("Ingress.networking.k8s.io", name, ns)
            expect(updated.metadata.uid).to.be.equal(uid)
        })

        describe('Remove our ingresses to cleanup', () => {
            kubernetes.delete("Ingress.networking.k8s.io", name, ns)
        })
    })

}
