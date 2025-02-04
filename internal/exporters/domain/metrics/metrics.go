package metrics

import (
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"

	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters/domain/gpus"
	"github.com/openinnovationai/k8s-amd-exporter/internal/exporters/domain/pods"
	"github.com/prometheus/client_golang/prometheus"
)

// CustomMetric defines data required to build a metric.
type CustomMetric struct {
	Name      string
	Namespace string
	Subsystem string
	HelpText  string
	Type      prometheus.ValueType
	// Labels The metric's variable label dimensions.
	Labels []string
	// Divide indicates that the metric value should be divided by the given divisor.
	Divide  bool
	Divisor float64
}

// AMDMetrics set of prometheus metrics to be collected from amd resources.
type AMDMetrics struct {
	DataDesc       *CustomMetric
	CoreEnergy     *CustomMetric
	SocketEnergy   *CustomMetric
	BoostLimit     *CustomMetric
	SocketPower    *CustomMetric
	PowerLimit     *CustomMetric
	ProchotStatus  *CustomMetric
	Sockets        *CustomMetric
	Threads        *CustomMetric
	ThreadsPerCore *CustomMetric
	NumGPUs        *CustomMetric
	GPUDevID       *CustomMetric
	GPUPowerCap    *CustomMetric
	GPUPower       *CustomMetric
	GPUTemperature *CustomMetric
	GPUSCLK        *CustomMetric
	GPUMCLK        *CustomMetric
	GPUUsage       *CustomMetric
	GPUMemoryUsage *CustomMetric
	CardsInfo      [gpus.MaxNumGPUDevices]gpus.Card
	K8SResources   map[string][]pods.PodInfo
	Data           gpus.AMDParamsHandler // This is the Scan() function handle
	logger         *slog.Logger
	withKubernetes bool
}

// Setup contains objects required to process metrics.
type Setup struct {
	AMDParamsHandler gpus.AMDParamsHandler
	Logger           *slog.Logger
	WithKubernetes   bool
}

// metric labels.
const (
	podNameLabel       string = "exported_pod"
	namespaceNameLabel string = "exported_namespace"
	containerNameLabel string = "exported_container"
	nodeNameLabel      string = "exported_node"
	productNameLabel   string = "productname"
	deviceNameLabel    string = "device"

	deviceIDPrefix           string = "amd"
	amdMetricHelpTextDefault string = "AMD Params" // The metric's help text.
)

// metric common values.
const amdNamespace string = "amd"

// labelPrefix in case you want prefix your labels with "label" word.
const labelPrefixPattern = "label_%s"

// Prometheus label naming convention regex.
var validLabelRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// Replace invalid characters with underscores in labels.
var replacer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// NewAMDMetrics creates AMD metrics based on given handler and a flag to indicate
// if k8s resources are needed.
func NewAMDMetrics(settings *Setup) *AMDMetrics {
	newAMDMetrics := &AMDMetrics{
		withKubernetes: settings.WithKubernetes,
		Data:           settings.AMDParamsHandler,
		logger:         settings.Logger,
	}

	return newAMDMetrics.initializeMetrics()
}

// initializeMetrics initializes prometheus metric descriptions.
func (a *AMDMetrics) initializeMetrics() *AMDMetrics {
	a.DataDesc = &CustomMetric{
		Name:     "amd_data",
		HelpText: amdMetricHelpTextDefault,
		Labels:   []string{"socket"},
	}
	a.CoreEnergy = newAMDCounterMetric("core_energy", "thread")
	a.SocketEnergy = newAMDCounterMetric("socket_energy", "socket")
	a.BoostLimit = newAMDGaugeMetric("boost_limit", "thread")
	a.SocketPower = newAMDGaugeMetric("socket_power", "socket")
	a.PowerLimit = newAMDGaugeMetricWithName("power_limit")
	a.ProchotStatus = newAMDGaugeMetricWithName("prochot_status")
	a.Sockets = newAMDGaugeMetricWithName("num_sockets")
	a.Threads = newAMDGaugeMetricWithName("num_threads")
	a.ThreadsPerCore = newAMDGaugeMetricWithName("num_threads_per_core")
	a.NumGPUs = newAMDGaugeMetricWithName("num_gpus")
	a.GPUDevID = newAMDGPUGaugeMetric("gpu_dev_id")
	a.GPUPowerCap = newAMDGPUGaugeMetric("gpu_power_cap").
		WithDivisor(1e6)
	a.GPUPower = newAMDGPUCounterMetric("gpu_power").
		WithDivisor(1e6)
	a.GPUTemperature = newAMDGPUGaugeMetric("gpu_current_temperature").
		WithDivisor(1e3)
	a.GPUSCLK = newAMDGPUGaugeMetric("gpu_SCLK").
		WithDivisor(1e6)
	a.GPUMCLK = newAMDGPUGaugeMetric("gpu_MCLK").
		WithDivisor(1e6)
	a.GPUUsage = newAMDGPUGaugeMetric("gpu_use_percent")
	a.GPUMemoryUsage = newAMDGPUGaugeMetric("gpu_memory_use_percent")

	return a
}

