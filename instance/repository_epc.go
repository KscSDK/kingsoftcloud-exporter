package instance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/epc"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("EPC", NewInstanceEPCRepository)
}

//InstanceEPCRepository
type InstanceEPCRepository struct {
	credential config.Credential
	client     *epc.Epc
	logger     log.Logger
}

func (repo *InstanceEPCRepository) GetInstanceKey() string {
	return "EPC"
}

func (repo *InstanceEPCRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceEPCRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceEPCRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	return nil, nil
}

type DescribeEPCsResponse struct {
	HostSet    []*InstanceEPCMeta `json:"HostSet"`
	NextToken  string             `json:"NextToken"`
	TotalCount int64              `json:"TotalCount"`
	RequestId  string             `json:"RequestId"`
}

func (repo *InstanceEPCRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	var nextToken int64 = 0

	var maxResults int64 = 300

	level.Info(repo.logger).Log("msg", "EPC 资源开始加载")

getMoreInstances:

	filters["NextToken"] = nextToken
	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeEpcs(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no epc.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeEPCsResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse listener instance list err, %+v", err)
	}

	for _, v := range response.HostSet {
		instance, err := NewInstanceEPC(v.HostId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get epc instance fail", "id", v.HostId)
			continue
		}
		instances = append(instances, instance)
	}

	var responseNextToken int64 = 0
	if response.NextToken != "" || len(response.NextToken) > 0 {
		responseNextToken, _ = strconv.ParseInt(response.NextToken, 10, 64)
	}

	nextToken = responseNextToken
	if nextToken > 0 {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "EPC 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceEPCRepository
func NewInstanceEPCRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := epc.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceEPCRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
