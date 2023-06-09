package metric

import (
	"fmt"
	"strings"
	"sync"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
)

// 代表一个云监控指标的元数据
type Meta struct {
	Id                string
	Namespace         string
	ProductName       string
	MetricName        string
	MetricReName      string
	SupportDimensions []string
	m                 *MetricSet
}

func FilterByMetricName(metricName string) string {
	if len(metricName) <= 0 {
		return metricName
	}

	left := strings.Index(metricName, "[")

	if left <= -1 {
		return metricName
	}

	right := strings.Index(metricName, "]")
	if right <= -1 {
		return metricName
	}

	return metricName[:left]
}

func ParseDimensionsByName(metricName string) (string, []string) {
	if len(metricName) <= 0 {
		return metricName, nil
	}

	left := strings.Index(metricName, "[")

	if left <= -1 {
		return metricName, nil
	}

	right := strings.Index(metricName, "]")
	if right <= -1 {
		return metricName, nil
	}

	if left >= right {
		return metricName, nil
	}

	dimensionsStr := metricName[left+1 : right]

	dimensions := strings.Split(dimensionsStr, ",")

	return metricName[:left], dimensions
}

func NewMultiDimensionMeta(conf *config.KscExporterConfig, m *MetricSet) (*Meta, error) {

	metricName, dimensionValues := ParseDimensionsByName(*m.MetricName)

	id := fmt.Sprintf("%s-%s", *m.Namespace, *m.MetricName)

	var lock sync.Mutex

	lock.Lock()
	metricConfigMap := conf.GetMetricConfigMap(*m.Namespace)
	lock.Unlock()

	labels := make(map[string]string)

	metricReName := ""
	if len(dimensionValues) > 0 {
		//根据配置的查找对应的labels
		if _, isExist := metricConfigMap[*m.MetricName]; isExist {
			metricDesc := metricConfigMap[*m.MetricName].MetricDesc
			if len(metricDesc) >= 0 {
				m.MetricDesc = &metricDesc
			}

			for i := 0; i < len(metricConfigMap[*m.MetricName].Labels); i++ {
				if i < len(dimensionValues) {
					labels[metricConfigMap[*m.MetricName].Labels[i]] = dimensionValues[i]
				}
			}
		}

		//根据配置的查找对应的labels
		if _, isExist := metricConfigMap[metricName]; isExist {
			metricDesc := metricConfigMap[metricName].MetricDesc
			if len(metricDesc) >= 0 {
				m.MetricDesc = &metricDesc
			}

			for i := 0; i < len(metricConfigMap[metricName].Labels); i++ {
				if i < len(dimensionValues) {
					labels[metricConfigMap[metricName].Labels[i]] = dimensionValues[i]
				}
			}

			if metricName == "net.if.in" && len(dimensionValues) >= 2 {
				metricReName = "net.if.in.pps"
				*m.MetricDesc = "网卡入包数"
				*m.Unit = "pps"
			}

			if metricName == "net.if.out" && len(dimensionValues) >= 2 {
				metricReName = "net.if.out.pps"
				*m.MetricDesc = "网卡出包数"
				*m.Unit = "pps"
			}
		}

		if *m.MetricName == "vm.memory.size[total]" {
			metricReName = "vm.memory.total.size"
			*m.MetricDesc = "总内存"
			*m.Unit = "B"
		}

		if *m.MetricName == "vm.memory.size[available]" {
			metricReName = "vm.memory.available.size"
			*m.MetricDesc = "可用内存"
			*m.Unit = "B"
		}

		if *m.MetricName == "vm.memory.size[used]" {
			metricReName = "vm.memory.used.size"
			*m.MetricDesc = "已用内存"
			*m.Unit = "B"
		}

		if *m.MetricName == "vm.memory.size[pavailable]" {
			metricReName = "vm.memory.available.utilization"
			*m.MetricDesc = "内存可用率"
			*m.Unit = "%"
		}

	}

	m.Dimensions = &labels

	if m.MetricDesc == nil || len(*m.MetricDesc) <= 0 {
		m.MetricDesc = &metricName
	}

	var supportDimensions []string

	meta := &Meta{
		Id:                id,
		Namespace:         *m.Namespace,
		ProductName:       *m.Namespace,
		MetricName:        *m.MetricName,
		MetricReName:      metricReName,
		SupportDimensions: supportDimensions,
		m:                 m,
	}

	return meta, nil
}

func NewMeta(m *MetricSet) (*Meta, error) {
	id := fmt.Sprintf("%s-%s", *m.Namespace, *m.MetricName)

	var supportDimensions []string

	labels := make(map[string]string)

	m.Dimensions = &labels

	if len(*m.MetricDesc) <= 0 {
		*m.MetricDesc = *m.MetricName
	}

	meta := &Meta{
		Id:                id,
		Namespace:         *m.Namespace,
		ProductName:       *m.Namespace,
		MetricName:        *m.MetricName,
		SupportDimensions: supportDimensions,
		m:                 m,
	}
	return meta, nil

}
