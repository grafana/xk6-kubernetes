
import { Kubernetes } from 'k6/x/kubernetes';

export default function () {
  const kubernetes = new Kubernetes({
    // config_path: "/path/to/kube/config", ~/.kube/config by default
  })
  const nodes = kubernetes.nodes.list()
  console.log(`${nodes.length} Nodes found:`)
  const info = nodes.map(function(node){
    let conditions = {}
    node.status.conditions.forEach(function(condition) {
	conditions[condition.type] = condition.status
    })
    return {
      name: node.name,
      conditions: conditions
    } 
  })
  console.log(JSON.stringify(info, null, 2))
}
