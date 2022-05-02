package nodes

import (
	"github.com/grafana/xk6-kubernetes/internal/testutils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

var (
	testName = "node-test"
)

func TestNodes_List(t *testing.T) {
	t.Parallel()
	fixture := New(fake.NewSimpleClientset(
		testutils.NewNodes("node-1"),
		testutils.NewNodes("node-2"),
		testutils.NewNodes("node-3"),
	), metav1.ListOptions{}, nil)

	result, err := fixture.List()
	if err != nil {
		t.Errorf("encountered an error: %v", err)
		return
	}
	if len(result) != 3 {
		t.Errorf("received %v node(s), expected %v", len(result), 3)
		return
	}
}
