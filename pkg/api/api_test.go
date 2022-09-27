package api

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
)

var podSpec = map[string]interface{}{
	"apiVersion": "v1",
	"kind":       "Pod",
	"metadata": map[string]interface{}{
		"name":      "busybox",
		"namespace": "testns",
	},
	"spec": map[string]interface{}{
		"containers": []interface{}{
			map[string]interface{}{
				"name":    "busybox",
				"image":   "busybox",
				"command": []interface{}{"sh", "-c", "sleep 30"},
			},
		},
	},
}

var JobSpec = map[string]interface{}{
	"apiVersion": "apps/v1",
	"kind":       "Job",
	"metadata": map[string]interface{}{
		"name":      "busybox",
		"namespace": "testns",
	},
	"spec": map[string]interface{}{
		"template": map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":    "busybox",
						"image":   "busybox",
						"command": []interface{}{"sh", "-c", "sleep 30"},
					},
				},
			},
		},
	},
}

func newForTest() *kubernetes {
	scheme := runtime.NewScheme()
	fake.AddToScheme(scheme)
	client := dynamicfake.NewSimpleDynamicClient(scheme)
	return &kubernetes{
		ctx:        context.TODO(),
		client:     client,
		serializer: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
	}
}

func TestCreateGetListDelete(t *testing.T) {
	testCases := []struct {
		test string
		obj  map[string]interface{}
		kind string
		name string
		ns   string
	}{
		{
			test: "Create Get List Delete Pods",
			obj:  podSpec,
			kind: "Pod",
			name: "busybox",
			ns:   "testns",
		},
		{
			test: "Create Get List Delete Jobs",
			obj:  JobSpec,
			kind: "Job",
			name: "busybox",
			ns:   "testns",
		},
	}

	for _, tc := range testCases {
		k := newForTest()
		_, err := k.Create(tc.obj)
		if err != nil {
			t.Errorf("failed %v", err)
			return
		}

		obj, err := k.Get(tc.kind, tc.name, tc.ns)
		if err != nil {
			t.Errorf("failed %v", err)
			return
		}
		if obj == nil {
			t.Errorf("invalid value returned")
			return
		}
		pods, err := k.List(tc.kind, tc.ns)
		if err != nil {
			t.Errorf("failed to get list of %ss: %v", tc.kind, err)
			return
		}

		if len(pods) == 0 {
			t.Errorf("expected one %s but none received", tc.kind)
			return
		}

		err = k.Delete(tc.kind, tc.name, tc.ns)
		if err != nil {
			t.Errorf("failed to delete pod: %v", err)
			return
		}
	}
}

var podManifest = `
apiVersion: v1
kind: Pod
metadata:
  name: busybox
  namespace: testns
spec:
  containers:
  - name: busybox
    image: busybox
    command: ["sleep", "300"]
`

func TestApply(t *testing.T) {
	testCases := []struct {
		test     string
		manifest string
		kind     string
		name     string
		ns       string
	}{
		{
			test:     "Apply pod manifest",
			manifest: podManifest,
			kind:     "Pod",
			name:     "busybox",
			ns:       "testns",
		},
	}

	for _, tc := range testCases {
		k := newForTest()
		err := k.Apply(tc.manifest)
		if err != nil {
			t.Errorf("failed %v", err)
			return
		}

		obj, err := k.Get(tc.kind, tc.name, tc.ns)
		if err != nil {
			t.Errorf("failed %v", err)
			return
		}
		if obj == nil {
			t.Errorf("invalid value returned")
			return
		}
	}
}
