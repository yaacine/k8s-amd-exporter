package pods_test

import (
	"testing"

	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters/domain/pods"
	"github.com/stretchr/testify/assert"
)

func TestNamespacedName(t *testing.T) {
	// Given
	tests := map[string]struct {
		podInfo pods.PodInfo
		want    string
	}{
		"normal": {
			podInfo: pods.PodInfo{
				Name:      "amd-pod-1",
				Namespace: "team-a",
			},
			want: "team-a/amd-pod-1",
		},
		"default": {
			podInfo: pods.PodInfo{
				Name: "amd-pod-1",
			},
			want: "default/amd-pod-1",
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			// When
			got := testData.podInfo.NamespacedName()
			// Then
			assert.Equal(t, testData.want, got)
		})
	}
}

func TestLabelsSortedKey(t *testing.T) {
	// Given
	tests := map[string]struct {
		labels pods.Labels
		want   []string
	}{
		"unsorted_keys": {
			labels: pods.Labels{
				"label-z": "value-z",
				"label-a": "value-a",
				"label-m": "value-m",
			},
			want: []string{"label-a", "label-m", "label-z"},
		},
		"nil": {
			labels: nil,
			want:   []string{},
		},
		"empty": {
			labels: pods.Labels{},
			want:   []string{},
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			// When
			got := testData.labels.SortKeys()
			// Then
			assert.Equal(t, testData.want, got)
		})
	}
}
