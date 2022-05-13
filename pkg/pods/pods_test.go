package pods

import (
	"github.com/grafana/xk6-kubernetes/internal/testutils"
	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
	"strings"
	"testing"
	"time"
)

var (
	testName      = "pod-test"
	testNamespace = "ns-test"
)

func TestPods_Create(t *testing.T) {
	t.Parallel()
	// TODO Figure out the rest.Config
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPod("existing", testNamespace),
	), nil, metav1.ListOptions{}, nil)

	options := PodOptions{
		Name:          testName,
		Namespace:     testNamespace,
		Image:         "busybox",
		Command:       []string{"sh", "-c", "sleep 300"},
		RestartPolicy: k8sTypes.RestartPolicyNever,
	}
	result, err := fixture.Create(options)

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if result.Name != options.Name || result.Namespace != options.Namespace {
		t.Errorf("incorrect instance was returned")
		return
	}
	pods, _ := fixture.List(testNamespace)
	if len(pods) != 2 {
		t.Errorf("expecting 2 pods in namespace, listing returned %v", len(pods))
		return
	}
}

func TestPods_Wait(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		test           string
		name           string
		namespace      string
		status         string
		delay          time.Duration
		expectedStatus string
		expectError    bool
		expectedResult bool
		timeout        string
	}

	testCases := []TestCase{
		{
			test:           "wait pod running",
			name:           "pod-running",
			namespace:      testNamespace,
			status:         "Running",
			delay:          1 * time.Second,
			expectedStatus: "Running",
			expectError:    false,
			expectedResult: true,
			timeout:        "5s",
		},
		{
			test:           "timeout waiting pod running",
			name:           "pod-running",
			namespace:      testNamespace,
			status:         "Running",
			delay:          10 * time.Second,
			expectedStatus: "Running",
			expectError:    false,
			expectedResult: false,
			timeout:        "5s",
		},
		{
			test:           "wait failed pod",
			name:           "pod-running",
			namespace:      testNamespace,
			status:         "Failed",
			delay:          1 * time.Second,
			expectedStatus: "Running",
			expectError:    true,
			expectedResult: false,
			timeout:        "5s",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			// TODO Figure out the rest.Config
			client := fake.NewSimpleClientset()
			watcher := watch.NewFake()
			client.PrependWatchReactor("pods", k8stest.DefaultWatchReactor(watcher, nil))
			fixture := New(client, nil, metav1.ListOptions{}, nil)
			go func(tc TestCase) {
				time.Sleep(tc.delay)
				watcher.Modify(testutils.NewPodWithStatus(tc.name, tc.namespace, tc.status))
			}(tc)

			result, err := fixture.Wait(WaitOptions{
				Name:      tc.name,
				Namespace: tc.namespace,
				Status:    tc.expectedStatus,
				Timeout:   tc.timeout,
			})

			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but none returned")
				return
			}
			if result != tc.expectedResult {
				t.Errorf("expected result %t but %t returned", tc.expectedResult, result)
				return
			}
		})
	}
}

func TestPods_List(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPod("pod-1", testNamespace),
		testutils.NewPod("pod-2", testNamespace),
		testutils.NewPod("pod-3", testNamespace),
	), nil, metav1.ListOptions{}, nil)

	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 pods",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 pods",
			namespace:     "ns-empty",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.List(tc.namespace)
			if err != nil {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if len(result) != tc.expectedCount {
				t.Errorf("received %v pod(s), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, pod := range result {
				if tc.namespace != pod.Namespace {
					t.Errorf("received pod from %v namespace, only expected %v", pod.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestPods_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPod(testName, testNamespace),
	), nil, metav1.ListOptions{}, nil)

	err := fixture.Delete(testName, testNamespace, metav1.DeleteOptions{})

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	pods, _ := fixture.List(testNamespace)
	if len(pods) != 0 {
		t.Errorf("expecting 0 pods in namespace, listing returned %v", len(pods))
		return
	}
}

func TestPods_Get(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewPod(testName, testNamespace),
		testutils.NewPod("pod-other", "ns-2"),
	), nil, metav1.ListOptions{}, nil)

	testCases := []struct {
		testID       string
		name         string
		namespace    string
		expectToFind bool
	}{
		{
			testID:       "fetching valid name within namespace returns correctly",
			name:         testName,
			namespace:    testNamespace,
			expectToFind: true,
		},
		{
			testID:       "fetching valid name from incorrect namespace returns nothing",
			name:         "pod-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "pod-unknown",
			namespace:    "any-namespace",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.Get(tc.name, tc.namespace)
			if err != nil && !strings.Contains(err.Error(), "pod not found") {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find pod %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable pod %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received pod with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received pod with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}
