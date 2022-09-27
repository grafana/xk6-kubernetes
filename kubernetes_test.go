package kubernetes

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/dop251/goja"
	localutils "github.com/grafana/xk6-kubernetes/internal/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"

	"k8s.io/client-go/kubernetes/fake"
)

type testEnv struct {
	Runtime *goja.Runtime
	Module  *ModuleInstance
}

// setupTestEnv should be called from each test to build the execution environment for the test
func setupTestEnv(t *testing.T) testEnv {
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper{})

	testLog := logrus.New()
	testLog.AddHook(&testutils.SimpleLogrusHook{
		HookedLevels: []logrus.Level{logrus.WarnLevel},
	})
	testLog.SetOutput(ioutil.Discard)

	state := &lib.State{
		Options: lib.Options{
			SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
		},
		Logger: testLog,
		Tags:   lib.NewTagMap(nil),
	}

	root := &RootModule{}
	m, ok := root.NewModuleInstance(
		&modulestest.VU{
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment{},
			CtxField:     context.Background(),
			StateField:   state,
		},
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("Kubernetes", m.Exports().Named["Kubernetes"]))

	return testEnv{
		Module:  m,
		Runtime: rt,
	}
}

// TestConfigMapsScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestConfigMapsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewConfigMap("cm-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const cms = k8s.config_maps.list("ns-test")
if (cms === undefined || cms.length < 1) {
	throw new Error("Expected listing with at least 1 ConfigMap")
}
const initialCount = cms.length

var created = k8s.config_maps.create({name: "cm-new", data: {"key-1": "value-1", "key-2": "value-2"}}, "ns-test")
if (k8s.config_maps.list(created.namespace).length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.config_maps.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the ConfigMap. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.config_maps.delete(created.name, created.namespace)
if (k8s.config_maps.list("ns-test").length != initialCount) {
	throw new Error("Deletion failed to remove ConfigMap from list")
}
`)
	require.NoError(t, err)
}

// TestDeploymentsScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestDeploymentsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewDeployment("deployment-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const deployments = k8s.deployments.list("ns-test")
if (deployments === undefined || deployments.length < 1) {
	throw new Error("Expected listing with at least 1 Deployment")
}
const initialCount = deployments.length

var created = k8s.deployments.create(
  {name: "deployment-new", spec: {replicas: 3, template: {spec: {containers: [{name: "nginx", image: "nginx:1.14.2"}]}}}}, 
  "ns-test"
)
if (k8s.deployments.list(created.namespace).length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.deployments.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Deployment. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.deployments.delete(created.name, created.namespace)
if (k8s.deployments.list("ns-test").length != initialCount) {
	throw new Error("Deletion failed to remove Deployment from list")
}
`)
	require.NoError(t, err)
}

// TestIngressesScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestIngressesScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewIngress("ingress-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const ingresses = k8s.ingresses.list("ns-test")
if (ingresses === undefined || ingresses.length < 1) {
	throw new Error("Expected listing with at least 1 Ingress")
}
const initialCount = ingresses.length

var created = k8s.ingresses.create(
  {name: "ingress-new", spec: {ingressClassname: "nginx-example", rules: [{http: {paths: [{path: "/testpath"}]}}]}}, 
  "ns-test"
)
if (k8s.ingresses.list(created.namespace).length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.ingresses.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Ingress. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.ingresses.delete(created.name, created.namespace)
if (k8s.ingresses.list("ns-test").length != initialCount) {
	throw new Error("Deletion failed to remove Ingress from list")
}
`)
	require.NoError(t, err)
}

// TestJobsScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestJobsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewJob("job-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const jobs = k8s.jobs.list("ns-test")
if (jobs === undefined || jobs.length < 1) {
	throw new Error("Expected listing with at least 1 Job")
}
const initialCount = jobs.length

var created = k8s.jobs.create(
  {name: "job-new", node_name: "node-1", image: "perl", command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]}, 
  "ns-test"
)
if (k8s.jobs.list(created.namespace).length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.jobs.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Job. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.jobs.delete(created.name, created.namespace)
if (k8s.jobs.list("ns-test").length != initialCount) {
	throw new Error("Deletion failed to remove Job from list")
}
`)
	require.NoError(t, err)
}

// TestNamespaceScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestNamespaceScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewNamespace("ns-1"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const namespaces = k8s.namespaces.list()
if (namespaces === undefined || namespaces.length < 1) {
	throw new Error("Expected listing with at least 1 Namespace")
}
const initialCount = namespaces.length

var created = k8s.namespaces.create({name: "ns-new"})
if (k8s.namespaces.list().length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.namespaces.get(created.name)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Namespace. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.namespaces.delete(created.name)
if (k8s.namespaces.list().length != initialCount) {
	throw new Error("Deletion failed to remove Namespace from list")
}
`)
	require.NoError(t, err)
}

// TestNodesScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestNodesScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewNodes("node-1"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const nodes = k8s.nodes.list()
if (nodes === undefined || nodes.length < 1) {
	throw new Error("Expected listing with at least 1 Node")
}
`)
	require.NoError(t, err)
}

// TestPersistentVolumeClaimsScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestPersistentVolumeClaimsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewPersistentVolumeClaim("pvc-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const pvcs = k8s.persistent_volume_claims.list()
if (pvcs === undefined || pvcs.length < 1) {
	throw new Error("Expected listing with at least 1 PersistentVolumeClaim")
}
const initialCount = pvcs.length

