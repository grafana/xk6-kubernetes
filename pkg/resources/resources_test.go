package resources

import (
	"context"
	"testing"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
)

func podSpec() map[string]interface{} {
	return map[string]interface{}{
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
}

func jobSpec() map[string]interface{} {
	return map[string]interface{}{
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
}

func newForTest() (*Client, error) {
	dynamic, err := testutils.NewFakeDynamic()
	if err != nil {
		return nil, err
	}
	return NewFromClient(context.TODO(), dynamic), nil
}

func TestCreateGetListDelete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		test string
		obj  map[string]interface{}
		kind string
		name string
		ns   string
	}{
		{
			test: "Create Get List Delete Pods",
			obj:  podSpec(),
			kind: "Pod",
			name: "busybox",
			ns:   "testns",
		},
		{
			test: "Create Get List Delete Jobs",
			obj:  jobSpec(),
			kind: "Job",
			name: "busybox",
			ns:   "testns",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()
			c, err := newForTest()
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}
			_, err = c.Create(tc.obj)
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}

			obj, err := c.Get(tc.kind, tc.name, tc.ns)
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}
			if obj == nil {
				t.Errorf("invalid value returned")
				return
			}
			pods, err := c.List(tc.kind, tc.ns)
			if err != nil {
				t.Errorf("failed to get list of %ss: %v", tc.kind, err)
				return
			}

			if len(pods) == 0 {
				t.Errorf("expected one %s but none received", tc.kind)
				return
			}

			err = c.Delete(tc.kind, tc.name, tc.ns)
			if err != nil {
				t.Errorf("failed to delete pod: %v", err)
				return
			}
		})
	}
}

func podManifest() string {
	return `
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
}

func TestApply(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		test     string
		manifest string
		kind     string
		name     string
		ns       string
	}{
		{
			test:     "Apply pod manifest",
			manifest: podManifest(),
			kind:     "Pod",
			name:     "busybox",
			ns:       "testns",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()
			c, err := newForTest()
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}
			err = c.Apply(tc.manifest)
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}

			obj, err := c.Get(tc.kind, tc.name, tc.ns)
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}
			if obj == nil {
				t.Errorf("invalid value returned")
				return
			}
		})
	}
}
