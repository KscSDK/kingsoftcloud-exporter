package instance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/slb"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("LISTENER7", NewInstanceListener7Repository)
}

//InstanceListenerHTTPSRepository
type InstanceListener7Repository struct {
	credential config.Credential
	client     *slb.Slb
	logger     log.Logger
}

func (repo *InstanceListener7Repository) GetNamespace() string {
	return "LISTENER7"
}

func (repo *InstanceListener7Repository) GetInstanceKey() string {
	return "LISTENER7"
}

func (repo *InstanceListener7Repository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceListener7Repository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceListener7Repository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	//24,26
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreHTTPInstances:

	instancesHTTP, instancesHTTPCount, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		24,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}
	for _, v := range instancesHTTP {
		meta := &InstanceListenerMeta{
			ListenerId:       v.InstanceID,
			ListenerName:     v.InstanceName,
			ListenerProtocol: "HTTP",
		}
		ins := &InstanceListener{
			InstanceBase: InstanceBase{
				InstanceID: v.InstanceID,
			},
			meta: meta,
		}
		instances = append(instances, ins)
	}

	if totalCount == -1 {
		totalCount = instancesHTTPCount
	}

	if (marker * maxResults) < totalCount {
		marker++
		goto getMoreHTTPInstances
	}

	level.Info(repo.logger).Log("msg", "HTTP 资源加载完毕")

	marker = 1

	maxResults = 300

	totalCount = -1

getMoreHTTPSInstances:
	instancesHTTPS, instancesHTTPSCount, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		27,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range instancesHTTPS {
		meta := &InstanceListenerMeta{
			ListenerId:       v.InstanceID,
			ListenerName:     v.InstanceName,
			ListenerProtocol: "HTTPS",
		}
		ins := &InstanceListener{
			InstanceBase: InstanceBase{
				InstanceID: v.InstanceID,
			},
			meta: meta,
		}
		instances = append(instances, ins)
	}

	if totalCount == -1 {
		totalCount = instancesHTTPSCount
	}

	if (marker * maxResults) < totalCount {
		marker++
		goto getMoreHTTPSInstances
	}

	level.Info(repo.logger).Log("msg", "HTTPS 资源加载完毕")

	return
}

func (repo *InstanceListener7Repository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	var nextToken int64 = 1

	var maxResults int64 = 300

	level.Info(repo.logger).Log("msg", "LISTENER7 资源开始加载")

getMoreInstances:

	filters["NextToken"] = nextToken
	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeListeners(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no listener.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeListenersResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse listener instance list err, %+v", err)
	}

	for _, v := range response.ListenerSet {
		if v.ListenerProtocol == "HTTPS" || v.ListenerProtocol == "HTTP" {
			instance, err := NewInstanceListener(v.ListenerId, v)
			if err != nil {
				level.Error(repo.logger).Log("msg", "get listener instance fail", "id", v.ListenerId)
				continue
			}
			instances = append(instances, instance)
		}
	}

	var responseNextToken int64 = 0
	if response.NextToken != "" || len(response.NextToken) > 0 {
		responseNextToken, _ = strconv.ParseInt(response.NextToken, 10, 64)
	}

	nextToken = responseNextToken
	if nextToken > 0 {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "LISTENER7 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceListener7Repository
func NewInstanceListener7Repository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := slb.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceListener7Repository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
