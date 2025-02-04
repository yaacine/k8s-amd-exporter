package kubernetes_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters/domain/pods"
	"github.com/openinnovationai/k8s-amd-exporter/internal/kubernetes"
	"github.com/openinnovationai/k8s-amd-exporter/internal/sdk/e2etests"
	"github.com/openinnovationai/k8s-amd-exporter/internal/sdk/unittests/fakekubelet"
	"github.com/openinnovationai/k8s-amd-exporter/internal/sdk/unittests/k8sfixtures"
	"github.com/openinnovationai/k8s-amd-exporter/internal/sdk/unittests/testlogs"
)

func TestListPods(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	podResourcesSocket := "/var/lib/kubelet/pod-resources/kubelet.sock"
	grpcConn, err := kubernetes.CreateGRPCConn(podResourcesSocket)
	require.NoError(t, err)
	logger := testlogs.NewLogger()
	config := kubernetes.Setup{
		Logger:              logger,
		KubeletConn:         grpcConn,
		CustomResourceNames: []string{},
	}
	k8sClient := kubernetes.NewClient(&config)
	defer k8sClient.Close()

	// When
	pods, err := k8sClient.GetPodsUsingDevices(ctx)
	// Then
	require.NoError(t, err)
	assert.NotEmpty(t, pods)
}

func TestGetPodsUsingDevices(t *testing.T) {
	// Given
	ctx := context.TODO()
	k8sClient := fakekubelet.New(t,
		fakekubelet.WithPodResources(k8sfixtures.NewPodResourcesFixture(t)),
		fakekubelet.WithAMDCustomResourceNames([]string{"amd-custom-resource-name"}),
		fakekubelet.WithNodeName("oi-wn-gpu-amd-01.test.oiai.corp"),
	)
	defer k8sClient.Close()
	want := pods.PodPerDevices{
		"0001:34:00.0": {
			Name:      "pod-1",
			Container: "container-1",
			Namespace: "team-a",
			NodeName:  "oi-wn-gpu-amd-01.test.oiai.corp",
		},

		"80EE/vgpu-2345": {
			Name:      "pod-b",
			Container: "container-1",
			Namespace: "team-2",
			NodeName:  "oi-wn-gpu-amd-01.test.oiai.corp",
		},

		"0000:8e:00.0": {
			Name:      "pod-c",
			Container: "container-1",
			Namespace: "team-2",
			NodeName:  "oi-wn-gpu-amd-01.test.oiai.corp",
		},

		"amd123/gi456": {
			Name:      "pod-i",
			Container: "container-1",
			Namespace: "team-b",
			NodeName:  "oi-wn-gpu-amd-01.test.oiai.corp",
		},

		"0000:b3:00.0": {
			Name:      "pod-ii",
			Namespace: "team-b",
			Container: "container-1",
			NodeName:  "oi-wn-gpu-amd-01.test.oiai.corp",
		},
	}

	// When
	got, err := k8sClient.GetPodsUsingDevices(ctx)

	// Then
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetPodsLabelsOneByOne(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	existingPods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "team-a",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-2",
				Namespace: "team-a",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-a",
					"label-2": "value-b",
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-3",
				Namespace: "team-b",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-i",
					"label-2": "value-ii",
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-4",
				Namespace: "team-c",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-5",
				Namespace: "team-c",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
		},
	}

	logger := testlogs.NewLogger()
	k8sClient := fake.NewClientset(existingPods...)
	setup := kubernetes.Setup{
		Logger:    logger,
		K8SClient: k8sClient,
	}

	kubeClient := kubernetes.NewClient(&setup)

	requiredLabels := []string{
		"label-1", "label-2",
	}
	podsInfo := []pods.PodInfo{
		{
			Name:      "pod-1",
			Namespace: "team-a",
		},
		{
			Name:      "pod-2",
			Namespace: "team-a",
		},
		{
			Name:      "pod-3",
			Namespace: "team-b",
		},
	}

	want := map[string]pods.Labels{
		"team-a/pod-1": {
			"label-1": "value-1",
			"label-2": "value-2",
		},
		"team-a/pod-2": {
			"label-1": "value-a",
			"label-2": "value-b",
		},
		"team-b/pod-3": {
			"label-1": "value-i",
			"label-2": "value-ii",
		},
	}

	// When
	got, err := kubeClient.GetPodsLabels(ctx, podsInfo, requiredLabels)

	// Then
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetPodsLabelsWithinNode(t *testing.T) {
	// Given
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	existingPods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "team-a",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-2",
				Namespace: "team-a",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-a",
					"label-2": "value-b",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-3",
				Namespace: "team-b",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-i",
					"label-2": "value-ii",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-4",
				Namespace: "team-c",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-1",
			},
		},
		&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pod-5",
				Namespace: "team-c",
				Labels: map[string]string{
					"key":     "value",
					"label-1": "value-1",
					"label-2": "value-2",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node-2",
			},
		},
	}

	logger := testlogs.NewLogger()
	k8sClient := fake.NewClientset(existingPods...)
	setup := kubernetes.Setup{
		Logger:    logger,
		K8SClient: k8sClient,
		NodeName:  "node-1",
	}

	kubeClient := kubernetes.NewClient(&setup)

	requiredLabels := []string{
		"label-1", "label-2",
	}
	podsInfo := []pods.PodInfo{
		{
			Name:      "pod-1",
			Namespace: "team-a",
		},
		{
			Name:      "pod-2",
			Namespace: "team-a",
		},
		{
			Name:      "pod-3",
			Namespace: "team-b",
		},
	}

	want := map[string]pods.Labels{
		"team-a/pod-1": {
			"label-1": "value-1",
			"label-2": "value-2",
		},
		"team-a/pod-2": {
			"label-1": "value-a",
			"label-2": "value-b",
		},
		"team-b/pod-3": {
			"label-1": "value-i",
			"label-2": "value-ii",
		},
	}

	// When
	got, err := kubeClient.GetPodsLabels(ctx, podsInfo, requiredLabels)

	// Then
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
