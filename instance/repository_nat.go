package instance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
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

	level.Info(repo.logger).Log("msg", "NAT 资源开始加载")

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
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

	level.Info(repo.logger).Log("msg", "NAT 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeNatSetResponse struct {
	NatSet    []*InstanceNATMeta `json:"NatSet"`
	NextToken string             `json:"NextToken"`
	RequestId string             `json:"RequestId"`
}

func (repo *InstanceNATRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	var nextToken int64 = 1

	var maxResults int64 = 300

	level.Info(repo.logger).Log("msg", "NAT 资源开始加载")

	if len(iam.IAMProjectIDs) > 0 || len(iam.IAMProjectIDs) <= 100 {
		for i := 0; i < len(iam.IAMProjectIDs); i++ {
			filters[fmt.Sprintf("ProjectId.%d", i+1)] = iam.IAMProjectIDs[i]
		}
	}

getMoreInstances:

	filters["NextToken"] = nextToken
	filters["MaxResults"] = maxResults

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

	var responseNextToken int64 = 0
	if response.NextToken != "" || len(response.NextToken) > 0 {
		responseNextToken, _ = strconv.ParseInt(response.NextToken, 10, 64)
	}

	nextToken = responseNextToken
	if nextToken > 0 {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "NAT 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceNATRepository
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
