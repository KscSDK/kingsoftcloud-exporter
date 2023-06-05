package metric

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/util"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type SeriesCache struct {
	Series map[string]*Series // 包含的多个时间线
	// need cache it, because some cases DescribeBaseMetrics/GetMonitorData dims not match
	LabelNames map[string]struct{}
}

func newCache() *SeriesCache {
	return &SeriesCache{
		Series:     make(map[string]*Series),
		LabelNames: make(map[string]struct{}),
	}
}

type Desc struct {
	FQName string
	Help   string
}

type MetricSet struct {
	// Namespace, each cloud product will have a namespace
	Namespace *string `json:"namespace,omitempty"`

	MetricName *string `json:"metricName,omitempty"`

	MetricDesc *string `json:"metricDesc,omitempty"`

	InstanceID *string `json:"InstanceId,omitempty"`

	//metric collection or push interval
	Interval *string `json:"interval,omitempty"`

	Statistics *[]string `json:"statistics,omitempty"`

	Dimensions *map[string]string `json:"-"`

	//metric value type
	ValueType *string `json:"type,omitempty"`

	Unit *string `json:"Unit,omitempty" name:"unit"`
}

//Metric 代表一个指标, 包含多个时间线
type Metric struct {
	Id           string
	Meta         *Meta  // 指标元数据
	Labels       *Label // 指标labels
	SeriesCache  *SeriesCache
	StatPromDesc map[string]Desc // 按统计纬度的Desc, max、min、avg、last
	Conf         *MetricConf
	seriesLock   sync.Mutex
}

func GetLatestPromMetricsV2(repo MetricRepository, metrics map[string]*Metric, logger log.Logger) (map[string][]prometheus.Metric, error) {
	// var st int64
	now := time.Now().Unix()
	st := now - int64(120)
	et := now - int64(60)

	metricSamples, err := repo.ListBatchSamples(metrics, st, et)
	if err != nil {
		return nil, err
	}

	pms := make(map[string][]prometheus.Metric)

	for metricID, samplesList := range metricSamples {
		m, isExist := metrics[metricID]
		if !isExist {
			continue
		}

		for _, samples := range samplesList {
			for st, desc := range m.StatPromDesc {
				var point *Sample
				switch st {
				case "last":
					point, err = samples.GetLatestPoint()
					if err != nil {
						return nil, err
					}
				case "max":
					point, err = samples.GetMaxPoint()
					if err != nil {
						return nil, err
					}
				case "min":
					point, err = samples.GetMinPoint()
					if err != nil {
						return nil, err
					}
				case "avg":
					point, err = samples.GetAvgPoint()
					if err != nil {
						return nil, err
					}
				}

				var names []string
				var values []string

				labels := m.Labels.GetValues(samples.Series.QueryLabels, samples.Series.Instance)

				labels = map[string]string{
					"namespace":    m.Meta.Namespace,
					"region":       config.ExporterRunningRegion,
					"instancename": samples.Series.Instance.GetInstanceName(),
					"instanceid":   samples.Series.Instance.GetInstanceID(),
					"instanceip":   samples.Series.Instance.GetInstanceIP(),
				}

				for labelName, labelValue := range labels {
					names = append(names, util.ToUnderlineLower(labelName))
					values = append(values, labelValue)
				}

				if m.Meta.m.Dimensions != nil || len(*m.Meta.m.Dimensions) > 0 {
					for k, v := range *m.Meta.m.Dimensions {
						names = append(names, util.ToUnderlineLower(k))
						values = append(values, v)
					}
				}

				newDesc := prometheus.NewDesc(
					desc.FQName,
					desc.Help,
					names,
					nil,
				)

				var pm prometheus.Metric
				if m.Conf.StatDelaySeconds > 0 {
					pm = prometheus.NewMetricWithTimestamp(time.Unix(int64(point.Timestamp), 0), prometheus.MustNewConstMetric(
						newDesc,
						prometheus.GaugeValue,
						point.Value,
						values...,
					))
				} else {
					pm = prometheus.MustNewConstMetric(
						newDesc,
						prometheus.GaugeValue,
						point.Value,
						values...,
					)
				}
				var lock sync.Mutex

				lock.Lock()
				if _, isExists := pms[metricID]; !isExists {
					pms[metricID] = make([]prometheus.Metric, 0, len(samplesList))
				}
				pms[metricID] = append(pms[metricID], pm)
				lock.Unlock()
			}
		}
	}
	return pms, nil
}

