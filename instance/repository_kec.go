package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/kec"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("KEC", NewInstanceKECRepository)
}

//InstanceKECRepository
type InstanceKECRepository struct {
	credential config.Credential
	client     *kec.Kec
	logger     log.Logger
}

func (repo *InstanceKECRepository) GetNamespace() string {
	return "KEC"
}

func (repo *InstanceKECRepository) GetInstanceKey() string {
	return "KEC"
}

func (repo *InstanceKECRepository) Get(id string) (instance KscInstance, err error) {
	return
}

//ListByIds
func (repo *InstanceKECRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceKECRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "KEC 资源开始加载")

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		0,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstancesKECMeta{
			AvailabilityZone: v.Region,
			InstanceId:       v.InstanceID,
			InstanceName:     v.InstanceName,
			PrivateIpAddress: v.InstanceIP,
		}
		ins := &InstanceKEC{
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

	level.Info(repo.logger).Log("msg", "KEC 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeInstancesResponse struct {
	InstancesSet  []*InstancesKECMeta `json:"InstancesSet"`
	InstanceCount int64               `json:"InstanceCount"`
	Marker        int64               `json:"Marker"`
	RequestId     string              `json:"RequestId"`
}

//ListByFilters
func (repo *InstanceKECRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	var marker int64 = 0

	var maxResults int64 = 300

	var totalCount int64 = -1

	namespace := repo.GetNamespace()
	if _, isOK := iam.OnlyIncludeProjectIDs[namespace]; isOK && !hasIncludeInstances {
		for i := 0; i < len(iam.OnlyIncludeProjectIDs[namespace]); i++ {
			filters[fmt.Sprintf("ProjectId.%d", i)] = iam.OnlyIncludeProjectIDs[namespace][i]
		}
	} else {
		if len(iam.IAMProjectIDs) > 0 || len(iam.IAMProjectIDs) <= 100 {
			for i := 0; i < len(iam.IAMProjectIDs); i++ {
				filters[fmt.Sprintf("ProjectId.%d", i)] = iam.IAMProjectIDs[i]
			}
		}
	}

	level.Info(repo.logger).Log("msg", "KEC 资源开始加载")

getMoreInstances:

	filters["Marker"] = marker
	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeInstances(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no kec instances.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeInstancesResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse kec instance list err, %+v", err)
	}

	if totalCount == -1 {
		totalCount = response.InstanceCount
	}

	for _, v := range response.InstancesSet {
		instance, err := NewInstanceKEC(v.InstanceId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "Create kec instance fail", "id", v.InstanceId)
			continue
		}
		instances = append(instances, instance)
	}

	marker += maxResults
	if marker < totalCount {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "KEC 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceKECRepository
func NewInstanceKECRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := kec.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceKECRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
