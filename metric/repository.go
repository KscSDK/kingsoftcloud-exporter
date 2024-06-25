package metric

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/monitor"
	v2 "github.com/KscSDK/ksc-sdk-go/service/monitorv2"
	v3 "github.com/KscSDK/ksc-sdk-go/service/monitorv3"
	v5 "github.com/KscSDK/ksc-sdk-go/service/monitorv5"
	"github.com/google/uuid"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/time/rate"
)

//MetricRepository 云监控指标 Repository
type MetricRepository interface {
	//ListMetrics 获取多维度指标
	ListMetrics(namespace, instanceId string) ([]*Meta, error)

	//ListMetrics 从本地配置文件中获取指标
	ListLocalMetrics(namespace, instanceId string) ([]*Meta, error)

	// 获取指标的元数据
	GetMeta(conf config.KscMetricConfig, instanceId string) (*Meta, error)

	// 按时间范围批量获取数据点
	ListBatchSamples(namespace string, metric map[string]*Metric, period int64, startTime int64, endTime int64) (metricSamples map[string][]*Samples, err error)

	// 按时间范围获取单个指标下所有时间线的数据点
	ListSamples(metric *Metric, startTime int64, endTime int64) (samplesList []*Samples, err error)

	DescribeMonitorData(namespace string, metric map[string]*Metric, period int64, startTime int64, endTime int64) (metricSamples map[string][]*Samples, err error)
}

//MetricRepositoryImpl
type MetricRepositoryImpl struct {
	exporterConf *config.KscExporterConfig

	monitorClient *monitor.Monitor

	monitorClientV2 *v2.Monitorv2

	monitorClientV3 *v3.Monitorv3

	monitorClientV5 *v5.Monitorv5

	limiter *rate.Limiter // 限速

	ctx context.Context

	queryMetricBatchSize int

	logger log.Logger
}

func (repo *MetricRepositoryImpl) GetMeta(conf config.KscMetricConfig, instanceId string) (meta *Meta, err error) {
	interval := strconv.FormatInt(conf.MinPeriodSeconds, 10)

	m := &MetricSet{
		Namespace:  &conf.Namespace,
		MetricName: &conf.MetricName,
		MetricDesc: &conf.MetricDesc,
		Statistics: &conf.Statistics,
		InstanceID: &instanceId,
		Interval:   &interval,
		Unit:       &conf.Unit,
	}

	meta, err = NewMeta(m)
	if err != nil {
		return
	}

	return
}

func (repo *MetricRepositoryImpl) ListLocalMetrics(namespace, instanceId string) ([]*Meta, error) {

	metricsConf, err := config.GetMetricConfigs(namespace)
	if err != nil {
		return nil, err
	}

	metaSlice := make([]*Meta, 0, len(metricsConf))

	for i := 0; i < len(metricsConf); i++ {
		m, err := repo.GetMeta(metricsConf[i], instanceId)
		if err != nil {
			return nil, err
		}
		metaSlice = append(metaSlice, m)
	}

	return metaSlice, nil
}

//ListMetrics
func (repo *MetricRepositoryImpl) ListMetrics(namespace, instanceId string) ([]*Meta, error) {
	// 限速
	ctx, cancel := context.WithCancel(repo.ctx)
	defer cancel()

	if err := repo.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	var metricSets []*MetricSet
	var err error
	if config.ExporterRunningMode == config.ExporterMode_Mock {
		metricSets, err = repo.describeMetricsMetaRequest(namespace, instanceId)
	} else {
		metricSets, err = repo.listMetricsRequest(namespace, instanceId)
	}
	if err != nil {
		return nil, err
	}

	onlyIncludeMetricsMaps := config.GetOnlyIncludeMetrics(namespace)
	metaSlice := make([]*Meta, 0, len(metricSets))
	for _, metricSet := range metricSets {
		filterName := FilterByMetricName(*metricSet.MetricName)
		if len(onlyIncludeMetricsMaps) > 0 {
			if _, isOK := onlyIncludeMetricsMaps[filterName]; !isOK {
				continue
			}
		}
		m, e := NewMultiDimensionMeta(repo.exporterConf, metricSet)
		if e != nil {
			return nil, err
		}
		metaSlice = append(metaSlice, m)
	}
	return metaSlice, nil
}

//ListMetricsResponse
type ListMetricsResponse struct {
	List     ListMetricsResult `json:"listMetricsResult"`
	Metadata ResponseMetadata  `json:"responseMetadata"`
}

