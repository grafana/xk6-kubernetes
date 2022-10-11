package resources

import (
	"context"
	"testing"

	"github.com/grafana/xk6-kubernetes/internal/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func buildUnstructuredPod() map[string]interface{} {
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

func buildUnstructuredJob() map[string]interface{} {
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

func buildPod() corev1.Pod {
	return corev1.Pod{
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
}

func newForTest(objs ...runtime.Object) (*Client, error) {
	dynamic, err := testutils.NewFakeDynamic(objs...)
	if err != nil {
		return nil, err
	}
	return NewFromClient(context.TODO(), dynamic), nil
}

func TestCreate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		test     string
		obj      map[string]interface{}
		kind     string
		resource schema.GroupVersionResource
		name     string
		ns       string
	}{
		{
			test:     "Create Get List Delete Pods",
			obj:      buildUnstructuredPod(),
			kind:     "Pod",
			resource: schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			name:     "busybox",
			ns:       "testns",
		},
		{
			test:     "Create Get List Delete Jobs",
			obj:      buildUnstructuredJob(),
			kind:     "Job",
			resource: schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"},
			name:     "busybox",
			ns:       "testns",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()
			fake, _ := testutils.NewFakeDynamic()
			c := NewFromClient(context.TODO(), fake)

			created, err := c.Create(tc.obj)
			if err != nil {
				t.Errorf("failed %v", err)
				return
			}

			name, found, err := unstructured.NestedString(created, "metadata", "name")
			if err != nil {
				t.Errorf("error retrieving pod name %v", err)
				return
			}

			if !found {
				t.Errorf("object created has no name field")
				return
			}

			if name != tc.name {
				t.Errorf("wrong object retrieved. Expected %s Received %s", tc.name, name)
				return
			}

			// check the object was added to the fake client's object tracker
			_, err = fake.Tracker().Get(tc.resource, tc.ns, tc.name)
			if err != nil {
				t.Errorf("error retrieving object %v", err)
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

	// initialize with pod
	obj := buildPod()
	c, err := newForTest(&obj)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	// set the status
	pod := buildUnstructuredPod()
	pod["status"] = map[string]interface{}{
		"phase": string(corev1.PodFailed),
	}

	updated, err := c.Update(pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	// get the status
	status, found, err := unstructured.NestedString(updated, "status", "phase")
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}
	if !found || status != string(corev1.PodFailed) {
		t.Errorf("pod phase was not updated")
	}
}

func TestStructuredCreate(t *testing.T) {
	t.Parallel()

	fake, _ := testutils.NewFakeDynamic()
	c := NewFromClient(context.TODO(), fake)

	pod := buildPod()
	created, err := c.Structured().Create(pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	if created.(corev1.Pod).Name != pod.Name {
		t.Errorf("invalid pod returned. Expected %s Returned %s", pod.Name, created.(corev1.Pod).Name)
		return
	}
}

func TestStructuredGet(t *testing.T) {
	t.Parallel()
	// initialize with pod
	initPod := buildPod()
	c, err := newForTest(&initPod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	pod := &corev1.Pod{}
	err = c.Structured().Get("Pod", "busybox", "testns", pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}
	if pod.Name != initPod.Name {
		t.Errorf("invalid pod returned. Expected %s Returned %s", initPod.Name, pod.Name)
		return
	}
}

func TestStructuredList(t *testing.T) {
	t.Parallel()
	// initialize with pod
	pod := buildPod()
	c, err := newForTest(&pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	podList := []corev1.Pod{}
	err = c.Structured().List("Pod", "testns", &podList)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	if len(podList) != 1 {
		t.Errorf("one pod expected but %d returned", len(podList))
		return
	}

	if podList[0].Name != pod.Name {
		t.Errorf("invalid pod returned. Expected %s Returned %s", pod.Name, podList[0].Name)
		return
	}
}

func TestStructuredDelete(t *testing.T) {
	t.Parallel()
	// initialize with pod
	obj := buildPod()
	c, err := newForTest(&obj)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	err = c.Structured().Delete("Pod", "busybox", "testns")
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}
}

func TestStructuredUpdate(t *testing.T) {
	t.Parallel()
	// initialize with pod
	pod := buildPod()
	c, err := newForTest(&pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	// change status
	pod.Status.Phase = corev1.PodFailed
	updated, err := c.Structured().Update(pod)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}

	status := updated.(corev1.Pod).Status.Phase
	if status != corev1.PodFailed {
		t.Errorf("pod status not updated")
		return
	}
}
