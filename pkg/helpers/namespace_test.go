package helpers

import (
	"context"
	"testing"

	"github.com/grafana/xk6-kubernetes/internal/testutils"
	"github.com/grafana/xk6-kubernetes/pkg/resources"
)

func Test_CreateRandomNs(t *testing.T) {
	t.Parallel()

	fake, _ := testutils.NewFakeDynamic()
	client := resources.NewFromClient(context.TODO(), fake)
	h := NewHelper(context.TODO(), client, "default")

	prefix := "test-"
	// Ignore returned value as the fake client will not assign a name to a resource
	// created with a GenerateName parameter.
	_, err := h.CreateRandomNamespace(prefix)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}
}
