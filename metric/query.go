package metric

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/log"
)

// 负责一个指标的查询管理
type Query struct {
	Metric            *Metric
	LatestQueryStatus int
	repo              MetricRepository
}

type QuerySet []*Query

func (q *Query) GetPromMetrics() (pms []prometheus.Metric, err error) {

	q.LatestQueryStatus = 2

	pms, err = q.Metric.GetLatestPromMetrics(q.repo)
	if err != nil {
		return
	}

	q.LatestQueryStatus = 1
	return
}

func GetPromMetricsByQueries(queries []*Query, logger log.Logger) (pms []prometheus.Metric, err error) {
	if len(queries) <= 0 {
		return pms, fmt.Errorf("queries is empty")
	}

	queryMetrics := make(map[string]*Metric)
	for i := 0; i < len(queries); i++ {
		queryMetrics[queries[i].Metric.Id] = queries[i].Metric
	}

	pms, err = GetLatestPromMetrics(queries[0].repo, queryMetrics, logger)
	if err != nil {
		return
	}
	return
}

func (qs QuerySet) SplitByBatch(batch int) (steps [][]*Query) {
	total := len(qs)
	for i := 0; i < total/batch+1; i++ {
		s := i * batch
		if s >= total {
			continue
		}
		e := i*batch + batch
		if e >= total {
			e = total
		}
		steps = append(steps, qs[s:e])
	}
	return
}

func NewQuery(m *Metric, repo MetricRepository) (query *Query, err error) {
	query = &Query{
		Metric: m,
		repo:   repo,
	}
	return
}
