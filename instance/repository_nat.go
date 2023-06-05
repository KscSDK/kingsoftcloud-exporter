package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("NAT", NewInstanceNATRepository)
}

//InstanceNATRepository
type InstanceNATRepository struct {
	credential config.Credential
	client     *vpc.Vpc
	logger     log.Logger
}

func (repo *InstanceNATRepository) GetInstanceKey() string {
	return "NAT"
}

func (repo *InstanceNATRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceNATRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceNATRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.MockInnerURL,
		repo.credential.MockAccountId,
		10,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceNATMeta{
			NatId:   v.InstanceID,
			NatName: v.InstanceName,
		}
		ins := &InstanceNAT{
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

	level.Info(repo.logger).Log("msg", "NAT 资源加载完毕")

	return
}

type DescribeNatSetResponse struct {
	NatSet    []*InstanceNATMeta `json:"NatSet"`
	NextToken int64              `json:"NextToken"`
	RequestId string             `json:"RequestId"`
}

func (repo *InstanceNATRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	// var nextToken int64 = 0

	// var maxResults int64 = 10

	level.Info(repo.logger).Log("msg", "NAT 资源开始加载")

	// getMoreInstances:

	// 	filters["Marker"] = nextToken
	// 	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeNats(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no nat resources.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeNatSetResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse nat instance list err, %+v", err)
	}

	for _, v := range response.NatSet {
		instance, err := NewInstanceNAT(v.NatId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get nat instance fail", "id", v.NatId)
			continue
		}
		instances = append(instances, instance)
	}

	// nextToken = response.NextToken
	// if nextToken > 0 {
	// 	goto getMoreInstances
	// }
	level.Info(repo.logger).Log("msg", "NAT 资源加载完毕")

	return nil, nil
}

//NewInstanceSLBRepository
func NewInstanceNATRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := vpc.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceNATRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
