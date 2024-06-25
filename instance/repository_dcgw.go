package instance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("DCGW", NewInstanceDCGWRepository)
}

//InstanceDCGWRepository
type InstanceDCGWRepository struct {
	credential config.Credential
	client     *vpc.Vpc
	logger     log.Logger
}

func (repo *InstanceDCGWRepository) GetNamespace() string {
	return "DCGW"
}

func (repo *InstanceDCGWRepository) GetInstanceKey() string {
	return "DCGW"
}

func (repo *InstanceDCGWRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceDCGWRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceDCGWRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "DC_GW 资源开始加载")
getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		17,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceDCGWMeta{
			DirectConnectGatewayId:   v.InstanceID,
			DirectConnectGatewayName: v.InstanceName,
		}
		ins := &InstanceDCGW{
			InstanceBase: InstanceBase{
				InstanceID: v.InstanceID,
			},
			meta: meta,
		}
		instances = append(instances, ins)
	}

	if totalCount == -1 {
		totalCount = count
	}

	if (marker * maxResults) < totalCount {
		marker++
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "DC_GW 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeDirectConnectGatewaysResponse struct {
	DirectConnectGatewaySet []*InstanceDCGWMeta `json:"DirectConnectGatewaySet"`
	NextToken               string              `json:"NextToken"`
	RequestId               string              `json:"RequestId"`
}

func (repo *InstanceDCGWRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	var nextToken int64 = 1

	var maxResults int64 = 300

	level.Info(repo.logger).Log("msg", "DC_GW 资源开始加载")

getMoreInstances:

	filters["NextToken"] = nextToken
	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeDirectConnectGateways(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no listener.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeDirectConnectGatewaysResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse dcgw instance list err, %+v", err)
	}

	for _, v := range response.DirectConnectGatewaySet {
		instance, err := NewInstanceDCGW(v.DirectConnectGatewayId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get dcgw instance fail", "id", v.DirectConnectGatewayId)
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

	level.Info(repo.logger).Log("msg", "DC_GW 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceDCGWRepository
func NewInstanceDCGWRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := vpc.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	repo := &InstanceDCGWRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
