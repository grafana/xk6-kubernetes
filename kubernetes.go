package kubernetes

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/loadimpact/k6/js/modules"
	k8sTypes "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const version = "v0.0.1"

type Kubernetes struct {
	Version     string
	Client      *kubernetes.Clientset
	meraOptions metav1.ListOptions
	ctx         context.Context
}

type KubernetesOptions struct {
	ConfigPath string
}

func (obj *Kubernetes) Init(options KubernetesOptions) error {
	kubeconfig := options.ConfigPath
	if kubeconfig == "" {
		home := homedir.HomeDir()
		if home == "" {
			return errors.New("Home dir not found")
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	obj.Client = clientset
	obj.meraOptions = metav1.ListOptions{}
	obj.ctx = context.Background()
	return nil
}

func (obj *Kubernetes) GetPods() ([]k8sTypes.Pod, error) {
	pods, err := obj.Client.CoreV1().Pods("").List(obj.ctx, obj.meraOptions)
	if err != nil {
		return []k8sTypes.Pod{}, err
	}
	return pods.Items, nil
}

func init() {
	k8s := &Kubernetes{
		Version: version,
	}
	modules.Register("k6/x/kubernetes", k8s)
}