type ListMetricsResult struct {
	Metrics Member `json:"metrics"`
}

type ResponseMetadata struct {
	RequestId string `json:"requestId"`
}

type Member struct {
	Member []*MetricSet `json:"member"`
}

//ListMetricsRequest
func (repo *MetricRepositoryImpl) listMetricsRequest(namespace, instanceId string) ([]*MetricSet, error) {
	requestParams := make(map[string]interface{})
	requestParams["Namespace"] = namespace
	requestParams["InstanceID"] = instanceId
	requestParams["PageIndex"] = "1"

	var resp *map[string]interface{}
	var err error
	if namespace == "KCM" {
		resp, err = repo.monitorClientV5.ListMetrics(&requestParams)
	} else {
		resp, err = repo.monitorClient.ListMetrics(&requestParams)
	}

	if err != nil {
		return nil, err
	}

	respBytes, _ := json.Marshal(resp)

	var response ListMetricsResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, err
	}

	if len(response.List.Metrics.Member) <= 0 {
		return nil, fmt.Errorf("response metricSet size <= 0")
	}

	return response.List.Metrics.Member, nil
}

//describeMetricsMetaRequest
func (repo *MetricRepositoryImpl) describeMetricsMetaRequest(namespace, instanceId string) ([]*MetricSet, error) {
	if len(repo.exporterConf.Credential.AccessMetricMetaURL) <= 0 {
		return nil, fmt.Errorf("mock inner url is empty.")
	}

	apiURL := fmt.Sprintf("%s&InstanceID=%s&Namespace=%s&PageIndex=1",
		repo.exporterConf.Credential.AccessMetricMetaURL,
		instanceId,
		namespace,
	)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	requestID := uuid.New().String()

	req.Header = http.Header{
		textproto.CanonicalMIMEHeaderKey("Content-Type"):     []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("Accept"):           []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("X-KSC-ACCOUNT-ID"): []string{repo.exporterConf.Credential.AccessAccount},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Region"):     []string{repo.exporterConf.Credential.Region},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Request-Id"): []string{requestID},
	}

	c := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 100,
		},
		Timeout: 60 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var response ListMetricsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse monitor data err, %+v", err)
	}

	if len(response.List.Metrics.Member) <= 0 {
		return nil, fmt.Errorf("response metricSet size <= 0")
	}

	return response.List.Metrics.Member, nil
}

func (repo *MetricRepositoryImpl) ListSamples(m *Metric, st int64, et int64) ([]*Samples, error) {
	return nil, nil
}

func (repo *MetricRepositoryImpl) ListBatchSamples(namespace string, m map[string]*Metric, period int64, st int64, et int64) (map[string][]*Samples, error) {
	return repo.getMetricStatisticsBatch(namespace, m, period, st, et)
}

func (repo *MetricRepositoryImpl) GetMetricStatisticsBatchRequest(m map[string]*Metric,
	st int64,
	et int64,
) ([]*Samples, error) {
	return nil, nil
}

//GetMetricStatisticsBatch
func (repo *MetricRepositoryImpl) getMetricStatisticsBatch(
	namespace string,
	ms map[string]*Metric,
	period int64,
	st int64,
	et int64,
) (map[string][]*Samples, error) {

	ctx, cancel := context.WithCancel(repo.ctx)
	defer cancel()

	err := repo.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	requestParams := repo.buildGetMonitorRequest(ms, period, st, et)

	var dataPoints []*DataPoint
	if namespace == "KCM" {
		resp, err := repo.getMetricStatisticsBatchV5(requestParams)
		if err != nil {
			return nil, err
		}

		for _, result := range resp.GetMetricStatisticsBatchResults {
			for _, point := range result.Points {
				dataPoints = append(dataPoints, &DataPoint{
					InstanceId: result.Instance,
					MetricName: result.MetricName,
					Points:     point.Values,
					Dimensions: point.Dimensions,
				})
			}
		}
	} else {
		resp, err := repo.getMetricStatisticsBatchV2(requestParams)
		if err != nil {
			return nil, err
		}

		for _, point := range resp.Result {
			dataPoint := &DataPoint{
				InstanceId: point.InstanceId,
				MetricName: point.Label,
				Points:     point.Data.Points,
			}

			if namespace == "MONGO" {
				dataPoint.InstanceId = dataPoint.InstanceId[4:]
			}

			dataPoints = append(dataPoints, dataPoint)
		}
	}

	metricSamplesList := make(map[string][]*Samples)
	for _, point := range dataPoints {
		id := fmt.Sprintf("%s.%s", point.MetricName, point.InstanceId)
		if _, isExist := ms[id]; isExist {
			samples, ql, e := repo.buildSamples(ms[id], point)
			if e != nil {
				level.Error(repo.logger).Log(
					"msg", e.Error(),
					"metric", ms[id].Meta.MetricName,
					"dimension", fmt.Sprintf("%v", ql))
				continue
			}
			metricSamplesList[id] = append(metricSamplesList[id], samples)
		}
	}

	return metricSamplesList, nil
}

