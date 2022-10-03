package resources

import (
	"context"
	"testing"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"github.com/grafana/xk6-kubernetes/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func newForTest(objs ...runtime.Object) (*Client, error) {
	dynamic, err := testutils.NewFakeDynamic(objs...)
	if err != nil {
		return nil, err
	}
	return NewFromClient(context.TODO(), dynamic), nil
}

func TestCreateGetUpdateListDelete(t *testing.T) {
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

func TestUpdate(t *testing.T) {
	t.Parallel()

	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "busybox",
			Namespace: "testns",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sh", "-c", "sleep 30"},
				},
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
		},
	}
	c, err := newForTest(&pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	pod.Status.Phase = corev1.PodFailed
	podObj, err := utils.RuntimeToUnstructured(&pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}
	updated, err := c.Update(podObj.UnstructuredContent())
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	updatedPod := &corev1.Pod{}
	err = utils.GenericToRuntime(updated, updatedPod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	if updatedPod.Status.Phase != corev1.PodFailed {
		t.Errorf("pod phase was not updated")
	}
}