//GetLatestPromMetrics
func GetLatestPromMetrics(repo MetricRepository, metrics map[string]*Metric, logger log.Logger) (pms []prometheus.Metric, err error) {
	// var st int64
	now := time.Now().Unix()
	st := now - int64(120)
	et := now - int64(60)

	metricSamples, err := repo.ListBatchSamples(metrics, st, et)
	if err != nil {
		return nil, err
	}

	for metricID, samplesList := range metricSamples {
		m, isExist := metrics[metricID]
		if !isExist {
			continue
		}

		for _, samples := range samplesList {
			for st, desc := range m.StatPromDesc {
				var point *Sample
				switch st {
				case "last":
					point, err = samples.GetLatestPoint()
					if err != nil {
						return nil, err
					}
				case "max":
					point, err = samples.GetMaxPoint()
					if err != nil {
						return nil, err
					}
				case "min":
					point, err = samples.GetMinPoint()
					if err != nil {
						return nil, err
					}
				case "avg":
					point, err = samples.GetAvgPoint()
					if err != nil {
						return nil, err
					}
				}

				var names []string
				var values []string

				labels := m.Labels.GetValues(samples.Series.QueryLabels, samples.Series.Instance)

				labels = map[string]string{
					"namespace":    m.Meta.Namespace,
					"region":       config.ExporterRunningRegion,
					"instancename": samples.Series.Instance.GetInstanceName(),
					"instanceid":   samples.Series.Instance.GetInstanceID(),
					"instanceip":   samples.Series.Instance.GetInstanceIP(),
				}

				for labelName, labelValue := range labels {
					names = append(names, util.ToUnderlineLower(labelName))
					values = append(values, labelValue)
				}

				if m.Meta.m.Dimensions != nil || len(*m.Meta.m.Dimensions) > 0 {
					for k, v := range *m.Meta.m.Dimensions {
						names = append(names, util.ToUnderlineLower(k))
						values = append(values, v)
					}
				}

				newDesc := prometheus.NewDesc(
					desc.FQName,
					desc.Help,
					names,
					nil,
				)

				var pm prometheus.Metric
				if m.Conf.StatDelaySeconds > 0 {
					pm = prometheus.NewMetricWithTimestamp(time.Unix(int64(point.Timestamp), 0), prometheus.MustNewConstMetric(
						newDesc,
						prometheus.GaugeValue,
						point.Value,
						values...,
					))
				} else {
					pm = prometheus.MustNewConstMetric(
						newDesc,
						prometheus.GaugeValue,
						point.Value,
						values...,
					)
				}
				pms = append(pms, pm)
			}
		}
	}
	return
}