type RequestMonitorMetric struct {
	InstanceID string `json:"InstanceID"`
	MetricName string `json:"MetricName"`
}

func (repo *MetricRepositoryImpl) buildGetMonitorRequest(
	ms map[string]*Metric,
	period, st, et int64,
) map[string]interface{} {
	requestParams := make(map[string]interface{})

	requestParams["StartTime"] = time.Unix(st, 0).Format("2006-01-02T15:04:05Z")
	requestParams["EndTime"] = time.Unix(et, 0).Format("2006-01-02T15:04:05Z")
	requestParams["Period"] = period
	requestParams["Aggregate"] = []string{"Max"}

	requestMetrics := make([]*RequestMonitorMetric, 0, len(ms))

	for k, v := range ms {

		requestMetric := &RequestMonitorMetric{
			InstanceID: *ms[k].Meta.m.InstanceID,
			MetricName: *&ms[k].Meta.MetricName,
		}

		requestParams["Namespace"] = *v.Meta.m.Namespace

		if *v.Meta.m.Namespace == "LISTENER7" {
			requestParams["Namespace"] = "LISTENER"
		}

		if *v.Meta.m.Namespace == "KCS" {
			requestParams["Namespace"] = "KCS2"
		}

		if *v.Meta.m.Namespace == "MONGO" {
			requestParams["Namespace"] = "MONDB"
			requestMetric.InstanceID = fmt.Sprintf("user%+v", requestMetric.InstanceID)
		}

		requestMetrics = append(requestMetrics, requestMetric)
	}

	requestParams["Metrics"] = requestMetrics

	return requestParams
}

