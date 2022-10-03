package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"github.com/grafana/xk6-kubernetes/pkg/resources"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	k8stest "k8s.io/client-go/testing"
)

const (
	testNamespace = "ns-test"
)

func buildPod(name, namespace string, status corev1.PodPhase) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sh", "-c", "sleep 300"},
				},
			},
		},
		Status: corev1.PodStatus{
			Phase: status,
		},
	}
}

func TestPods_Wait(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		test           string
		name           string
		status         corev1.PodPhase
		delay          time.Duration
		expectError    bool
		expectedResult bool
		timeout        uint
	}

	testCases := []TestCase{
		{
			test:           "wait pod running",
			name:           "pod-running",
			delay:          1 * time.Second,
			status:         corev1.PodRunning,
			expectError:    false,
			expectedResult: true,
			timeout:        5,
		},
		{
			test:           "timeout waiting pod running",
			name:           "pod-running",
			status:         corev1.PodRunning,
			delay:          10 * time.Second,
			expectError:    false,
			expectedResult: false,
			timeout:        5,
		},
		{
			test:           "wait failed pod",
			name:           "pod-running",
			status:         corev1.PodFailed,
			delay:          1 * time.Second,
			expectError:    true,
			expectedResult: false,
			timeout:        5,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()
			fake, _ := testutils.NewFakeDynamic()
			watcher := watch.NewRaceFreeFake()
			fake.PrependWatchReactor("pods", k8stest.DefaultWatchReactor(watcher, nil))
			client := resources.NewFromClient(context.TODO(), fake)
			h := NewHelper(context.TODO(), client, testNamespace)
			go func(tc TestCase) {
				pod := buildPod(tc.name, testNamespace, tc.status)
				time.Sleep(tc.delay)
				watcher.Modify(pod)
			}(tc)

			result, err := h.WaitPodRunning(
				tc.name,
				tc.timeout,
			)

			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tc.expectError && err == nil {
				t.Error("expected an error but none returned")
				return
			}
			if result != tc.expectedResult {
				t.Errorf("expected result %t but %t returned", tc.expectedResult, result)
				return
			}
		})
	}
}
