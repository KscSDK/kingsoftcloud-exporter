package metric

import "strconv"

// 代表一个数据点
type Sample struct {
	Timestamp int64
	Value     float64
}

// 代表一个时间线的多个数据点
type Samples struct {
	Series  *Series
	Samples []*Sample
}

func (s *Samples) GetLatestPoint() (point *Sample, err error) {
	if len(s.Samples) == 1 {
		return s.Samples[0], nil
	} else {
		return s.Samples[len(s.Samples)-1], nil
	}
}

func (s *Samples) GetMaxPoint() (point *Sample, err error) {
	maxValue := s.Samples[0].Value
	var maxIdx int
	for idx, sample := range s.Samples {
		if sample.Value > maxValue {
			maxValue = sample.Value
			maxIdx = idx
		}
	}
	return s.Samples[maxIdx], nil
}

func (s *Samples) GetMinPoint() (point *Sample, err error) {
	minValue := s.Samples[0].Value
	var minIdx int
	for idx, sample := range s.Samples {
		if sample.Value < minValue {
			minValue = sample.Value
			minIdx = idx
		}
	}
	return s.Samples[minIdx], nil
}

func (s *Samples) GetAvgPoint() (point *Sample, err error) {
	var sum float64
	for _, sample := range s.Samples {
		sum = sum + sample.Value
	}
	avg := sum / float64(len(s.Samples))
	sample := &Sample{
		Timestamp: s.Samples[len(s.Samples)-1].Timestamp,
		Value:     avg,
	}
	return sample, nil
}

func SplitBySamplesBatch(l []*Samples, batch int) (steps [][]*Samples) {
	total := len(l)
	for i := 0; i < total/batch+1; i++ {
		s := i * batch
		if s >= total {
			continue
		}
		e := i*batch + batch
		if e >= total {
			e = total
		}
		steps = append(steps, l[s:e])
	}
	return
}

func NewSamples(series *Series, mSeries MonitorSeries) (s *Samples, err error) {
	s = &Samples{
		Series:  series,
		Samples: []*Sample{},
	}

	for i := 0; i < len(mSeries.Data.Points); i++ {
		value, err := strconv.ParseFloat(mSeries.Data.Points[i].Avg, 64)
		if err != nil {
			continue
		}

		unixTimestamp := mSeries.Data.Points[i].UnixTimestamp / 1000

		s.Samples = append(s.Samples, &Sample{
			Timestamp: unixTimestamp,
			Value:     value,
		})
	}
	return
}
