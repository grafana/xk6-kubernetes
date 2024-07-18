package kubernetes

import (
	"context"
	"io"
	"testing"

	"github.com/grafana/sobek"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/k3s"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// setupE2eTestEnv should be called from each test to build the execution environment for e2e tests
func setupE2eTestEnv(t *testing.T) *sobek.Runtime {
	rt := sobek.New()
	rt.SetFieldNameMapper(common.FieldNameMapper{})

	testLog := logrus.New()
	testLog.AddHook(&testutils.SimpleLogrusHook{
		HookedLevels: []logrus.Level{logrus.WarnLevel},
	})
	testLog.SetOutput(io.Discard)

	state := &lib.State{
		Options: lib.Options{
			SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
		},
		Logger: testLog,
		Tags:   lib.NewVUStateTags(metrics.NewRegistry().RootTagSet()),
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

	ctx := context.Background()

	// start a k3s container
	k3sContainer, err := k3s.RunContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container
	defer func() {
		if err := k3sContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()
	
	kubeConfigYaml, err := k3sContainer.GetKubeConfig(context.TODO())
	if err != nil {
		t.Fatalf("failed to get kubeconfig: %s", err)
	}

	restcfg, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigYaml)
	if err != nil {
		t.Fatalf("failed to create rest config: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(restcfg)
	if err != nil {
		t.Fatalf("failed to create k8s client: %s", err)
	}

	m.clientset = clientset
	m.config = restcfg

	return rt
}

// E2eTest runs helpers
func TestE2e(t *testing.T) {
	t.Parallel()

	rt := setupE2eTestEnv(t)

	_, err := rt.RunString(`
const k8s = new Kubernetes()

let pod = {
	apiVersion: "v1",
	kind:       "Pod",
	metadata: {
	    name:      "busybox",
	    namespace:  "default"
	},
	spec: {
	    containers: [
		{
		    name:    "busybox",
		    image:   "busybox",
		    command: ["sh", "-c", "sleep 30"]
		}
	    ]
	},
	status: {
		phase: "Running"
	}
}

// create pod in test namespace
k8s.create(pod)

// get helpers for test namespace
const helpers = k8s.helpers()

// wait for pod to be running
if (!helpers.waitPodRunning(pod.metadata.name, 5)) {
	throw new Error("should not timeout")
}
`)
	require.NoError(t, err)
}