func (repo *MetricRepositoryImpl) getMetricStatisticsBatchV2(requestParams map[string]interface{}) (*GetMetricStatisticsBatchResponseV2, error) {
	start := time.Now()
	resp, err := repo.monitorClientV2.GetMetricStatisticsBatch(&requestParams)
	if err != nil {
		level.Error(repo.logger).Log("request start time ", requestParams["StartTime"], "duration ", time.Since(start).Seconds(), "err ", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("get monitor data is empty")
	}

	respBytes, _ := json.Marshal(&resp)
	var rep GetMetricStatisticsBatchResponseV2
	if err := json.Unmarshal(respBytes, &rep); err != nil {
		return nil, err
	}

	if len(rep.ErrorMessage) > 0 {
		for _, v := range rep.ErrorMessage {
			level.Debug(repo.logger).Log("msg", v)
		}
	}

	return &rep, nil
}

func (repo *MetricRepositoryImpl) getMetricStatisticsBatchV5(requestParams map[string]interface{}) (*GetMetricStatisticsBatchResponseV5, error) {
	start := time.Now()
	resp, err := repo.monitorClientV5.GetMetricStatisticsBatch(&requestParams)
	if err != nil {
		level.Error(repo.logger).Log("request start time ", requestParams["StartTime"], "duration ", time.Since(start).Seconds(), "err ", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("get monitor data is empty")
	}

	respBytes, _ := json.Marshal(&resp)
	var rep GetMetricStatisticsBatchResponseV5
	if err := json.Unmarshal(respBytes, &rep); err != nil {
		return nil, err
	}

	if len(rep.ErrorMessage) > 0 {
		for _, v := range rep.ErrorMessage {
			level.Debug(repo.logger).Log("msg", v)
		}
	}

	return &rep, nil
}

func (repo *MetricRepositoryImpl) buildSamples(
	m *Metric,
	p *DataPoint,
) (*Samples, map[string]string, error) {

	ql := map[string]string{}

	sid, e := GetSeriesId(m, p.InstanceId)
	if e != nil {
		return nil, ql, fmt.Errorf("get series id fail")
	}
	s, ok := m.SeriesCache.Series[sid]
	if !ok {
		return nil, ql, fmt.Errorf("response data point not match series")
	}
	samples, e := NewSamples(s, p)
	if e != nil {
		return nil, ql, fmt.Errorf("this instance may not have metric data")
	}
	return samples, ql, nil
}

//DescribeMonitorData
func (repo *MetricRepositoryImpl) DescribeMonitorData(namespace string, m map[string]*Metric, period int64, st int64, et int64) (map[string][]*Samples, error) {
	ctx, cancel := context.WithCancel(repo.ctx)
	defer cancel()

	err := repo.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	requestParams := repo.buildGetMonitorRequest(m, period, st, et)

	start := time.Now()
	points, err := repo.describeMonitorDataRequest(namespace, &requestParams)
	if err != nil {
		level.Error(repo.logger).Log("request start time ", requestParams["StartTime"], "duration ", time.Since(start).Seconds(), "err ", err.Error())
		return nil, err
	}

	metricSamplesList := make(map[string][]*Samples)
	for _, point := range points {

		instanceId := point.InstanceId

		if namespace == "MONGO" && len(instanceId) > 36 {
			instanceId = instanceId[len(instanceId)-36:]
		}

		id := fmt.Sprintf("%s.%s", point.MetricName, instanceId)
		if _, isExist := m[id]; isExist {
			samples, ql, e := repo.buildSamples(m[id], point)
			if e != nil {
				level.Error(repo.logger).Log(
					"msg", e.Error(),
					"metric", m[id].Meta.MetricName,
					"dimension", fmt.Sprintf("%v", ql))
				continue
			}
			metricSamplesList[id] = append(metricSamplesList[id], samples)
		}
	}

	return metricSamplesList, nil
}

func (repo *MetricRepositoryImpl) describeMonitorDataRequest(namespace string, input *map[string]interface{}) ([]*DataPoint, error) {
	if len(repo.exporterConf.Credential.AccessMonitorMetaURL) <= 0 {
		return nil, fmt.Errorf("mock inner url is empty")
	}

	dataJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, repo.exporterConf.Credential.AccessMonitorMetaURL, bytes.NewBuffer(dataJSON))
	if err != nil {
		return nil, err
	}

	requestID := uuid.New().String()

	req.Header = http.Header{
		textproto.CanonicalMIMEHeaderKey("Content-Type"):     []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("Accept"):           []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("X-KSC-ACCOUNT-ID"): []string{repo.exporterConf.Credential.AccessAccount},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Region"):     []string{repo.exporterConf.Credential.Region},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Request-Id"): []string{requestID},
	}

	c := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 100,
		},
		Timeout: 60 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var dataPoints []*DataPoint
	if namespace == "KCM" {
		var response GetMetricStatisticsBatchResponseV5
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("parse monitor data err, %+v", err)
		}

		for _, result := range response.GetMetricStatisticsBatchResults {
			for _, point := range result.Points {
				dataPoints = append(dataPoints, &DataPoint{
					InstanceId: result.Instance,
					MetricName: result.MetricName,
					Points:     point.Values,
					Dimensions: point.Dimensions,
				})
			}
		}
	} else {
		var response GetMetricStatisticsBatchResponseV2
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("parse monitor data err, %+v", err)
		}

		for _, v := range response.ErrorMessage {
			level.Debug(repo.logger).Log("msg", v)
		}

		for _, point := range response.Result {
			dataPoints = append(dataPoints, &DataPoint{
				InstanceId: point.InstanceId,
				MetricName: point.Label,
				Points:     point.Data.Points,
			})
		}
	}

	return dataPoints, nil
}

//NewMetricRepository
func NewMetricRepository(
	conf *config.KscExporterConfig,
	logger log.Logger,
) (repo MetricRepository, err error) {

	client := monitor.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	clientV2 := v2.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	clientV3 := v3.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	clientV5 := v5.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	repo = &MetricRepositoryImpl{
		exporterConf:         conf,
		monitorClient:        client,
		monitorClientV2:      clientV2,
		monitorClientV3:      clientV3,
		monitorClientV5:      clientV5,
		limiter:              rate.NewLimiter(rate.Limit(conf.RateLimit), 1),
		ctx:                  context.Background(),
		queryMetricBatchSize: conf.MetricQueryBatchSize,
		logger:               logger,
	}

	return
}
