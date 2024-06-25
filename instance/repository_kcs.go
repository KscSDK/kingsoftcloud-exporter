package instance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/iam"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	kcsv1 "github.com/KscSDK/ksc-sdk-go/service/kcsv1"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("KCS", NewInstanceKCSRepository)
}

//InstanceKCSRepository
type InstanceKCSRepository struct {
	credential config.Credential
	client     *kcsv1.Kcsv1
	logger     log.Logger
}

func (repo *InstanceKCSRepository) GetNamespace() string {
	return "KCS"
}

func (repo *InstanceKCSRepository) GetInstanceKey() string {
	return "KCS"
}

func (repo *InstanceKCSRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceKCSRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceKCSRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "KCS 资源开始加载")

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		97,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceKCSMeta{
			Region:  v.Region,
			CacheId: v.InstanceID,
			Name:    v.InstanceName,
			Vip:     v.InstanceIP,
		}
		ins := &InstanceKCS{
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

	level.Info(repo.logger).Log("msg", "KCS 资源加载完毕", "instance_num", len(instances))

	return
}

type RedisInstances struct {
	InstanceSet []*InstanceKCSMeta `json:"list"`
	Marker      int64              `json:"offset"`
	MaxRecords  int64              `json:"limit"`
	TotalCount  int64              `json:"total"`
}

type DescribeCacheClustersResponse struct {
	Data      RedisInstances `json:"Data"`
	RequestId string         `json:"RequestId"`
}

func (repo *InstanceKCSRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	var marker int64 = 0

	var maxResults int64 = 100

	var totalCount int64 = -1

	namespace := repo.GetNamespace()
	if len(iam.OnlyIncludeProjectIDs[namespace]) > 0 && !hasIncludeInstances {
		projectIDs := make([]string, 0, len(iam.OnlyIncludeProjectIDs[namespace]))
		for _, v := range iam.OnlyIncludeProjectIDs[namespace] {
			projectIDs = append(projectIDs, strconv.FormatInt(v, 10))
		}
		filters["IamProjectId"] = strings.Join(projectIDs, ",")
	}

	level.Info(repo.logger).Log("msg", "KCS 资源开始加载")

getMoreInstances:

	filters["Offset"] = marker
	filters["Limit"] = maxResults

	resp, err := repo.client.DescribeCacheClusters(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no ksc.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeCacheClustersResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse rds instance list err, %+v", err)
	}

	if totalCount == -1 {
		totalCount = response.Data.TotalCount
	}

	for _, v := range response.Data.InstanceSet {
		instance, err := NewInstanceKCS(v.CacheId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get rds instance fail", "id", v.CacheId)
			continue
		}
		instances = append(instances, instance)
	}

	marker += maxResults
	if marker < totalCount {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "KCS 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceKCSRepository
func NewInstanceKCSRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := kcsv1.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	repo := &InstanceKCSRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
