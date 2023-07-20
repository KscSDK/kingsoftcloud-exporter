package metric

import (
	"strings"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
)

type MetricConf struct {
	CustomProductName    string
	CustomMetricName     string
	MetricNameType       int32
	InstanceLabelNames   []string
	ConstLabels          map[string]string
	StatTypes            []string
	StatPeriodSeconds    int64
	StatNumSamples       int64
	StatDelaySeconds     int64
	AllInstances         bool
	InstanceFilters      map[string]interface{}
	OnlyIncludeInstances []string
	ExcludeInstances     []string
}

func (c *MetricConf) IsIncludeOnlyInstance() bool {
	return len(c.OnlyIncludeInstances) > 0
}

func (c *MetricConf) IsIncludeAllInstance() bool {
	if c.IsIncludeOnlyInstance() {
		return false
	}
	return c.AllInstances
}

func NewMetricConfigWithMetricYaml(c config.KscMetricConfig, meta *Meta) (*MetricConf, error) {
	conf := &MetricConf{}

	conf.CustomMetricName = c.MetricName

	conf.InstanceLabelNames = c.Labels
	if len(c.Statistics) != 0 {
		for _, st := range c.Statistics {
			conf.StatTypes = append(conf.StatTypes, strings.ToLower(st))
		}
	} else {
		conf.StatTypes = []string{"last"}
	}
	// 自动获取支持的统计周期
	// period, err := meta.GetPeriod(c.PeriodSeconds)
	// if err != nil {
	// 	return nil, err
	// }
	period := int64(config.DefaultPeriodSeconds)
	conf.StatPeriodSeconds = period
	conf.StatNumSamples = (c.RangeSeconds / period) + 1
	// 至少采集4个点的数据
	if conf.StatNumSamples < 4 {
		conf.StatNumSamples = 4
	}
	conf.StatDelaySeconds = c.DelaySeconds

	if len(c.Dimensions) == 0 {
		conf.AllInstances = true
	}

	return conf, nil

}

func NewMetricConfigWithProductYaml(c config.KscProductConfig, meta *Meta) (*MetricConf, error) {
	conf := &MetricConf{}

	conf.CustomProductName = c.Namespace
	conf.CustomMetricName = ""
	if c.MetricNameType != 0 {
		conf.MetricNameType = c.MetricNameType
	} else {
		conf.MetricNameType = 2
	}

	conf.InstanceLabelNames = c.ExtraLabels
	if len(c.Statistics) != 0 {
		for _, st := range c.Statistics {
			conf.StatTypes = append(conf.StatTypes, strings.ToLower(st))
		}
	} else {
		conf.StatTypes = []string{"last"}
	}

	// period, err := meta.GetPeriod(c.PeriodSeconds)
	// if err != nil {
	// 	return nil, err
	// }
	period := int64(config.DefaultPeriodSeconds)
	conf.StatPeriodSeconds = period
	conf.StatNumSamples = (c.RangeSeconds / period) + 1
	if conf.StatNumSamples < 4 {
		conf.StatNumSamples = 4
	}
	conf.StatDelaySeconds = c.DelaySeconds
	conf.AllInstances = c.AllInstances
	conf.InstanceFilters = c.InstanceFilters
	conf.OnlyIncludeInstances = c.OnlyIncludeInstances
	conf.ExcludeInstances = c.ExcludeInstances

	return conf, nil
}
