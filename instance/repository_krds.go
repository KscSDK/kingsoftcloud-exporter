package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/krds"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("KRDS", NewInstanceRDSRepository)
}

//InstanceRDSRepository
type InstanceRDSRepository struct {
	credential config.Credential
	client     *krds.Krds
	logger     log.Logger
}

func (repo *InstanceRDSRepository) GetInstanceKey() string {
	return "KRDS"
}

func (repo *InstanceRDSRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceRDSRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceRDSRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "RDS 资源开始加载")

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		2,
		marker,
		maxResults,
		repo.credential.Region,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range l {
		meta := &InstanceRDSMeta{
			Region:               v.Region,
			DBInstanceIdentifier: v.InstanceID,
			DBInstanceName:       v.InstanceName,
			Vip:                  v.InstanceIP,
		}
		ins := &InstanceRDS{
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

	level.Info(repo.logger).Log("msg", "RDS 资源加载完毕", "instance_num", len(instances))

	return
}

type DBInstances struct {
	InstanceSet []*InstanceRDSMeta `json:"Instances"`
	Marker      int64              `json:"Marker"`
	MaxRecords  int64              `json:"MaxRecords"`
	TotalCount  int64              `json:"TotalCount"`
}

type DescribeRdsResponse struct {
	Data      DBInstances `json:"Data"`
	RequestId string      `json:"RequestId"`
}

func (repo *InstanceRDSRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	var marker int64 = 0

	var maxResults int64 = 100

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "RDS 资源开始加载")

getMoreInstances:

	filters["Marker"] = marker
	filters["MaxRecords"] = maxResults

	resp, err := repo.client.DescribeDBInstances(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no rds.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeRdsResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse rds instance list err, %+v", err)
	}

	if totalCount == -1 {
		totalCount = response.Data.TotalCount
	}

	for _, v := range response.Data.InstanceSet {
		instance, err := NewInstanceRDS(v.DBInstanceIdentifier, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get rds instance fail", "id", v.DBInstanceIdentifier)
			continue
		}
		instances = append(instances, instance)
	}

	marker += maxResults
	if marker < totalCount {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "RDS 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceRDSRepository
func NewInstanceRDSRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := krds.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceRDSRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
