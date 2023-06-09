package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("PEER", NewInstancePEERRepository)
}

//InstancePEERRepository
type InstancePEERRepository struct {
	credential config.Credential
	client     *vpc.Vpc
	logger     log.Logger
}

func (repo *InstancePEERRepository) GetInstanceKey() string {
	return "PEER"
}

func (repo *InstancePEERRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstancePEERRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstancePEERRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	//16
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.MockInnerURL,
		repo.credential.MockAccountId,
		16,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstancePEERMeta{
			VpcPeeringConnectionId: v.InstanceID,
			PeeringName:            v.InstanceName,
		}
		ins := &InstancePEER{
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

	level.Info(repo.logger).Log("msg", "PEER 资源加载完毕")

	return
}

type DescribeVpcPeeringConnectionsResponse struct {
	VpcPeeringConnectionSet []*InstancePEERMeta `json:"VpcPeeringConnectionSet"`
	RequestId               string              `json:"RequestId"`
}

func (repo *InstancePEERRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	var nextToken int64 = 1

	var maxResults int64 = 10

	level.Info(repo.logger).Log("msg", "Peering 资源开始加载")

	if len(iam.IAMProjectIDs) > 0 || len(iam.IAMProjectIDs) <= 100 {
		for i := 0; i < len(iam.IAMProjectIDs); i++ {
			filters[fmt.Sprintf("ProjectId.%d", i+1)] = iam.IAMProjectIDs[i]
		}
	}

getMoreInstances:

	filters["NextToken"] = nextToken
	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeVpcPeeringConnections(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no slb.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeVpcPeeringConnectionsResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse peering instance list err, %+v", err)
	}

	for _, v := range response.VpcPeeringConnectionSet {
		instance, err := NewInstancePEER(v.VpcPeeringConnectionId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get peering instance fail", "id", v.VpcPeeringConnectionId)
			continue
		}
		instances = append(instances, instance)
	}

	if len(response.VpcPeeringConnectionSet) == filters["MaxResults"] {
		nextToken += maxResults
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "Peering 资源加载完毕")

	return
}

//NewInstancePEERRepository
func NewInstancePEERRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := vpc.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstancePEERRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