var created = k8s.persistent_volume_claims.create(
  {name: "pvc-new", spec: {storageClassName: "local-storage"}},
  "ns-test")
if (k8s.persistent_volume_claims.list().length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.persistent_volume_claims.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the PersistentVolumeClaim. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.persistent_volume_claims.delete(created.name, created.namespace)
if (k8s.persistent_volume_claims.list().length != initialCount) {
	throw new Error("Deletion failed to remove PersistentVolumeClaim from list")
}
`)
	require.NoError(t, err)
}

// TestPersistentVolumesScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestPersistentVolumesScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewPersistentVolume("pv-1"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const pvs = k8s.persistent_volumes.list()
if (pvs === undefined || pvs.length < 1) {
	throw new Error("Expected listing with at least 1 PersistentVolume")
}
const initialCount = pvs.length

var created = k8s.persistent_volumes.create({name: "pv-new", spec: {storageClassName: "local-storage", hostPath: {path: "/tmp/xk6-test"}}})
if (k8s.persistent_volumes.list().length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.persistent_volumes.get(created.name)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the PersistentVolume. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.persistent_volumes.delete(created.name)
if (k8s.persistent_volumes.list().length != initialCount) {
	throw new Error("Deletion failed to remove PersistentVolume from list")
}
`)
	require.NoError(t, err)
}

// TestPodsScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestPodsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewPod("pod-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const pods = k8s.pods.list()
if (pods === undefined || pods.length < 1) {
	throw new Error("Expected listing with at least 1 Pod")
}
const initialCount = pods.length

var created = k8s.pods.create(
  {name: "pod-new", image: "busybox", command: ["sh", "-c", "sleep 300"], restart_policy: "Never"},
  "ns-test")
if (k8s.pods.list().length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.pods.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Pod. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.pods.delete(created.name, created.namespace)
if (k8s.pods.list().length != initialCount) {
	throw new Error("Deletion failed to remove Pod from list")
}
`)
	require.NoError(t, err)
}

// TestSecretsScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestSecretsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewSecret("secret-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const secrets = k8s.secrets.list()
if (secrets === undefined || secrets.length < 1) {
	throw new Error("Expected listing with at least 1 Secret")
}
const initialCount = secrets.length

var created = k8s.secrets.create(
  {name: "secret-new", data: {"secret-key": "MWYyZDFlMmU2N2Rm"}},
  "ns-test")
if (k8s.secrets.list().length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.secrets.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Secret. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.secrets.delete(created.name, created.namespace)
if (k8s.secrets.list().length != initialCount) {
	throw new Error("Deletion failed to remove Secret from list")
}
`)
	require.NoError(t, err)
}

// TestServicesScriptable runs through listing, creating, fetching, and deleting to ensure scripting
func TestServicesScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)
	tenv.Module.clientset = fake.NewSimpleClientset(
		localutils.NewService("svc-1", "ns-test"),
	)

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const svcs = k8s.services.list()
if (svcs === undefined || svcs.length < 1) {
	throw new Error("Expected listing with at least 1 Service")
}
const initialCount = svcs.length

var created = k8s.services.create(
  {name: "svc-new", spec: {selector: {app: "MyApp"}, ports: [{protocol: "TCP", port: 80}]}},
  "ns-test")
if (k8s.services.list().length != (initialCount + 1)) {
	throw new Error("Expected list to have increased given addition")
}

var fetched = k8s.services.get(created.name, created.namespace)
if (created.name != fetched.name) {
	throw new Error("Fetch by name did not return the Service. Expected: " + created.name + " but got: " + fetched.name)
}

k8s.services.delete(created.name, created.namespace)
if (k8s.services.list().length != initialCount) {
	throw new Error("Deletion failed to remove Service from list")
}
`)
	require.NoError(t, err)
}

// TestGenericApiIsScriptable runs through creating, getting, listing and deleting an object
func TestGenericApiIsScriptable(t *testing.T) {
	t.Parallel()

	tenv := setupTestEnv(t)

	tenv.Module.clientset = fake.NewSimpleClientset()
	tenv.Module.dynamic = localutils.NewFakeDynamic()

	_, err := tenv.Runtime.RunString(`
const k8s = new Kubernetes()

const podSpec = {
    apiVersion: "v1",
    kind:       "Pod",
    metadata: {
        name:      "busybox",
        namespace: "testns"
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

var created = k8s.create(podSpec)

var pod = k8s.get(podSpec.kind, podSpec.metadata.name, podSpec.metadata.namespace)
if (podSpec.metadata.name != pod.metadata.name) {
	throw new Error("Fetch by name did not return the Service. Expected: " + podSpec.metadata.name + " but got: " + fetched.name)
}

const pods = k8s.list(podSpec.kind, podSpec.metadata.namespace)
if (pods === undefined || pods.length < 1) {
	throw new Error("Expected listing with 1 Pod")
}

k8s.delete(podSpec.kind, podSpec.metadata.name, podSpec.metadata.namespace)
if (k8s.list(podSpec.kind, podSpec.metadata.namespace).length != 0) {
	throw new Error("Deletion failed to remove pod")
}
`)
	require.NoError(t, err)
}
