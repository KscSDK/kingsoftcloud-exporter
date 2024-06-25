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
	"github.com/KscSDK/ksc-sdk-go/service/mongodb"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("MONGO", NewInstanceMongoRepository)
}

//InstanceMongoRepository
type InstanceMongoRepository struct {
	credential config.Credential
	client     *mongodb.Mongodb
	logger     log.Logger
}

func (repo *InstanceMongoRepository) GetNamespace() string {
	return "MONGO"
}

func (repo *InstanceMongoRepository) GetInstanceKey() string {
	return "MONGO"
}

func (repo *InstanceMongoRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceMongoRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceMongoRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "MongoDB 资源开始加载")

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		13,
		marker,
		maxResults,
		repo.credential.Region,
	)

	if err != nil {
		return nil, err
	}

	for _, v := range l {
		instanceId := v.InstanceID
		if len(instanceId) > 36 {
			instanceId = instanceId[len(instanceId)-36:]
		}

		meta := &InstanceMongoNodeMeta{
			InstanceId:   instanceId,
			InstanceName: v.InstanceName,
		}
		ins := &InstanceMongo{
			InstanceBase: InstanceBase{
				InstanceID: instanceId,
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

	level.Info(repo.logger).Log("msg", "MongoDB 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeMongoDBInstancesResponse struct {
	MongoDBInstancesResult []*InstanceMongoMeta `json:"MongoDBInstancesResult"`
	Limit                  int64                `json:"Limit"`
	Offset                 int64                `json:"Offset"`
	Total                  int64                `json:"Total"`
	RequestId              string               `json:"RequestId"`
}

func (repo *InstanceMongoRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	var marker int64 = 0

	var maxResults int64 = 100

	var totalCount int64 = -1

	level.Info(repo.logger).Log("msg", "MongoDB 资源开始加载")

	namespace := repo.GetNamespace()

	var projectIDs []string
	if _, isOK := iam.OnlyIncludeProjectIDs[namespace]; isOK && !hasIncludeInstances {
		projectIDs = make([]string, 0, len(iam.OnlyIncludeProjectIDs))
		for _, v := range iam.OnlyIncludeProjectIDs[namespace] {
			projectIDs = append(projectIDs, strconv.FormatInt(v, 10))
		}
	} else {
		if len(iam.IAMProjectIDs) > 0 {
			projectIDs = make([]string, 0, len(iam.IAMProjectIDs))
			for i := 0; i < len(iam.IAMProjectIDs); i++ {
				projectIDs = append(projectIDs, strconv.FormatInt(iam.IAMProjectIDs[i], 10))
			}
		}
	}

	if len(projectIDs) > 0 {
		filters["IamProjectId"] = strings.Join(projectIDs, ",")
	}

getMoreInstances:

	filters["Offset"] = marker
	filters["Limit"] = maxResults

	resp, err := repo.client.DescribeMongoDBInstances(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no nat resources.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeMongoDBInstancesResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse mongodb instance list err, %+v", err)
	}

	if totalCount == -1 {
		totalCount = response.Total
	}

	for _, v := range response.MongoDBInstancesResult {
		var nodes []KscInstance
		if v.InstanceType == MongoDBType_ReplicaSet {
			nodes, err = repo.ListNodeByInstanceID(v.InstanceId)
			if err != nil {
				level.Error(repo.logger).Log("msg", "get mongodb instance fail", "id", v.InstanceId)
				continue
			}
		}

		if v.InstanceType == MongoDBType_Sharding {
			nodes, err = repo.ListShardNodeByInstanceID(v.InstanceId)
			if err != nil {
				level.Error(repo.logger).Log("msg", "get mongodb instance fail", "id", v.InstanceId)
				continue
			}
		}

		instances = append(instances, nodes...)
	}

	marker += maxResults
	if marker < totalCount {
		goto getMoreInstances
	}

	level.Info(repo.logger).Log("msg", "MongoDB 资源加载完毕", "instance_num", len(instances))

	return
}

//DescribeMongoDBShardNodeResponse
type DescribeMongoDBShardNodeResponse struct {
	MongosNodeResult []*InstanceMongoNodeMeta `json:"MongosNodeResult"`
	ShardNodeResult  []*InstanceMongoNodeMeta `json:"ShardNodeResult"`
	RequestId        string                   `json:"RequestId"`
}

//ListShardNodeByInstanceID
func (repo *InstanceMongoRepository) ListShardNodeByInstanceID(instanceId string) (instances []KscInstance, err error) {
	filters := make(map[string]interface{})

	filters["InstanceId"] = instanceId

	level.Info(repo.logger).Log("msg", "MongoDB Shard Node 资源开始加载", "instance_id", instanceId)

	resp, err := repo.client.DescribeMongoDBShardNode(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no nat resources.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeMongoDBShardNodeResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse mongodb instance list err, %+v", err)
	}

	for _, v := range response.MongosNodeResult {
		instance, err := NewInstanceMongo(v.InstanceId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get mongos node instance fail", "id", v.InstanceId)
			continue
		}
		instances = append(instances, instance)
	}

	for _, v := range response.ShardNodeResult {
		instance, err := NewInstanceMongo(v.InstanceId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get shard node instance fail", "id", v.InstanceId)
			continue
		}
		instances = append(instances, instance)
	}

	level.Info(repo.logger).Log("msg", "MongoDB Shard Node 资源加载完毕", "instance_num", len(instances))

	return

}

//DescribeMongoDBInstanceNodeResponse
type DescribeMongoDBInstanceNodeResponse struct {
	MongoDBInstanceNodeResult []*InstanceMongoNodeMeta `json:"MongoDBInstanceNodeResult"`
	RequestId                 string                   `json:"RequestId"`
}

//ListNodeByInstanceID
func (repo *InstanceMongoRepository) ListNodeByInstanceID(instanceId string) (instances []KscInstance, err error) {
	filters := make(map[string]interface{})

	filters["InstanceId"] = instanceId

	level.Info(repo.logger).Log("msg", "MongoDB ReplicaSet Node 资源开始加载", "instance_id", instanceId)

	resp, err := repo.client.DescribeMongoDBInstanceNode(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no nat resources.")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeMongoDBInstanceNodeResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse mongodb instance list err, %+v", err)
	}

	for _, v := range response.MongoDBInstanceNodeResult {
		instance, err := NewInstanceMongo(v.InstanceId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get mongodb instance fail", "id", v.InstanceId)
			continue
		}
		instances = append(instances, instance)
	}

	level.Info(repo.logger).Log("msg", "MongoDB ReplicaSet Node 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceMongoRepository
func NewInstanceMongoRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := mongodb.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	repo := &InstanceMongoRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
