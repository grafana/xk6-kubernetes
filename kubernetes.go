package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"path/filepath"

	"github.com/grafana/xk6-kubernetes/pkg/configmaps"
	"github.com/grafana/xk6-kubernetes/pkg/deployments"
	"github.com/grafana/xk6-kubernetes/pkg/ingresses"
	"github.com/grafana/xk6-kubernetes/pkg/jobs"
	"github.com/grafana/xk6-kubernetes/pkg/namespaces"
	"github.com/grafana/xk6-kubernetes/pkg/nodes"
	"github.com/grafana/xk6-kubernetes/pkg/persistentvolumeclaims"
	"github.com/grafana/xk6-kubernetes/pkg/persistentvolumes"
	"github.com/grafana/xk6-kubernetes/pkg/pods"
	"github.com/grafana/xk6-kubernetes/pkg/secrets"
	"github.com/grafana/xk6-kubernetes/pkg/services"

	"go.k6.io/k6/js/modules"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func init() {
	modules.Register("k6/x/kubernetes", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/kubernetes` module instances for each VU.
type RootModule struct{}

// ModuleInstance represents an instance of the JS module.
type ModuleInstance struct {
	vu modules.VU
}

// Kubernetes is the exported object used within JavaScript.
type Kubernetes struct {
	client                 *kubernetes.Clientset
	metaOptions            metav1.ListOptions
	ctx                    context.Context
	ConfigMaps             *configmaps.ConfigMaps
	Ingresses              *ingresses.Ingresses
	Deployments            *deployments.Deployments
	Pods                   *pods.Pods
	Namespaces             *namespaces.Namespaces
	Nodes                  *nodes.Nodes
	Jobs                   *jobs.Jobs
	Services               *services.Services
	Secrets                *secrets.Secrets
	PersistentVolumes      *persistentvolumes.PersistentVolumes
	PersistentVolumeClaims *persistentvolumeclaims.PersistentVolumeClaims
}

// KubeConfig represents the initialization settings for the kubernetes api client.
type KubeConfig struct {
	ConfigPath string
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu: vu,
	}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"Kubernetes": mi.newClient,
		},
	}
}

func (mi *ModuleInstance) newClient(c goja.ConstructorCall) *goja.Object {
	rt := mi.vu.Runtime()
	ctx := mi.vu.Context()

	var options KubeConfig
	err := rt.ExportTo(c.Argument(0), &options)
	if err != nil {
		common.Throw(rt, fmt.Errorf("Kubernetes constructor expect KubernetesOption as it's argument: %w", err))
	}

	kubeconfig := options.ConfigPath
	if kubeconfig == "" {
		home := homedir.HomeDir()
		if home == "" {
			common.Throw(rt, errors.New("Home dir not found"))
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		common.Throw(rt, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.Throw(rt, err)
	}

	obj := &Kubernetes{}
	obj.client = clientset
	obj.metaOptions = metav1.ListOptions{}
	obj.ctx = ctx

	obj.ConfigMaps = configmaps.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Ingresses = ingresses.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Deployments = deployments.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Pods = pods.New(obj.client, config, obj.metaOptions, obj.ctx)
	obj.Namespaces = namespaces.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Nodes = nodes.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Jobs = jobs.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Services = services.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Secrets = secrets.New(obj.client, obj.metaOptions, obj.ctx)
	obj.PersistentVolumes = persistentvolumes.New(obj.client, obj.metaOptions, obj.ctx)
	obj.PersistentVolumeClaims = persistentvolumeclaims.New(obj.client, obj.metaOptions, obj.ctx)

	return rt.ToValue(obj).ToObject(rt)
}
