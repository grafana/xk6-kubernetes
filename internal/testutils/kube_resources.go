package testutils

import (
	appsV1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	networkV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewConfigMap creates a basic ConfigMap instance having the provided namespace and name
func NewConfigMap(name string, namespace string) *coreV1.ConfigMap {
	return &coreV1.ConfigMap{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Data: map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	}
}

// NewDeployment is a helper to build a new Deployment instance
func NewDeployment(name string, namespace string) *appsV1.Deployment {
	return &appsV1.Deployment{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: nil,
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "xk6-kubernetes/unit-test",
				},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
							Ports: []coreV1.ContainerPort{
								{ContainerPort: 80},
							},
						},
					},
				},
			},
		},
	}
}

// NewIngress is a helper to build a new Ingress instance
func NewIngress(name string, namespace string) *networkV1.Ingress {
	return &networkV1.Ingress{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: networkV1.IngressSpec{
			IngressClassName: nil,
			Rules:            []networkV1.IngressRule{},
		},
	}
}

// NewJob is a helper to build a new Job instance
func NewJob(name string, namespace string) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: nil,
			Template:     coreV1.PodTemplateSpec{},
		},
		Status: batchv1.JobStatus{
			Conditions: []batchv1.JobCondition{},
		},
	}
}

// NewJobWithStatus creates a job with a given condition
func NewJobWithStatus(name string, namespace string, status string) *batchv1.Job {
	job := NewJob(name, namespace)
	job.Status.Conditions = []batchv1.JobCondition{
		{
			Type:   batchv1.JobConditionType(status),
			Status: coreV1.ConditionTrue,
		},
	}
	return job
}

// NewNamespace is a helper to build a new Namespace instance
func NewNamespace(name string) *coreV1.Namespace {
	return &coreV1.Namespace{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
	}
}

// NewNodes is a helper to build a new Node instance
func NewNodes(name string) *coreV1.Node {
	return &coreV1.Node{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Node",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: coreV1.NodeSpec{},
	}
}

// NewPersistentVolumeClaim is a helper to build a new PersistentVolumeClaim instance
func NewPersistentVolumeClaim(name string, namespace string) *coreV1.PersistentVolumeClaim {
	return &coreV1.PersistentVolumeClaim{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			StorageClassName: nil,
			AccessModes: []coreV1.PersistentVolumeAccessMode{
				coreV1.ReadWriteMany,
			},
			Resources: coreV1.ResourceRequirements{},
		},
	}
}

// NewPersistentVolume is a helper to build a new PersistentVolume instance
func NewPersistentVolume(name string) *coreV1.PersistentVolume {
	return &coreV1.PersistentVolume{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolume",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: coreV1.PersistentVolumeSpec{
			Capacity: map[coreV1.ResourceName]resource.Quantity{
				coreV1.ResourceStorage: {
					Format: "10G",
				},
			},
			PersistentVolumeSource: coreV1.PersistentVolumeSource{
				HostPath: &coreV1.HostPathVolumeSource{
					Path: "/tmp/" + name,
				},
			},
			StorageClassName: "local-storage",
		},
	}
}

// NewPod is a helper to build a new Pod instance
func NewPod(name string, namespace string) *coreV1.Pod {
	return &coreV1.Pod{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sh", "-c", "sleep 300"},
				},
			},
			EphemeralContainers: nil,
		},
		Status: coreV1.PodStatus{
			Phase: coreV1.PodRunning,
		},
	}
}

// NewPodWithStatus is a helper for building Pods with a given Status
func NewPodWithStatus(name string, namespace string, phase string) *coreV1.Pod {
	pod := NewPod(name, namespace)
	pod.Status.Phase = coreV1.PodPhase(phase)
	return pod
}

// NewSecret is a helper to build a new Secret instance
func NewSecret(name string, namespace string) *coreV1.Secret {
	return &coreV1.Secret{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		StringData: map[string]string{
			"secret-key-1": "secret-value-1",
			"secret-key-2": "secret-value-2",
		},
		Type: "Opaque",
	}
}

// NewService is a helper to build a new Service instance
func NewService(name string, namespace string) *coreV1.Service {
	return &coreV1.Service{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "xk6-kubernetes/unit-test",
			},
		},
		Spec: coreV1.ServiceSpec{
			Selector: map[string]string{
				"app": "MyApp",
			},
			Ports: []coreV1.ServicePort{
				{Port: 80, Protocol: coreV1.ProtocolTCP},
			},
		},
	}
}
