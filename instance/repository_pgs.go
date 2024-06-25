package instance

import (
	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/krds"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func init() {
	registerRepository("PGS", NewInstanceRDSRepository)
}

//InstancePGSRepository
type InstancePGSRepository struct {
	credential config.Credential
	client     *krds.Krds
	logger     log.Logger
}

func (repo *InstancePGSRepository) GetNamespace() string {
	return "PGS"
}

func (repo *InstancePGSRepository) GetInstanceKey() string {
	return "PGS"
}

func (repo *InstancePGSRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstancePGSRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstancePGSRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	var marker int64 = 1

	var maxResults int64 = 300

	var totalCount int64 = -1

getMoreInstances:

	l, count, err := DescribeMonitorInstances(
		repo.credential.AccessInstancesURL,
		repo.credential.AccessAccount,
		42,
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

	level.Info(repo.logger).Log("msg", "KRDS 资源加载完毕")

	return
}

func (repo *InstancePGSRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {
	return
}

//NewInstancePGSRepository
func NewInstancePGSRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := krds.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	repo := &InstancePGSRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
