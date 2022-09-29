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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

// setupTestEnv should be called from each test to build the execution environment for the test
func setupTestEnv(t *testing.T, objs ...runtime.Object) *goja.Runtime {
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

	m.clientset = fake.NewSimpleClientset(objs...)

	dynamic, err := localutils.NewFakeDynamic()
	if err != nil {
		t.Errorf("unexpected error creating fake client %v", err)
	}
	m.dynamic = dynamic

	return rt
}

// TestGenericApiIsScriptable runs through creating, getting, listing and deleting an object
func TestGenericApiIsScriptable(t *testing.T) {
	t.Parallel()

	rt := setupTestEnv(t)

	_, err := rt.RunString(`
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
