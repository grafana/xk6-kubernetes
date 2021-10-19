package kubernetes

import (
	"context"
	"errors"
	"github.com/k6io/xk6-kubernetes/pkg/configmaps"
	"github.com/k6io/xk6-kubernetes/pkg/jobs"
	"github.com/k6io/xk6-kubernetes/pkg/pods"
	"path/filepath"

	"go.k6.io/k6/js/modules"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const version = "v0.0.1"

type Kubernetes struct {
	Version     string
	client      *kubernetes.Clientset
	metaOptions metav1.ListOptions
	ctx         context.Context
	ConfigMaps  *configmaps.ConfigMaps
	Pods        *pods.Pods
	Jobs        *jobs.Jobs
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
	obj.Pods = pods.New(obj.client, obj.metaOptions, obj.ctx)
	obj.Jobs = jobs.New(obj.client, obj.metaOptions, obj.ctx)

	return obj, nil
}

func init() {
	k8s := &Kubernetes{
		Version: version,
	}
	modules.Register("k6/x/kubernetes", k8s)
}