func newAMDGaugeMetric(name string, label ...string) *CustomMetric {
	return newAMDMetric(name, prometheus.GaugeValue, label...)
}

func newAMDCounterMetric(name string, label ...string) *CustomMetric {
	return newAMDMetric(name, prometheus.CounterValue, label...)
}

func newAMDMetric(name string, mType prometheus.ValueType, label ...string) *CustomMetric {
	return &CustomMetric{
		Name:      name,
		Namespace: amdNamespace,             // metric namespace
		HelpText:  amdMetricHelpTextDefault, // The metric's help text.
		Labels:    label,                    // The metric's variable label dimensions.
		Type:      mType,
	}
}

func newAMDGPUGaugeMetric(name string) *CustomMetric {
	return newAMDGPUMetric(name, prometheus.GaugeValue)
}

func newAMDGPUCounterMetric(name string) *CustomMetric {
	return newAMDGPUMetric(name, prometheus.CounterValue)
}

func newAMDGaugeMetricWithName(name string) *CustomMetric {
	return newAMDMetric(name, prometheus.GaugeValue, name)
}

func newAMDGPUMetric(name string, mType prometheus.ValueType) *CustomMetric {
	return newAMDMetric(name, mType, name, productNameLabel, deviceNameLabel)
}

// k8sVariableLabels return list of kubernetes labels required in metrics.
func k8sVariableLabels() []string {
	return []string{podNameLabel, containerNameLabel, namespaceNameLabel, nodeNameLabel}
}

// WithDivisor enable dividing metric value by given divisor.
func (c *CustomMetric) WithDivisor(divisor float64) *CustomMetric {
	c.Divide = true
	c.Divisor = divisor

	return c
}

// buildPrometheusMetric builds prometheus metric based on given value and metric configuration.
func (c *CustomMetric) buildPrometheusMetric(value float64, labelValues ...string) prometheus.Metric {
	value = c.transformValue(value)

	return prometheus.MustNewConstMetric(
		c.NewDesc(),
		c.Type,
		value,
		labelValues...,
	)
}

// transformValue transform metric value to the desired format.
func (c *CustomMetric) transformValue(value float64) float64 {
	if c.divisionRequired() {
		value /= c.Divisor
	}

	return value
}

// divisionRequired checks if metric value should be divided.
func (c *CustomMetric) divisionRequired() bool {
	return c.Divide && c.Divisor > 0
}

// NewDesc allocates and initializes a new prometheus Desc.
func (c *CustomMetric) NewDesc(additionalLabels ...string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(c.Namespace, c.Subsystem, c.Name),
		c.HelpText,                            // The metric's help text.
		append(c.Labels, additionalLabels...), // The metric's variable label dimensions.
		nil,                                   // The metric's constant label dimensions.
	)
}

// CollectAndBuildMetrics scans amd data and build a collection of metrics.
func (a *AMDMetrics) CollectAndBuildMetrics() []prometheus.Metric {
	data := a.Data() // Scan AMD metrics

	a.logger.Debug("scanning", slog.Any("amd-params", data))

	metrics := make([]prometheus.Metric, 0)

	metrics = append(metrics, buildMetrics(data.CoreEnergy[:data.Threads], data.Threads, a.CoreEnergy)...)
	metrics = append(metrics, buildMetrics(data.CoreBoost[:data.Threads], data.Threads, a.BoostLimit)...)
	metrics = append(metrics, buildMetrics(data.SocketEnergy[:data.Sockets], data.Sockets, a.SocketEnergy)...)
	metrics = append(metrics, buildMetrics(data.SocketPower[:data.Sockets], data.Sockets, a.SocketPower)...)
	metrics = append(metrics, buildMetrics(data.PowerLimit[:data.Sockets], data.Sockets, a.PowerLimit)...)
	metrics = append(metrics, buildMetrics(data.ProchotStatus[:data.Sockets], data.Sockets, a.ProchotStatus)...)

	// GPU metrics
	metrics = append(metrics, a.buildGPUMetrics(data.GPUDevID[:data.NumGPUs], data.NumGPUs, a.GPUDevID)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUPowerCap[:data.NumGPUs], data.NumGPUs, a.GPUPowerCap)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUPower[:data.NumGPUs], data.NumGPUs, a.GPUPower)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUTemperature[:data.NumGPUs], data.NumGPUs, a.GPUTemperature)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUSCLK[:data.NumGPUs], data.NumGPUs, a.GPUSCLK)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUMCLK[:data.NumGPUs], data.NumGPUs, a.GPUMCLK)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUUsage[:data.NumGPUs], data.NumGPUs, a.GPUUsage)...)
	metrics = append(metrics, a.buildGPUMetrics(data.GPUMemoryUsage[:data.NumGPUs], data.NumGPUs, a.GPUMemoryUsage)...)

	metrics = append(metrics, a.resourceGroupMetrics(&data)...)

	return metrics
}

