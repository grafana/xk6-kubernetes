package helpers

import (
	"context"
	"github.com/grafana/xk6-kubernetes/pkg/resources"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake" //nolint:typecheck
	k8stest "k8s.io/client-go/testing"
	"testing"
	"time"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
)

const (
	jobName = "test-job"
)

func TestWaitJobCompleted(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		test           string
		status         string
		delay          time.Duration
		expectError    bool
		expectedResult bool
		timeout        uint
	}

	testCases := []TestCase{
		{
			test:           "job completed before timeout",
			status:         "Complete",
			delay:          1 * time.Second,
			expectError:    false,
			expectedResult: true,
			timeout:        60,
		},
		{
			test:           "timeout waiting for job to complete",
			status:         "Complete",
			delay:          10 * time.Second,
			expectError:    false,
			expectedResult: false,
			timeout:        5,
		},
		{
			test:           "job failed before timeout",
			status:         "Failed",
			delay:          1 * time.Second,
			expectError:    true,
			expectedResult: false,
			timeout:        60,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()

			clientset := testutils.NewFakeClientset().(*fake.Clientset)
			watcher := watch.NewRaceFreeFake()
			clientset.PrependWatchReactor("jobs", k8stest.DefaultWatchReactor(watcher, nil))

			fake, _ := testutils.NewFakeDynamic()
			client := resources.NewFromClient(context.TODO(), fake).WithMapper(&testutils.FakeRESTMapper{})

			fixture := NewHelper(context.TODO(), clientset, client, nil, "default")
			job := testutils.NewJob(jobName, "default")
			_, err := client.Structured().Create(job)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			go func(tc TestCase) {
				time.Sleep(tc.delay)
				job = testutils.NewJobWithStatus(jobName, "default", tc.status)
				_, e := client.Structured().Update(job)
				if e != nil {
					t.Errorf("unexpected error: %v", e)
					return
				}
				watcher.Modify(job)
			}(tc)

			result, err := fixture.WaitJobCompleted(
				jobName,
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
