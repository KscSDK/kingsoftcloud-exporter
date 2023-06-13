package instance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/slb"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("SLB", NewInstanceSLBRepository)
}

//InstanceSLBRepository
type InstanceSLBRepository struct {
	credential config.Credential
	client     *slb.Slb
	logger     log.Logger
}

func (repo *InstanceSLBRepository) GetInstanceKey() string {
	return "KRDS"
}

func (repo *InstanceSLBRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceSLBRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceSLBRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	if len(iam.IAMProjectIDs) > 0 || len(iam.IAMProjectIDs) <= 100 {
		for i := 0; i < len(iam.IAMProjectIDs); i++ {
			filters[fmt.Sprintf("ProjectId.%d", i+1)] = iam.IAMProjectIDs[i]
		}
	}

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		7,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceSLBMeta{
			LoadBalancerId:   v.InstanceID,
			LoadBalancerName: v.InstanceName,
			PublicIp:         v.InstanceIP,
		}
		ins := &InstanceSLB{
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

	level.Info(repo.logger).Log("msg", "SLB 资源加载完毕", "instance_count", len(instances))

	return
}

type DescribeLoadBalancersResponse struct {
	InstanceSet []*InstanceSLBMeta `json:"LoadBalancerDescriptions"`
	NextToken   string             `json:"NextToken"`
	TotalCount  int64              `json:"TotalCount"`
	RequestId   string             `json:"RequestId"`
}

func (repo *InstanceSLBRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	var nextToken int64 = 1

	var maxResults int64 = 10

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "SLB 资源开始加载")

getMoreInstances:

	filters["NextToken"] = nextToken
	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeLoadBalancers(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no slb.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeLoadBalancersResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse slb instance list err, %+v", err)
	}

	if totalCount == -1 {
		totalCount = response.TotalCount
	}

	for _, v := range response.InstanceSet {
		instance, err := NewInstanceSLB(v.LoadBalancerId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get slb instance fail", "id", v.LoadBalancerId)
			continue
		}
		instances = append(instances, instance)
	}

	var responseNextToken int64 = 0
	if response.NextToken != "" || len(response.NextToken) > 0 {
		responseNextToken, _ = strconv.ParseInt(response.NextToken, 10, 64)
	}

	nextToken = responseNextToken
	if nextToken < totalCount && nextToken > 0 {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "SLB 资源加载完毕")

	return
}

//NewInstanceSLBRepository
func NewInstanceSLBRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := slb.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceSLBRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
