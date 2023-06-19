package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/bws"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("BWS", NewInstanceBWSRepository)
}

//InstanceBWSRepository
type InstanceBWSRepository struct {
	credential config.Credential
	client     *bws.Bws
	logger     log.Logger
}

func (repo *InstanceBWSRepository) GetNamespace() string {
	return "BWS"
}

func (repo *InstanceBWSRepository) GetInstanceKey() string {
	return "BWS"
}

func (repo *InstanceBWSRepository) Get(id string) (instance KscInstance, err error) {
	return
}

//ListByIds
func (repo *InstanceBWSRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceBWSRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		11,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceBWSMeta{
			BandWidthShareId: v.InstanceID,
		}
		ins := &InstanceBWS{
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

	level.Info(repo.logger).Log("msg", "BWS 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeBandWidthSharesResponse struct {
	BandWidthShareSet []*InstanceBWSMeta `json:"BandWidthShareSet"`
	NextToken         string             `json:"NextToken"`
	RequestId         string             `json:"RequestId"`
}

func (repo *InstanceBWSRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	var maxResults int64 = 300

	level.Info(repo.logger).Log("msg", "BWS 资源开始加载")

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

	filters["MaxResults"] = maxResults

	resp, err := repo.client.DescribeBandWidthShares(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no bws instances.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeBandWidthSharesResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse bws instance list err, %+v", err)
	}

	for _, v := range response.BandWidthShareSet {
		instance, err := NewInstanceBWS(v.BandWidthShareId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "Create bws instance fail", "id", v.BandWidthShareId)
			continue
		}
		instances = append(instances, instance)
	}

	level.Info(repo.logger).Log("msg", "BWS 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceBWSRepository
func NewInstanceBWSRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := bws.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceBWSRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
