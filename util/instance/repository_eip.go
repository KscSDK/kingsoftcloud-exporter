package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/eip"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("EIP", NewInstanceEIPRepository)
}

//InstanceEIPRepository
type InstanceEIPRepository struct {
	credential config.Credential
	client     *eip.Eip
	logger     log.Logger
}

func (repo *InstanceEIPRepository) GetNamespace() string {
	return "EIP"
}

func (repo *InstanceEIPRepository) GetInstanceKey() string {
	return "EIP"
}

func (repo *InstanceEIPRepository) Get(id string) (instance KscInstance, err error) {
	return
}

//ListByIds
func (repo *InstanceEIPRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceEIPRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "EIP 资源开始加载")

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		4,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceEIPMeta{
			AllocationId: v.InstanceID,
			PublicIp:     v.InstanceIP,
		}
		ins := &InstanceEIP{
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

	level.Info(repo.logger).Log("msg", "EIP 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeAddressesResponse struct {
	AddressesSet []*InstanceEIPMeta `json:"AddressesSet"`
	RequestId    string             `json:"RequestId"`
}

//ListByFilters
func (repo *InstanceEIPRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	NextToken := 1
	MaxResults := 300

	level.Info(repo.logger).Log("msg", "EIP 资源开始加载")

	namespace := repo.GetNamespace()
	if _, isOK := iam.OnlyIncludeProjectIDs[namespace]; isOK && !hasIncludeInstances {
		for i := 0; i < len(iam.OnlyIncludeProjectIDs[namespace]); i++ {
			filters[fmt.Sprintf("ProjectId.%d", i+1)] = iam.OnlyIncludeProjectIDs[namespace][i]
		}
	} else {
		if len(iam.IAMProjectIDs) > 0 || len(iam.IAMProjectIDs) <= 100 {
			for i := 0; i < len(iam.IAMProjectIDs); i++ {
				filters[fmt.Sprintf("ProjectId.%d", i+1)] = iam.IAMProjectIDs[i]
			}
		}
	}

getMoreInstances:

	filters["NextToken"] = NextToken
	filters["MaxResults"] = MaxResults

	resp, err := repo.client.DescribeAddresses(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no eip.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeAddressesResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse eip instance list err, %+v", err)
	}

	for _, v := range response.AddressesSet {
		instance, err := NewInstanceEIP(v.AllocationId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "Create eip instance fail", "id", v.AllocationId)
			continue
		}
		instances = append(instances, instance)
	}

	if len(response.AddressesSet) == filters["MaxResults"] {
		NextToken = NextToken + MaxResults
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "EIP 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceEIPRepository
func NewInstanceEIPRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := eip.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: conf.Credential.UseSSL,
		},
	)

	repo := &InstanceEIPRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
