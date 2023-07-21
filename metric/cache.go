package metric

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

//MetricCache
type MetricCache struct {
	Raw                MetricRepository
	metaCache          map[string]map[string]*Meta //k1=namespace, k2=metricname(小写)
	metaLastReloadTime map[string]int64
	metaLock           sync.Mutex
	logger             log.Logger
}

func (c *MetricCache) GetMeta(conf config.KscMetricConfig, instanceId string) (*Meta, error) {
	err := c.checkMetaNeedReload(conf.Namespace, instanceId)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s-%s", conf.Namespace, instanceId)

	np, exists := c.metaCache[key]
	if !exists {
		return nil, fmt.Errorf("namespace cache not exists")
	}
	m, exists := np[strings.ToLower(conf.MetricName)]
	if !exists {
		return nil, fmt.Errorf("metric cache not exists")
	}
	return m, nil
}

//ListMetrics
func (c *MetricCache) ListMetrics(namespace, instanceId string) ([]*Meta, error) {
	if err := c.checkMetaNeedReload(namespace, instanceId); err != nil {
		return nil, err
	}

	var metaSlice []*Meta

	key := fmt.Sprintf("%s-%s", namespace, instanceId)

	for _, meta := range c.metaCache[key] {
		metaSlice = append(metaSlice, meta)
	}
	return metaSlice, nil
}

//ListMetrics
func (c *MetricCache) ListLocalMetrics(namespace, instanceId string) ([]*Meta, error) {
	if err := c.checkMetaNeedReload(namespace, instanceId); err != nil {
		return nil, err
	}

	var metaSlice []*Meta

	key := fmt.Sprintf("%s-%s", namespace, instanceId)

	for _, meta := range c.metaCache[key] {
		metaSlice = append(metaSlice, meta)
	}
	return metaSlice, nil
}

//ListSamples
func (c *MetricCache) ListSamples(metric *Metric, startTime int64, endTime int64) (samplesList []*Samples, err error) {
	return c.Raw.ListSamples(metric, startTime, endTime)
}

func (c *MetricCache) ListBatchSamples(namespace string, m map[string]*Metric, period int64, startTime int64, endTime int64) (metricSamples map[string][]*Samples, err error) {
	return c.Raw.ListBatchSamples(namespace, m, period, startTime, endTime)
}

func (c *MetricCache) DescribeMonitorData(namespace string, m map[string]*Metric, period int64, startTime int64, endTime int64) (metricSamples map[string][]*Samples, err error) {
	return c.Raw.DescribeMonitorData(namespace, m, period, startTime, endTime)
}

//checkMetaNeedReload 检测是否需要reload缓存的数据
func (c *MetricCache) checkMetaNeedReload(namespace, instanceId string) (err error) {
	key := fmt.Sprintf("%s-%s", namespace, instanceId)

	currentTime := time.Now().Unix()

	v, ok := c.metaLastReloadTime[key]
	if ok && currentTime-v < 60 {
		return nil
	}

	var metas []*Meta
	if config.IsSupportMultiDimensionNamespace(namespace) {
		metas, err = c.Raw.ListMetrics(namespace, instanceId)
	} else {
		metas, err = c.Raw.ListLocalMetrics(namespace, instanceId)
	}

	if err != nil {
		return err
	}
	np, ok := c.metaCache[key]
	if !ok {
		np = map[string]*Meta{}
		c.metaCache[key] = np
	}

	for _, meta := range metas {
		id := fmt.Sprintf("%s.%s", meta.MetricName, instanceId)
		np[id] = meta
	}

	c.metaLastReloadTime[key] = currentTime

	level.Debug(c.logger).Log("msg", "Reload metric meta cache", "namespace", namespace, "num", len(np))
	return
}

//NewMetricCache
func NewMetricCache(repo MetricRepository, logger log.Logger) MetricRepository {
	cache := &MetricCache{
		Raw:                repo,
		metaCache:          map[string]map[string]*Meta{},
		metaLastReloadTime: map[string]int64{},
		logger:             logger,
	}
	return cache
}
