package stats

import (
	"github.com/finitum/node-cli/provider"
	"github.com/finitum/node-cli/stats"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

// PodStatsSummaryHandlerFunc defines the handler for getting pod stats summaries
type PodStatsSummaryHandlerFunc func() (*stats.Summary, error)

type Collector struct {
	p *provider.PodMetricsProvider

	// Node metrics
	usageNanoCores *prometheus.Desc

	availableBytesMemory *prometheus.Desc
	usageBytesMemory     *prometheus.Desc

	availableBytesEphemeral *prometheus.Desc
	capacityBytesEphemeral  *prometheus.Desc
	usedBytesEphemeral      *prometheus.Desc

	availableBytesStorage *prometheus.Desc
	capacityBytesStorage  *prometheus.Desc
	usedBytesStorage      *prometheus.Desc

	pods *prometheus.Desc

	// Pod metrics
	usageNanoCoresPod *prometheus.Desc

	availableBytesMemoryPod *prometheus.Desc
	usageBytesMemoryPod     *prometheus.Desc

	availableBytesEphemeralPod *prometheus.Desc
	capacityBytesEphemeralPod  *prometheus.Desc
	usedBytesEphemeralPod      *prometheus.Desc

	availableBytesStoragePod *prometheus.Desc
	capacityBytesStoragePod  *prometheus.Desc
	usedBytesStoragePod      *prometheus.Desc
}

func NewCollector(p *provider.PodMetricsProvider) Collector {
	return Collector{
		p:                       p,
		usageNanoCores:          prometheus.NewDesc("usage_nano_cores_node", "CPU Usage", []string{"node_name"}, prometheus.Labels{}),
		availableBytesMemory:    prometheus.NewDesc("available_bytes_memory_node", "Available bytes (memory)", []string{"node_name"}, prometheus.Labels{}),
		usageBytesMemory:        prometheus.NewDesc("usage_bytes_memory_node", "Used bytes (memory)", []string{"node_name"}, prometheus.Labels{}),
		availableBytesEphemeral: prometheus.NewDesc("available_bytes_ephemeral_node", "Available bytes (ephemeral)", []string{"node_name"}, prometheus.Labels{}),
		capacityBytesEphemeral:  prometheus.NewDesc("capacity_bytes_ephemeral_node", "Capacity bytes (ephemeral)", []string{"node_name"}, prometheus.Labels{}),
		usedBytesEphemeral:      prometheus.NewDesc("used_bytes_ephemeral_node", "Used bytes (ephemeral)", []string{"node_name"}, prometheus.Labels{}),
		availableBytesStorage:   prometheus.NewDesc("available_bytes_storage_node", "Available bytes (storage)", []string{"node_name"}, prometheus.Labels{}),
		capacityBytesStorage:    prometheus.NewDesc("capacity_bytes_storage_node", "Capacity bytes (storage)", []string{"node_name"}, prometheus.Labels{}),
		usedBytesStorage:        prometheus.NewDesc("used_bytes_storage_node", "Used bytes (storage)", []string{"node_name"}, prometheus.Labels{}),
		pods:                    prometheus.NewDesc("pods_node", "Running pods", []string{"node_name"}, prometheus.Labels{}),

		usageNanoCoresPod:          prometheus.NewDesc("usage_nano_cores_pod", "CPU Usage", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		availableBytesMemoryPod:    prometheus.NewDesc("available_bytes_memory_pod", "Available bytes (memory)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		usageBytesMemoryPod:        prometheus.NewDesc("usage_bytes_memory_pod", "Used bytes (memory)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		availableBytesEphemeralPod: prometheus.NewDesc("available_bytes_ephemeral_pod", "Available bytes (ephemeral)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		capacityBytesEphemeralPod:  prometheus.NewDesc("capacity_bytes_ephemeral_pod", "Capacity bytes (ephemeral)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		usedBytesEphemeralPod:      prometheus.NewDesc("used_bytes_ephemeral_pod", "Used bytes (ephemeral)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		availableBytesStoragePod:   prometheus.NewDesc("available_bytes_storage_pod", "Available bytes (storage)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		capacityBytesStoragePod:    prometheus.NewDesc("capacity_bytes_storage_pod", "Capacity bytes (storage)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
		usedBytesStoragePod:        prometheus.NewDesc("used_bytes_storage_pod", "Used bytes (storage)", []string{"pod_name", "pod_namespace", "pod_uid"}, prometheus.Labels{}),
	}
}

func (s Collector) Collect(metrics chan<- prometheus.Metric) {
	log.Printf("got metrics collect\n")

	// Collect st
	if s.p == nil {
		return
	}

	st, err := (*s.p).GetStatsSummary()
	if err != nil {
		return
	}

	// Add node stats
	metrics <- prometheus.MustNewConstMetric(s.usageNanoCores, prometheus.GaugeValue, float64(st.Node.UsageNanoCores), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.availableBytesMemory, prometheus.GaugeValue, float64(st.Node.AvailableBytesMemory), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.usageBytesMemory, prometheus.GaugeValue, float64(st.Node.UsageBytesMemory), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.availableBytesEphemeral, prometheus.GaugeValue, float64(st.Node.AvailableBytesEphemeral), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.capacityBytesEphemeral, prometheus.GaugeValue, float64(st.Node.CapacityBytesEphemeral), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.usedBytesEphemeral, prometheus.GaugeValue, float64(st.Node.UsedBytesEphemeral), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.availableBytesStorage, prometheus.GaugeValue, float64(st.Node.AvailableBytesStorage), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.capacityBytesStorage, prometheus.GaugeValue, float64(st.Node.CapacityBytesStorage), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.usedBytesStorage, prometheus.GaugeValue, float64(st.Node.UsedBytesStorage), st.Node.Name)
	metrics <- prometheus.MustNewConstMetric(s.pods, prometheus.GaugeValue, float64(st.Node.Pods), st.Node.Name)

	// Add pod stats
	for _, pod := range st.Pods {
		ref := pod.PodRef
		metrics <- prometheus.MustNewConstMetric(s.usageNanoCoresPod, prometheus.GaugeValue, float64(pod.UsageNanoCores), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.availableBytesMemoryPod, prometheus.GaugeValue, float64(pod.AvailableBytesMemory), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.usageBytesMemoryPod, prometheus.GaugeValue, float64(pod.UsageBytesMemory), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.availableBytesEphemeralPod, prometheus.GaugeValue, float64(pod.AvailableBytesEphemeral), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.capacityBytesEphemeralPod, prometheus.GaugeValue, float64(pod.CapacityBytesEphemeral), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.usedBytesEphemeralPod, prometheus.GaugeValue, float64(pod.UsedBytesEphemeral), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.availableBytesStoragePod, prometheus.GaugeValue, float64(pod.AvailableBytesStorage), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.capacityBytesStoragePod, prometheus.GaugeValue, float64(pod.CapacityBytesStorage), ref.Name, ref.Namespace, ref.UID)
		metrics <- prometheus.MustNewConstMetric(s.usedBytesStoragePod, prometheus.GaugeValue, float64(pod.UsedBytesStorage), ref.Name, ref.Namespace, ref.UID)
	}
}

func (s Collector) Describe(desc chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(s, desc)
}
