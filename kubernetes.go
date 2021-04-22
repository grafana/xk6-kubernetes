package kubernetes

import (
	"github.com/loadimpact/k6/js/modules"
)

const version = "v0.0.1"

type Kubernetes struct {
	Version string
}

func init() {
	modules.Register("k6/x/kubernetes", &Kubernetes{
		Version: version,
	})
}
