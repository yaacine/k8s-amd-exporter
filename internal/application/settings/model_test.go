package settings_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/openinnovationai/platform/infra/amd/amd_smi_exporter_v2/internal/application/settings"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	// Given
	err := os.Setenv("AMD_EXPORTER_LOG_LEVEL", "development")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_WEB_SERVER_PORT", "8080")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_KUBELET_SOCKET_PATH", "/any/path")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_RESOURCE_NAMES", "custom1,custom2,custom3")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_WITH_KUBERNETES", "true")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_NODE_NAME", "oi-wn-gpu-amd-01.test.oiai.corp")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_POD_NAME", "amd-smi-exporter-v2-2")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_NAMESPACE", "amdexporter-amdsmiexporter")
	require.NoError(t, err)
	err = os.Setenv("AMD_EXPORTER_POD_LABELS", "label_1,label_2,label_3")
	require.NoError(t, err)

	want := &settings.Configuration{
		LogLevel:          "development",
		WebServerPort:     8080,
		KubeletSocketPath: "/any/path",
		AMDResourceNames:  []string{"custom1", "custom2", "custom3"},
		NodeName:          "oi-wn-gpu-amd-01.test.oiai.corp",
		PodName:           "amd-smi-exporter-v2-2",
		PodNamespace:      "amdexporter-amdsmiexporter",
		PodLabels:         []string{"label_1", "label_2", "label_3"},
		WithKubernetes:    true,
	}

	// When
	got, err := settings.Load()

	// Then
	require.NoError(t, err)
	assert.Equal(t, want, got)

	err = os.Unsetenv("AMD_EXPORTER_LOG_LEVEL")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_WEB_SERVER_PORT")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_KUBELET_SOCKET_PATH")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_RESOURCE_NAMES")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_WITH_KUBERNETES")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_NODE_NAME")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_POD_NAME")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_NAMESPACE")
	require.NoError(t, err)
	err = os.Unsetenv("AMD_EXPORTER_POD_LABELS")
	require.NoError(t, err)
}
