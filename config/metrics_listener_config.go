package config

var AllLISTENERMetricConfigs = []KscMetricConfig{
	{
		Namespace:        "LISTENER",
		MetricName:       "listener.bps.in",
		MetricDesc:       "监听器入网流量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "bps",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.bps.out",
		MetricDesc:       "监听器出网流量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "bps",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.pps.in",
		MetricDesc:       "监听器每秒流入包数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.pps.out",
		MetricDesc:       "监听器每秒流出包数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.cps",
		MetricDesc:       "监听器每秒新建连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.activeconn",
		MetricDesc:       "监听器当前活跃连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.inactiveconn",
		MetricDesc:       "监听器当前未活跃连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER",
		MetricName:       "listener.concurrentconn",
		MetricDesc:       "监听器并发连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	},
}

var AllLISTENER7MetricConfigs = []KscMetricConfig{
	{
		Namespace:        "LISTENER7",
		MetricName:       "listener.bps.in",
		MetricDesc:       "监听器入网流量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "bps",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.bps.out",
		MetricDesc:       "监听器出网流量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "bps",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.pps.in",
		MetricDesc:       "监听器每秒流入包数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.pps.out",
		MetricDesc:       "监听器每秒流出包数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.cps",
		MetricDesc:       "监听器每秒新建连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.activeconn",
		MetricDesc:       "监听器当前活跃连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.inactiveconn",
		MetricDesc:       "监听器当前未活跃连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.2xx",
		MetricDesc:       "监听器返回的2XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.3xx",
		MetricDesc:       "监听器返回的3XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.4xx",
		MetricDesc:       "监听器返回的4XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.5xx",
		MetricDesc:       "监听器返回的5XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.backend.2xx",
		MetricDesc:       "真实服务器返回的2XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.backend.3xx",
		MetricDesc:       "真实服务器返回的3XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.backend.4XX",
		MetricDesc:       "真实服务器返回的4XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.httpcode.backend.5xx",
		MetricDesc:       "真实服务器返回的5XX的状态码数量",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.latency",
		MetricDesc:       "HTTP请求到后端的延时，Average（单位时间内平均值）",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "ms",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.requestcount",
		MetricDesc:       "单位时间内完成的HTTP请求数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	}, {
		Namespace:        "LISTENER7",
		MetricName:       "listener.concurrentconn",
		MetricDesc:       "监听器并发连接数",
		MetricType:       1,
		Labels:           []string{},
		Statistics:       []string{"avg"},
		MinPeriodSeconds: 60,
		Unit:             "个",
	},
}