// buildMetrics builds prometheus metric based on given amd metric.
func buildMetrics(
	data []float64,
	attrValue uint,
	metric *CustomMetric,
) []prometheus.Metric {
	if attrValue == 0 {
		return nil
	}

	metrics := make([]prometheus.Metric, attrValue)

	for i := range data {
		metrics[i] = metric.buildPrometheusMetric(data[i], strconv.Itoa(i))
	}

	return metrics
}

// buildGPUMetrics builds prometheus metric based on given amd gpu metric.
func (a *AMDMetrics) buildGPUMetrics(
	data []float64,
	attrValue uint,
	metric *CustomMetric,
) []prometheus.Metric {
	if attrValue == 0 {
		return nil
	}

	var metrics []prometheus.Metric

	for i := range data {
		metrics = append(metrics, a.newMetricWithResources(metric, data[i], i)...)
	}

	return metrics
}

// resourceGroupMetrics build global metrics such as sockets, thread and number of GPUs.
func (a *AMDMetrics) resourceGroupMetrics(params *gpus.AMDParams) []prometheus.Metric {
	return []prometheus.Metric{
		a.Sockets.buildPrometheusMetric(float64(params.Sockets), ""),
		a.Threads.buildPrometheusMetric(float64(params.Threads), ""),
		a.ThreadsPerCore.buildPrometheusMetric(float64(params.ThreadsPerCore), ""),
		a.NumGPUs.buildPrometheusMetric(float64(params.NumGPUs), ""),
	}
}

// newMetricWithResources map given GPU card metric with pod
// using it. If there is no any pod using this card then
// a prometheus metric is created with pod labels.
func (a *AMDMetrics) newMetricWithResources(
	metric *CustomMetric,
	value float64, cardIndex int,
) []prometheus.Metric {
	labelValues := a.commonGPULabelValues(cardIndex)

	if !a.withKubernetes {
		return []prometheus.Metric{
			metric.buildPrometheusMetric(value, labelValues...),
		}
	}

	podsInfo, exist := a.K8SResources[a.CardsInfo[cardIndex].PCIBus]
	if !exist {
		return []prometheus.Metric{
			metric.buildPrometheusMetric(value, labelValues...),
		}
	}

	metrics := make([]prometheus.Metric, 0, len(podsInfo))

	for _, p := range podsInfo {
		newLabels, newLabelValues := buildK8SPodLabelValues(p, labelValues)

		metrics = append(metrics,
			prometheus.MustNewConstMetric(
				metric.NewDesc(newLabels...),
				metric.Type,
				metric.transformValue(value),
				newLabelValues...),
		)
	}

	return metrics
}

// commonGPULabelValues returns common GPU labels.
func (a *AMDMetrics) commonGPULabelValues(cardIndex int) []string {
	return []string{
		strconv.Itoa(cardIndex),
		a.CardsInfo[cardIndex].Cardseries,
		buildDeviceLabelValue(cardIndex),
	}
}

// buildDeviceIDLabelValue build device label name.
func buildDeviceLabelValue(cardIndex int) string {
	return fmt.Sprintf("%s%d", deviceIDPrefix, cardIndex)
}

// buildK8SPodLabelValues return 2 slices of pod labels and its respective values.
func buildK8SPodLabelValues(pod pods.PodInfo, existingLabelValues []string) ([]string, []string) {
	labels := k8sVariableLabels()
	values := slices.Concat(
		existingLabelValues,
		[]string{
			pod.Name,
			pod.Container,
			pod.Namespace,
			pod.NodeName,
		},
	)

	// adding existing labels within pod
	for _, key := range pod.Labels.SortKeys() {
		labels = append(labels, formatLabel(key, true))
		values = append(values, pod.Labels[key])
	}

	return labels, values
}

// formatLabel format label to follow prometheus conventions.
func formatLabel(label string, withPrefix bool) string {
	// Prometheus label naming convention regex.
	if validLabelRegex.MatchString(label) {
		return label // Label is already valid.
	}

	// Replace invalid characters with underscores.
	sanitized := replacer.ReplaceAllString(label, "_")

	// Ensure the label starts with a letter or underscore.
	if sanitized[0] >= '0' && sanitized[0] <= '9' {
		sanitized = "_" + sanitized
	}

	if withPrefix {
		sanitized = fmt.Sprintf(labelPrefixPattern, sanitized)
	}

	return sanitized
}
