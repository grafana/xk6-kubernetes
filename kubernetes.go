package kubernetes

import (
	"context"
	"errors"
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

const version = "v0.0.1"

type Kubernetes struct {
	Version                string
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

type KubernetesOptions struct {
	ConfigPath string
}

func (obj *Kubernetes) XKubernetes(ctx *context.Context, options KubernetesOptions) (*Kubernetes, error) {
	kubeconfig := options.ConfigPath
	if kubeconfig == "" {
		home := homedir.HomeDir()
		if home == "" {
			return nil, errors.New("Home dir not found")
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	obj.client = clientset
	obj.metaOptions = metav1.ListOptions{}
	obj.ctx = *ctx

	obj.ConfigMaps = configmaps.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Ingresses = ingresses.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Deployments = deployments.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Pods = pods.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Namespaces = namespaces.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Nodes = nodes.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Jobs = jobs.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Services = services.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Secrets = secrets.New(obj.client, obj.metaOptions, obj.ctx)
	obj.PersistentVolumes = persistentvolumes.New(obj.client, obj.metaOptions, obj.ctx)
	obj.PersistentVolumeClaims = persistentvolumeclaims.New(obj.client, obj.metaOptions, obj.ctx)

	return obj, nil
}

func init() {
	k8s := &Kubernetes{
		Version: version,
	}
	modules.Register("k6/x/kubernetes", k8s)
}
