package nodes

import (
	"context"
	"testing"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodes_List(t *testing.T) {
	t.Parallel()
	fixture := New(context.Background(), fake.NewSimpleClientset(
		testutils.NewNodes("node-1"),
		testutils.NewNodes("node-2"),
		testutils.NewNodes("node-3"),
	), metav1.ListOptions{})

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
