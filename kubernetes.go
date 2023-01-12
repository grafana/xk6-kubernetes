// Package kubernetes provides the xk6 Modules implementation for working with Kubernetes resources using Javascript
package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"k8s.io/client-go/rest"

	"github.com/grafana/xk6-kubernetes/pkg/api"
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
	"k8s.io/apimachinery/pkg/api/meta"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Required for access to GKE and AKS
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
	// clientset enables injection of a pre-configured Kubernetes environment for unit tests
	clientset kubernetes.Interface
	// dynamic enables injection of a fake dynamic client for unit tests
	dynamic dynamic.Interface
	// mapper enables injection of a fake RESTMapper for unit tests
	mapper meta.RESTMapper
}

// Kubernetes is the exported object used within JavaScript.
type Kubernetes struct {
	api.Kubernetes
	client                 kubernetes.Interface
	metaOptions            metaV1.ListOptions
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

	obj := &Kubernetes{}
	var config *rest.Config

	// if clientset was not injected for unit testing
	if mi.clientset == nil {
		var options KubeConfig
		err := rt.ExportTo(c.Argument(0), &options)
		if err != nil {
			common.Throw(rt,
				fmt.Errorf("Kubernetes constructor expects KubeConfig as it's argument: %w", err))
		}
		config, err = getClientConfig(options)
		if err != nil {
			common.Throw(rt, err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			common.Throw(rt, err)
		}
		obj.client = clientset
	} else {
		// Pre-configured clientset is being injected for unit testing
		obj.client = mi.clientset
	}

	// If dynamic client was not injected for unit testing
	// It is assumed rest config is set
	if mi.dynamic == nil {
		k8s, err := api.NewFromConfig(
			api.KubernetesConfig{
				Config:  config,
				Context: ctx,
			},
		)
		if err != nil {
			common.Throw(rt, err)
		}
		obj.Kubernetes = k8s
	} else {
		// Pre-configured dynamic client and RESTMapper are injected for unit testing
		k8s, err := api.NewFromConfig(
			api.KubernetesConfig{
				Client:  mi.dynamic,
				Mapper:  mi.mapper,
				Context: ctx,
			},
		)
		if err != nil {
			common.Throw(rt, err)
		}
		obj.Kubernetes = k8s
	}

	obj.metaOptions = metaV1.ListOptions{}
	obj.ctx = ctx

	obj.ConfigMaps = configmaps.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Ingresses = ingresses.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Deployments = deployments.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Pods = pods.New(obj.ctx, obj.client, config, obj.metaOptions)
	obj.Namespaces = namespaces.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Nodes = nodes.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Jobs = jobs.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Services = services.New(obj.ctx, obj.client, obj.metaOptions)
	obj.Secrets = secrets.New(obj.ctx, obj.client, obj.metaOptions)
	obj.PersistentVolumes = persistentvolumes.New(obj.ctx, obj.client, obj.metaOptions)
	obj.PersistentVolumeClaims = persistentvolumeclaims.New(obj.ctx, obj.client, obj.metaOptions)

	return rt.ToValue(obj).ToObject(rt)
}

func getClientConfig(options KubeConfig) (*rest.Config, error) {
	kubeconfig := options.ConfigPath
	if kubeconfig == "" {
		home := homedir.HomeDir()
		if home == "" {
			return nil, errors.New("home directory not found")
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
