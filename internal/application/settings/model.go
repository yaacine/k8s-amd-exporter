package settings

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

// environment variable names.
const (
	NodeNameEnvVar     string = "AMD_EXPORTER_NODE_NAME"
	PodNameEnvVar      string = "AMD_EXPORTER_POD_NAME"
	PodNamespaceEnvVar string = "AMD_EXPORTER_NAMESPACE"
)

// Configuration contains the parameters required for this exporter to work.
type Configuration struct {
	// LogLevel it could be 'production' or 'development'.
	LogLevel string `env:"AMD_EXPORTER_LOG_LEVEL" envDefault:"development"`
	// port for web server
	WebServerPort uint `env:"AMD_EXPORTER_WEB_SERVER_PORT" envDefault:"2021"`
	// port for web server
	KubeletSocketPath string `env:"AMD_EXPORTER_KUBELET_SOCKET_PATH" envDefault:"/var/lib/kubelet/pod-resources/kubelet.sock"`
	// AMD custom gpu resource names
	AMDResourceNames []string `env:"AMD_EXPORTER_RESOURCE_NAMES"`
	// Using kubernetes resources flag
	WithKubernetes bool `env:"AMD_EXPORTER_WITH_KUBERNETES" envDefault:"true"`
	// Kubernetes node name where this application is running
	NodeName string `env:"AMD_EXPORTER_NODE_NAME"`
	// Kubernetes pod name where this application is running
	PodName string `env:"AMD_EXPORTER_POD_NAME"`
	// Kubernetes pod namespace where this application is running
	PodNamespace string `env:"AMD_EXPORTER_NAMESPACE"`
	// Kubernetes pod labels to be added to exporter labels.
	PodLabels []string `env:"AMD_EXPORTER_POD_LABELS"`
}

func Load() (*Configuration, error) {
	var cfg Configuration

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to load settings: %w", err)
	}

	return &cfg, nil
}
