package metric

import (
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/instance"
)

// 代表某个指标的一个时间线
type Series struct {
	Id          string
	Metric      *Metric
	QueryLabels Labels
	Instance    instance.KscInstance
}

func GetSeriesId(m *Metric, instanceId string) (string, error) {
	return fmt.Sprintf("%s", m.Id), nil
}

func NewSeries(m *Metric, ql Labels, ins instance.KscInstance) (*Series, error) {
	id, err := GetSeriesId(m, ins.GetInstanceID())
	if err != nil {
		return nil, err
	}

	s := &Series{
		Id:          id,
		Metric:      m,
		QueryLabels: ql,
		Instance:    ins,
	}
	return s, nil

}
