package config

var AllEIPMetricConfigs = []KscMetricConfig{
	{
		Namespace:        "EIP",
		MetricName:       "eip.bps.in",
		MetricDesc:       "弹性IP入网流量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "bps",
	},
	{
		Namespace:        "EIP",
		MetricName:       "eip.bps.out",
		MetricDesc:       "弹性IP出网流量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "bps",
	},
	{
		Namespace:        "EIP",
		MetricName:       "eip.pps.in",
		MetricDesc:       "弹性IP每秒流入包数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "count",
	},
	{
		Namespace:        "EIP",
		MetricName:       "eip.pps.out",
		MetricDesc:       "弹性IP每秒流出包数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "count",
	},
	{
		Namespace:        "EIP",
		MetricName:       "eip.utilization.in",
		MetricDesc:       "弹性IP入向带宽使用百分比",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "%",
	},
	{
		Namespace:        "EIP",
		MetricName:       "eip.utilization.out",
		MetricDesc:       "弹性IP出向带宽使用百分比",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "%",
	},
}
