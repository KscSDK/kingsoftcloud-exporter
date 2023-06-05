package collector

import (
	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/instance"
	"github.com/KscSDK/kingsoftcloud-exporter/metric"
	"github.com/KscSDK/kingsoftcloud-exporter/util"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	handlerFactoryMap = make(map[string]func(*KscProductCollector, log.Logger) (ProductHandler, error))
)

// 每个产品的指标处理逻辑
type ProductHandler interface {
	// 获取云监控指标namespace
	GetNamespace() string

	// 对指标元数据做检验, true=可用, false=跳过
	// IsMetricMetaValid(meta *metric.Meta) bool

	// 修改指标元数据
	// ModifyMetricMeta(meta *metric.Meta) error

	// 对指标做校验, true=可用, false=跳过
	// IsMetricValid(m *metric.Metric) bool

	// 修改指标
	// ModifyMetric(m *metric.Metric) error

	GetInstances() ([]instance.KscInstance, error)

	// 获取该指标下符合条件的所有实例, 并生成所有的series
	GetSeriesByInstances(m *metric.Metric, instances []instance.KscInstance) (series []*metric.Series, err error)

	// 获取该指标下符合条件的所有实例, 并生成所有的series
	GetSeries(m *metric.Metric) (series []*metric.Series, err error)
}

// 将对应的产品handler注册到Factory中
func registerHandler(namespace string, _ bool, factory func(*KscProductCollector, log.Logger) (ProductHandler, error)) {
	handlerFactoryMap[namespace] = factory
}

type baseProductHandler struct {
	collector       *KscProductCollector
	monitorQueryKey string
	logger          log.Logger
}

//GetInstances
func (h *baseProductHandler) GetInstances() ([]instance.KscInstance, error) {
	filters := make(map[string]interface{})

	var instances []instance.KscInstance
	var err error
	if config.ExporterRunningMode == config.ExporterMode_Mock {
		instances, err = h.collector.InstanceRepo.ListByMonitors(filters)
	} else {
		instances, err = h.collector.InstanceRepo.ListByFilters(filters)
	}

	if err != nil {
		return nil, err
	}

	return instances, nil
}

//GetSeriesByInstance
func (h *baseProductHandler) GetSeriesByInstances(m *metric.Metric, instances []instance.KscInstance) ([]*metric.Series, error) {
	var seriesSlice []*metric.Series
	for _, i := range instances {
		if len(m.Conf.ExcludeInstances) != 0 && util.IsStrInList(m.Conf.ExcludeInstances, i.GetInstanceID()) {
			continue
		}

		ql := map[string]string{
			"namespace":    m.Meta.Namespace,
			"region":       config.ExporterRunningRegion,
			"instancename": i.GetInstanceName(),
			"instanceid":   i.GetInstanceID(),
			"instanceip":   i.GetInstanceIP(),
		}

		s, err := metric.NewSeries(m, ql, i)
		if err != nil {
			level.Error(h.logger).Log("msg", "Create metric series fail", "metric", m.Meta.MetricName, "instance", i.GetInstanceID())
			continue
		}
		seriesSlice = append(seriesSlice, s)
	}
	return seriesSlice, nil
}

func (h *baseProductHandler) GetSeries(m *metric.Metric) ([]*metric.Series, error) {
	return h.GetSeriesByAll(m)
}

func (h *baseProductHandler) GetSeriesByAll(m *metric.Metric) ([]*metric.Series, error) {
	var seriesSlice []*metric.Series

	var instances []instance.KscInstance
	var err error
	if config.ExporterRunningMode == config.ExporterMode_Mock {
		instances, err = h.collector.InstanceRepo.ListByMonitors(m.Conf.InstanceFilters)
	} else {
		instances, err = h.collector.InstanceRepo.ListByFilters(m.Conf.InstanceFilters)
	}
	if err != nil {
		return nil, err
	}

	for _, i := range instances {
		if len(m.Conf.ExcludeInstances) != 0 && util.IsStrInList(m.Conf.ExcludeInstances, i.GetInstanceID()) {
			continue
		}

		ql := map[string]string{
			"namespace":    m.Meta.Namespace,
			"region":       config.ExporterRunningRegion,
			"instancename": i.GetInstanceName(),
			"instanceid":   i.GetInstanceID(),
			"instanceip":   i.GetInstanceIP(),
		}

		s, err := metric.NewSeries(m, ql, i)
		if err != nil {
			level.Error(h.logger).Log("msg", "Create metric series fail", "metric", m.Meta.MetricName, "instance", i.GetInstanceID())
			continue
		}
		seriesSlice = append(seriesSlice, s)
	}
	return seriesSlice, nil
}
