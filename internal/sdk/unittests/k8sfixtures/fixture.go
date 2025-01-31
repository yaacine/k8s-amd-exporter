package k8sfixtures

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1alpha1"
)

func NewPodResourcesFixture(t *testing.T) []*podresourcesapi.PodResources {
	t.Helper()

	return []*podresourcesapi.PodResources{
		{
			Name:      "pod-1",
			Namespace: "team-a",
			Containers: []*podresourcesapi.ContainerResources{
				{
					Name: "container-1",
					Devices: []*podresourcesapi.ContainerDevices{
						{
							ResourceName: "amd.com/gpu",
							DeviceIds: []string{
								"0001:34:00.0",
							},
						},
					},
				},
			},
		},
		{
			Name:      "pod-b",
			Namespace: "team-2",
			Containers: []*podresourcesapi.ContainerResources{
				{
					Name: "container-1",
					Devices: []*podresourcesapi.ContainerDevices{
						{
							ResourceName: "amd-custom-resource-name",
							DeviceIds: []string{
								"80EE/vgpu-2345",
							},
						},
					},
				},
			},
		},
		{
			Name:      "pod-c",
			Namespace: "team-2",
			Containers: []*podresourcesapi.ContainerResources{
				{
					Name: "container-1",
					Devices: []*podresourcesapi.ContainerDevices{
						{
							ResourceName: "amd.com/gpu",
							DeviceIds: []string{
								"0000:8e:00.0",
							},
						},
					},
				},
			},
		},
		{
			Name:      "pod-i",
			Namespace: "team-b",
			Containers: []*podresourcesapi.ContainerResources{
				{
					Name: "container-1",
					Devices: []*podresourcesapi.ContainerDevices{
						{
							ResourceName: "amd.com/gpu",
							DeviceIds: []string{
								"amd123/gi456",
							},
						},
					},
				},
			},
		},
		{
			Name:      "pod-ii",
			Namespace: "team-b",
			Containers: []*podresourcesapi.ContainerResources{
				{
					Name: "container-1",
					Devices: []*podresourcesapi.ContainerDevices{
						{
							ResourceName: "amd.com/gpu",
							DeviceIds: []string{
								"0000:b3:00.0",
							},
						},
					},
				},
			},
		},
		{
			Name:      "pod-y",
			Namespace: "team-a",
			Containers: []*podresourcesapi.ContainerResources{
				{
					Name: "container-1",
					Devices: []*podresourcesapi.ContainerDevices{
						{
							ResourceName: "unknown/resource",
							DeviceIds: []string{
								"12345",
							},
						},
					},
				},
			},
		},
	}
}

func ExistingPodsWithLabelsFixture(t *testing.T) []runtime.Object {
	t.Helper()

	return []runtime.Object{
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "team-a",
				Labels: map[string]string{
					"key":     "value",
					"label_1": "value-1",
					"label_2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-b",
				Namespace: "team-2",
				Labels: map[string]string{
					"key":     "value",
					"label_1": "value-a",
					"label_2": "value-b",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-c",
				Namespace: "team-2",
				Labels: map[string]string{
					"key":     "value",
					"label_1": "value-i",
					"label_2": "value-ii",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-i",
				Namespace: "team-b",
				Labels: map[string]string{
					"key":     "value",
					"label_1": "value-1",
					"label_2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-ii",
				Namespace: "team-b",
				Labels: map[string]string{
					"key":     "value",
					"label_1": "value-1",
					"label_2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-y",
				Namespace: "team-a",
				Labels: map[string]string{
					"key":     "value",
					"label_1": "value-1",
					"label_2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
	}
}

func ExistingPodsWithoutLabelsFixture(t *testing.T) []runtime.Object {
	t.Helper()

	return []runtime.Object{
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "team-a",
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-b",
				Namespace: "team-2",
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-c",
				Namespace: "team-2",
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-i",
				Namespace: "team-b",
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-ii",
				Namespace: "team-b",
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-y",
				Namespace: "team-a",
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
	}
}
