package metric

type Point struct {
	Timestamp     string `json:"timestamp"`
	UnixTimestamp int64  `json:"unixTimestamp"`
	Avg           string `json:"average"`
	Max           string `json:"max"`
	Min           string `json:"min"`
}

type MonitorData struct {
	Points []Point `json:"member"`
}

type MonitorSeries struct {
	InstanceId string      `json:"Instance"`
	Data       MonitorData `json:"datapoints"`
	Label      string      `json:"label"`
}

type GetMetricStatisticsBatchResponse struct {
	Result   []MonitorSeries  `json:"getMetricStatisticsBatchResults"`
	Metadata ResponseMetadata `json:"responseMetadata"`
}
