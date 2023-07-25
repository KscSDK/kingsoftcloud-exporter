package metric

type DataPoint struct {
	InstanceId string
	MetricName string
	Points     []Point
	Dimensions []Dimension
}

type MonitorSeries struct {
	InstanceId string      `json:"Instance"`
	Data       MonitorData `json:"datapoints"`
	Label      string      `json:"label"`
}

type MonitorData struct {
	Points []Point `json:"member"`
}

type Point struct {
	Timestamp     string `json:"timestamp"`
	UnixTimestamp int64  `json:"unixTimestamp"`
	Avg           string `json:"average"`
	Max           string `json:"max"`
	Min           string `json:"min"`
}

type GetMetricStatisticsBatchResponseV2 struct {
	Result       []MonitorSeries  `json:"getMetricStatisticsBatchResults"`
	Metadata     ResponseMetadata `json:"responseMetadata"`
	ErrorMessage []string         `json:"errorMessage"`
}

type GetMetricStatisticsBatchResponseV5 struct {
	GetMetricStatisticsBatchResults []GetMetricStatisticsResult `json:"getMetricStatisticsBatchResults"`
	ResponseMetadata                ResponseMetadata            `json:"responseMetadata"`
	ErrorMessage                    []string                    `json:"errorMessage,omitempty"`
}

type GetMetricStatisticsResult struct {
	Instance   string            `json:"instance"`
	MetricName string            `json:"metricName"`
	Points     []MetricDataPoint `json:"points"`
}

type MetricDataPoint struct {
	Dimensions []Dimension `json:"dimensions"`
	Values     []Point     `json:"values"`
}

type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
