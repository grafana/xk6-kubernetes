package jobs

import (
	"strings"
	"testing"
	"time"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
	k8sTypes "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
)

var (
	testName      = "job-test"
	testNamespace = "ns-test"
)

func TestJobs_Apply(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewJob("existing", testNamespace),
	), metav1.ListOptions{}, nil)

	testCases := []struct {
		testID        string
		yaml          string
		namespace     string
		expectedError string
	}{
		{
			testID: "apply with invalid yaml",
			yaml: `
this is not yaml
`,
			namespace:     testNamespace,
			expectedError: "json parse error",
		},
		{
			testID: "apply with incorrect resource yaml",
			yaml: `
apiVersion: v1
kind: Namespace
metadata:
  name: ` + testName + `
`,
			namespace:     testNamespace,
			expectedError: "Yaml was not a Job",
		},
		{
			testID: "create new job from yaml",
			yaml: `
apiVersion: batch/v1
kind: Job
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  template:
    spec:
      containers:
      - name: ` + testName + `
        image: perl
        command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
`,
			namespace:     testNamespace,
			expectedError: "",
		},
		{
			testID: "update existing job from yaml",
			yaml: `
apiVersion: batch/v1
kind: Job
metadata:
  name: ` + testName + `
  annotations:
  labels:
    app: xk6-kubernetes/unit-test
spec:
  template:
    spec:
      containers:
      - name: ` + testName + `
        image: perl
        command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
`,
			namespace:     testNamespace,
			expectedError: "jobs.batch \"" + testName + "\" already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.Apply(tc.yaml, tc.namespace)

			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Errorf("encountered error: %v, expected: %v", err, tc.expectedError)
				return
			}
			if err == nil && tc.expectedError != "" {
				t.Errorf("expected error \"%v\", but got none", tc.expectedError)
				return
			}
			if err == nil && (result.Name != testName || result.Namespace != testNamespace) {
				t.Errorf("incorrect instance was returned")
				return
			}
		})
	}
}

func TestJobs_Create(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewJob("existing", testNamespace),
	), metav1.ListOptions{}, nil)

	options := JobOptions{
		Name:          testName,
		Namespace:     testNamespace,
		NodeName:      "node-1",
		Image:         "perl",
		Command:       []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
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
	jobs, _ := fixture.List(testNamespace)
	if len(jobs) != 2 {
		t.Errorf("expecting 2 jobs in namespace, listing returned %v", len(jobs))
		return
	}
}

func TestJobs_List(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewJob("job-1", testNamespace),
		testutils.NewJob("job-2", testNamespace),
		testutils.NewJob("job-3", testNamespace),
	), metav1.ListOptions{}, nil)

	testCases := []struct {
		testID        string
		namespace     string
		expectedCount int
	}{
		{
			testID:        "test namespace returns 3 jobs",
			namespace:     testNamespace,
			expectedCount: 3,
		},
		{
			testID:        "empty namespace returns 0 jobs",
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
				t.Errorf("received %v secret(s), expected %v", len(result), tc.expectedCount)
				return
			}
			for _, job := range result {
				if tc.namespace != job.Namespace {
					t.Errorf("received job from %v namespace, only expected %v", job.Namespace, tc.namespace)
					return
				}
			}
		})
	}
}

func TestJobs_Delete(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewJob(testName, testNamespace),
	), metav1.ListOptions{}, nil)

	err := fixture.Delete(testName, testNamespace)

	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	jobs, _ := fixture.List(testNamespace)
	if len(jobs) != 0 {
		t.Errorf("expecting 0 jobs in namespace, listing returned %v", len(jobs))
		return
	}
}

func TestJobs_Get(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewJob(testName, testNamespace),
		testutils.NewJob("job-other", "ns-2"),
	), metav1.ListOptions{}, nil)

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
			name:         "job-other",
			namespace:    testNamespace,
			expectToFind: false,
		},
		{
			testID:       "fetching unknown name from any namespace returns nothing",
			name:         "job-unknown",
			namespace:    "any-namespace",
			expectToFind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			result, err := fixture.Get(tc.name, tc.namespace)
			if err != nil && !errors.IsNotFound(err) {
				t.Errorf("encountered an error: %v", err)
				return
			}
			if tc.expectToFind && err != nil {
				t.Errorf("expected to find job %v in %v namespace, but received error: %v", tc.name, tc.namespace, err)
				return
			}
			if !tc.expectToFind && err == nil {
				t.Errorf("expected an error when trying to find unavailable job %v in %v namespace", tc.name, tc.namespace)
				return
			}
			if tc.expectToFind && result.Name != tc.name {
				t.Errorf("received job with name %v, expected %v", result.Name, tc.name)
				return
			}
			if tc.expectToFind && result.Namespace != tc.namespace {
				t.Errorf("received job with namespace %v, expected %v", result.Name, tc.name)
				return
			}
		})
	}
}

func TestJobs_Wait(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		test           string
		name           string
		namespace      string
		status         string
		delay          time.Duration
		expectError    bool
		expectedResult bool
		timeout        string
	}

	testCases := []TestCase{
		{
			test:           "wait job complete",
			name:           "pod-running",
			namespace:      testNamespace,
			status:         "Complete",
			delay:          1 * time.Second,
			expectError:    false,
			expectedResult: true,
			timeout:        "5s",
		},
		{
			test:           "timeout waiting job complete",
			name:           "pod-running",
			namespace:      testNamespace,
			status:         "Complete",
			delay:          10 * time.Second,
			expectError:    false,
			expectedResult: false,
			timeout:        "5s",
		},
		{
			test:           "wait job failed",
			name:           "pod-running",
			namespace:      testNamespace,
			status:         "Failed",
			delay:          1 * time.Second,
			expectError:    true,
			expectedResult: false,
			timeout:        "5s",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			client := fake.NewSimpleClientset()
			watcher := watch.NewFake()
			client.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))
			fixture := New(client, metav1.ListOptions{}, nil)
			go func(tc TestCase) {
				time.Sleep(tc.delay)
				watcher.Modify(testutils.NewJobWithStatus(tc.name, tc.namespace, tc.status))
			}(tc)

			result, err := fixture.Wait(WaitOptions{
				Name:      tc.name,
				Namespace: tc.namespace,
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
