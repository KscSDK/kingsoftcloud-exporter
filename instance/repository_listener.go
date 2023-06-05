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
	registerRepository("LISTENER", NewInstanceListenerRepository)
}

//InstanceListenerRepository
type InstanceListenerRepository struct {
	credential config.Credential
	client     *slb.Slb
	logger     log.Logger
}

func (repo *InstanceListenerRepository) GetNamespace() string {
	return "LISTENER"
}

func (repo *InstanceListenerRepository) GetInstanceKey() string {
	return "LISTENER"
}

func (repo *InstanceListenerRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceListenerRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceListenerRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	//25,27
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreTCPInstances:

	instancesTCP, instancesTCPCount, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		25, marker, maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range instancesTCP {
		meta := &InstanceListenerMeta{
			ListenerId:       v.InstanceID,
			ListenerName:     v.InstanceName,
			ListenerProtocol: "TCP",
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
		totalCount = instancesTCPCount
	}

	if (marker * maxResults) < totalCount {
		marker++
		goto getMoreTCPInstances
	}

	level.Info(repo.logger).Log("msg", "TCP 资源加载完毕", "instance_num", len(instances))

	marker = 1

	maxResults = 300

	totalCount = -1

getMoreUDPInstances:
	instancesUDP, instancesUDPCount, err := DescribeMonitorInstances(
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

	for _, v := range instancesUDP {
		meta := &InstanceListenerMeta{
			ListenerId:       v.InstanceID,
			ListenerName:     v.InstanceName,
			ListenerProtocol: "UDP",
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
		totalCount = instancesUDPCount
	}

	if (marker * maxResults) < totalCount {
		marker++
		goto getMoreUDPInstances
	}

	level.Info(repo.logger).Log("msg", "UDP 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeListenersResponse struct {
	ListenerSet []*InstanceListenerMeta `json:"ListenerSet"`
	NextToken   string                  `json:"NextToken"`
	RequestId   string                  `json:"RequestId"`
}

func (repo *InstanceListenerRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	var nextToken int64 = 1

	var maxResults int64 = 300

	level.Info(repo.logger).Log("msg", "LISTENER 资源开始加载")

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
		if v.ListenerProtocol == "TCP" || v.ListenerProtocol == "UDP" {
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

	level.Info(repo.logger).Log("msg", "LISTENER 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceListenerRepository
func NewInstanceListenerRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := slb.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceListenerRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
