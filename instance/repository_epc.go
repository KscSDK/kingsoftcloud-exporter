package instance

import (
	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/epc"

	"github.com/go-kit/log"
)

func init() {
	registerRepository("EPC", NewInstanceEPCRepository)
}

//InstanceEPCRepository
type InstanceEPCRepository struct {
	credential config.Credential
	client     *epc.Epc
	logger     log.Logger
}

func (repo *InstanceEPCRepository) GetInstanceKey() string {
	return "EPC"
}

func (repo *InstanceEPCRepository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceEPCRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceEPCRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	return nil, nil
}

type DescribeEPCsResponse struct {
	HostSet    []InstanceEPCMeta `json:"HostSet"`
	NextToken  int64             `json:"NextToken"`
	TotalCount int64             `json:"TotalCount"`
	RequestId  string            `json:"RequestId"`
}

func (repo *InstanceEPCRepository) ListByFilters(filters map[string]interface{}) (instances []KscInstance, err error) {

	return
}

//NewInstanceEPCRepository
func NewInstanceEPCRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := epc.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: true,
		},
	)

	repo := &InstanceEPCRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
