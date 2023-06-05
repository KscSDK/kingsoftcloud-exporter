package instance

import (
	"github.com/KscSDK/kingsoftcloud-exporter/config"
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

func (repo *InstanceBWSRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.MockInnerURL,
		repo.credential.MockAccountId,
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
			InstanceId:       v.InstanceID,
			PrivateIpAddress: v.InstanceIP,
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

	level.Info(repo.logger).Log("msg", "BWS 资源加载完毕", "instance_count", len(instances))

	return
}

func (repo *InstanceBWSRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	return nil, nil
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
