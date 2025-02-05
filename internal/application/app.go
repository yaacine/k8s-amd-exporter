package application

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/openinnovationai/k8s-amd-exporter/internal/amd"
	"github.com/openinnovationai/k8s-amd-exporter/internal/application/logs"
	"github.com/openinnovationai/k8s-amd-exporter/internal/application/settings"
	"github.com/openinnovationai/k8s-amd-exporter/internal/application/web"
	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters"
	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters/domain/gpus"
	"github.com/openinnovationai/k8s-amd-exporter/internal/kubernetes"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	version    string
	buildDate  string
	commitHash string
)

type Application struct {
	configuration *settings.Configuration
	logger        *slog.Logger
	webServer     *web.Server
	k8sClient     *kubernetes.Client
	exporter      *exporters.Exporter
	gpuCards      [24]gpus.Card

	version    string
	buildDate  string
	commitHash string
}

// NewApplication instantiates exporter application.
func NewApplication() *Application {
	newApplication := Application{
		version:    version,
		buildDate:  buildDate,
		commitHash: commitHash,
	}

	return &newApplication
}

// Run starts this application, loading settings and injecting dependencies.
func (a *Application) Run() error {
	a.printInfo()

	err := a.loadConfiguration()
	if err != nil {
		slog.Error("loading application configuration", slog.String("error", err.Error()))

		return fmt.Errorf("unable to start exporter: %w", err)
	}

	// initialize logger
	a.initializeLogger()

	// initialize k8s connection
	err = a.initializeK8SConnection()
	if err != nil {
		a.logger.Error("initializing kubernetes client", slog.String("error", err.Error()))

		return fmt.Errorf("unable to start exporter: %w", err)
	}

	defer a.closeResources()

	err = a.initializeGPUInformation()
	if err != nil {
		a.logger.Error("initializing gpu products information", slog.String("error", err.Error()))

		return fmt.Errorf("unable to start exporter: %w", err)
	}

	a.initializeExporter()
	a.registryPrometheusExporter()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startWebServer(ctx)

	return nil
}

func (a *Application) printInfo() {
	slog.Info(
		"starting amd exporter",
		slog.String("version", a.version),
		slog.String("commit", a.commitHash),
		slog.String("build-date", a.buildDate),
	)
}

// loadConfiguration loads exporter configuration.
func (a *Application) loadConfiguration() error {
	slog.Info("loading configuration")

	newConfiguration, err := settings.Load()
	if err != nil {
		return fmt.Errorf("unable to load exporter configuration: %w", err)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("configuration parameters", slog.Any("config", newConfiguration))

	a.configuration = newConfiguration

	return nil
}

func (a *Application) initializeLogger() {
	slog.Info("initializing logger")

	a.logger = logs.MakeLogger(a.configuration)

	slog.SetDefault(a.logger)
}

func (a *Application) startWebServer(ctx context.Context) {
	a.logger.Info("starting web server", slog.Uint64("port", uint64(a.configuration.WebServerPort)))

	webServerSetup := web.Setup{
		Logger: a.logger,
		Port:   a.configuration.WebServerPort,
	}

	a.webServer = web.NewServer(&webServerSetup)

	a.webServer.Start(ctx)
}

func (a *Application) initializeK8SConnection() error {
	if !a.configuration.WithKubernetes {
		a.logger.Info("application was configured not to run with the Kubernetes client")

		return nil
	}

	a.logger.Info("initializing the Kubernetes client")

	grpcConn, err := kubernetes.CreateGRPCConn(a.configuration.KubeletSocketPath)
	if err != nil {
		return fmt.Errorf("unable to create grpc connection to kubernetes: %w", err)
	}

	kubeConfig, err := kubernetes.NewKubeConfig()
	if err != nil {
		return fmt.Errorf("unable to create grpc kube config: %w", err)
	}

	clientSetConn, err := kubernetes.NewK8SClient(kubeConfig)
	if err != nil {
		return fmt.Errorf("unable to create kubernetes clientset: %w", err)
	}

	if a.configuration.NodeName == "" {
		a.logger.Info(
			"node name setup is empty",
			slog.String("envvar", settings.NodeNameEnvVar),
			slog.String("hint", "spec.[]containers.[]env.valueFrom.fieldRef.fieldPath: spec.nodeName"),
		)
	}

	if a.configuration.PodName == "" {
		a.logger.Info(
			"exporter pod name setup is empty",
			slog.String("envvar", settings.PodNameEnvVar),
			slog.String("hint", "spec.[]containers.[]env.valueFrom.fieldRef.fieldPath: metadata.name"),
		)
	}

	if a.configuration.PodNamespace == "" {
		a.logger.Info(
			"exporter pod namespace setup is empty",
			slog.String("envvar", settings.PodNamespaceEnvVar),
			slog.String("hint", "spec.[]containers.[]env.valueFrom.fieldRef.fieldPath: metadata.namespace"),
		)
	}

	k8sClientSettings := kubernetes.Setup{
		KubeletConn:         grpcConn,
		Logger:              a.logger,
		CustomResourceNames: a.configuration.AMDResourceNames,
		NodeName:            a.configuration.NodeName,
		PodName:             a.configuration.PodName,
		PodNamespace:        a.configuration.PodNamespace,
		K8SClient:           clientSetConn,
	}

	a.k8sClient = kubernetes.NewClient(&k8sClientSettings)

	return nil
}

func (a *Application) initializeGPUInformation() error {
	gpuCards, err := amd.GetGpuProductNames()
	if err != nil {
		return fmt.Errorf("unable to get gpu products from environment: %w", err)
	}

	a.gpuCards = gpuCards

	return nil
}

func (a *Application) initializeExporter() {
	a.logger.Info("initializing the metrics exporter")

	amdScanner := amd.NewScanner(a.logger)
	settings := exporters.Setup{
		K8SClient:      a.k8sClient,
		CardsInfo:      a.gpuCards,
		Logger:         a.logger,
		OIPLabels:      a.configuration.PodLabels,
		WithKubernetes: a.configuration.WithKubernetes,
		GetMetricsFunc: func() gpus.AMDParams {
			return amdScanner.Scan()
		},
	}

	a.exporter = exporters.NewExporter(&settings)
}

func (a *Application) registryPrometheusExporter() {
	a.logger.Info("registering exporter with prometheus")
	// Make Prometheus client aware of our collector.
	prometheus.MustRegister(a.exporter)
}

func (a *Application) closeResources() {
	a.logger.Info("closing application resources")

	if a.configuration.WithKubernetes {
		a.logger.Info("closing kubernetes connection")
		a.k8sClient.Close()
	}
}