//GetLatestPromMetrics
func (m *Metric) GetLatestPromMetrics(repo MetricRepository) (pms []prometheus.Metric, err error) {
	var st int64
	et := int64(0)
	now := time.Now().Unix()

	if m.Conf.StatDelaySeconds > 0 {
		st = now - m.Conf.StatNumSamples*m.Conf.StatPeriodSeconds - m.Conf.StatDelaySeconds
		et = now - m.Conf.StatDelaySeconds
	} else {
		st = now - m.Conf.StatNumSamples*m.Conf.StatPeriodSeconds
		et = now
	}

	samplesList, err := repo.ListSamples(m, st, et)
	if err != nil {
		return nil, err
	}

	for _, samples := range samplesList {
		for st, desc := range m.StatPromDesc {
			var point *Sample
			switch st {
			case "last":
				point, err = samples.GetLatestPoint()
				if err != nil {
					return nil, err
				}
			case "max":
				point, err = samples.GetMaxPoint()
				if err != nil {
					return nil, err
				}
			case "min":
				point, err = samples.GetMinPoint()
				if err != nil {
					return nil, err
				}
			case "avg":
				point, err = samples.GetAvgPoint()
				if err != nil {
					return nil, err
				}
			}

			var names []string
			var values []string

			labels := m.Labels.GetValues(samples.Series.QueryLabels, samples.Series.Instance)

			labels = map[string]string{
				"namespace":    m.Meta.Namespace,
				"region":       config.ExporterRunningRegion,
				"instancename": samples.Series.Instance.GetInstanceName(),
				"instanceid":   samples.Series.Instance.GetInstanceID(),
				"instanceip":   samples.Series.Instance.GetInstanceIP(),
			}

			for labelName, labelValue := range labels {
				names = append(names, util.ToUnderlineLower(labelName))
				values = append(values, labelValue)
			}

			if m.Meta.m.Dimensions != nil || len(*m.Meta.m.Dimensions) > 0 {
				for k, v := range *m.Meta.m.Dimensions {
					names = append(names, util.ToUnderlineLower(k))
					values = append(values, v)
				}
			}

			newDesc := prometheus.NewDesc(
				desc.FQName,
				desc.Help,
				names,
				nil,
			)

			var pm prometheus.Metric
			if m.Conf.StatDelaySeconds > 0 {
				pm = prometheus.NewMetricWithTimestamp(time.Unix(int64(point.Timestamp), 0), prometheus.MustNewConstMetric(
					newDesc,
					prometheus.GaugeValue,
					point.Value,
					values...,
				))
			} else {
				pm = prometheus.MustNewConstMetric(
					newDesc,
					prometheus.GaugeValue,
					point.Value,
					values...,
				)
			}
			pms = append(pms, pm)
		}
	}
	return
}

//LoadSeries
func (m *Metric) LoadSeries(series []*Series) error {
	m.seriesLock.Lock()
	defer m.seriesLock.Unlock()

	newSeriesCache := newCache()

	for _, s := range series {
		newSeriesCache.Series[s.Id] = s
	}
	m.SeriesCache = newSeriesCache
	return nil
}

func (m *Metric) GetSeriesSplitByBatch(batch int) (steps [][]*Series) {
	var series []*Series
	for _, s := range m.SeriesCache.Series {
		series = append(series, s)
	}

	total := len(series)
	for i := 0; i < total/batch+1; i++ {
		s := i * batch
		if s >= total {
			continue
		}
		e := i*batch + batch
		if e >= total {
			e = total
		}
		steps = append(steps, series[s:e])
	}
	return
}

// 创建Metric
func NewMetric(meta *Meta, conf *MetricConf) (*Metric, error) {
	id := fmt.Sprintf("%s.%s", meta.MetricName, *meta.m.InstanceID)
	labels, err := NewLabels(meta.SupportDimensions, conf.InstanceLabelNames, conf.ConstLabels)
	if err != nil {
		return nil, err
	}

	statPromDesc := make(map[string]Desc)

	statType := "avg"

	help := fmt.Sprintf("%s %s %s %s",
		meta.MetricName,
		meta.Namespace,
		*meta.m.MetricDesc,
		*meta.m.Unit,
	)

	for _, s := range conf.StatTypes {
		var st string
		if s == "last" {
			st = strings.ToLower(statType)
		} else {
			st = strings.ToLower(s)
		}

		displayMetricName := FilterByMetricName(meta.MetricName)
		// 显示的指标名称
		vmn := util.PointToUnderline(displayMetricName)
		vmn = util.MiddleToUnderline(vmn)

		fqName := fmt.Sprintf("%s_%s",
			vmn,
			st,
		)
		fqName = strings.ToLower(fqName)
		statPromDesc[strings.ToLower(s)] = Desc{
			FQName: fqName,
			Help:   help,
		}
	}

	m := &Metric{
		Id:           id,
		Meta:         meta,
		Labels:       labels,
		SeriesCache:  newCache(),
		StatPromDesc: statPromDesc,
		Conf:         conf,
	}
	return m, nil
}
