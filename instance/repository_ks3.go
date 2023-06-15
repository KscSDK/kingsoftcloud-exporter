package instance

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/ks3sdklib/aws-sdk-go/aws"
	"github.com/ks3sdklib/aws-sdk-go/aws/credentials"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type KS3Endpoint struct {
	Region     string
	PublicURL  string
	PrivateURL string
}

var KS3_Allow_EndPoints = map[string]*KS3Endpoint{
	"cn-beijing-6": &KS3Endpoint{
		Region:     "BEIJING",
		PublicURL:  "ks3-cn-beijing.ksyuncs.com",
		PrivateURL: "ks3-cn-beijing-internal.ksyuncs.com",
	},
	"cn-shanghai-2": &KS3Endpoint{
		Region:     "SHANGHAI",
		PublicURL:  "ks3-cn-shanghai.ksyuncs.com",
		PrivateURL: "ks3-cn-shanghai-internal.ksyuncs.com",
	},
	"cn-guangzhou-1": &KS3Endpoint{
		Region:     "GUANGZHOU",
		PublicURL:  "ks3-cn-guangzhou.ksyuncs.com",
		PrivateURL: "ks3-cn-guangzhou-internal.ksyuncs.com",
	},
	"cn-hongkong-2": &KS3Endpoint{
		Region:     "HONGKONG",
		PublicURL:  "ks3-cn-hk-1.ksyuncs.com",
		PrivateURL: "ks3-cn-hk-1-internal.ksyuncs.com",
	},
	"eu-east-1": &KS3Endpoint{
		Region:     "RUSSIA",
		PublicURL:  "ks3-rus.ksyuncs.com",
		PrivateURL: "ks3-rus-internal.ksyuncs.com",
	},
	"ap-singapore-1": &KS3Endpoint{
		Region:     "SINGAPORE",
		PublicURL:  "ks3-sgp.ksyuncs.com",
		PrivateURL: "ks3-sgp-internal.ksyuncs.com",
	},
	"cn-beijing-fin": &KS3Endpoint{
		Region:     "JR_BEIJING",
		PublicURL:  "ks3-jr-beijing.ksyuncs.com",
		PrivateURL: "ks3-jr-beijing-internal.ksyuncs.com",
	},
	"cn-shanghai-fin": &KS3Endpoint{
		Region:     "JR_SHANGHAI",
		PublicURL:  "ks3-jr-shanghai.ksyuncs.com",
		PrivateURL: "ks3-jr-shanghai-internal.ksyuncs.com",
	},
}

func init() {
	registerRepository("KS3", NewInstanceKS3Repository)
}

//InstanceKS3Repository
type InstanceKS3Repository struct {
	credential config.Credential
	client     *s3.S3
	logger     log.Logger
}

func (repo *InstanceKS3Repository) GetNamespace() string {
	return "KS3"
}

func (repo *InstanceKS3Repository) GetInstanceKey() string {
	return "KS3"
}

func (repo *InstanceKS3Repository) Get(id string) (instance KscInstance, err error) {
	return
}

func (repo *InstanceKS3Repository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceKS3Repository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	return
}

type ListBucketsResponse struct {
	Buckets []*InstancesKS3BucketMeta
}

func (repo *InstanceKS3Repository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	level.Info(repo.logger).Log("msg", "KS3 Buckets 资源开始加载")

	resp, err := repo.client.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no ks3 bucket list.")
	}

	respBytes, _ := json.Marshal(resp)

	var response ListBucketsResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse ks3 bucket list err, %+v", err)
	}

	for _, v := range response.Buckets {
		instance, err := NewInstanceKS3(v.Name, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "get bucket fail", "bucket_name", v.Name)
			continue
		}
		instances = append(instances, instance)
	}

	level.Info(repo.logger).Log("msg", "KS3 Buckets 资源加载完毕", "bucket_num", len(instances))

	return
}

//NewInstanceKS3Repository
func NewInstanceKS3Repository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {

	e, isExists := KS3_Allow_EndPoints[conf.Credential.Region]
	if !isExists {
		return nil, fmt.Errorf("no endpoint")
	}

	cre := credentials.NewStaticCredentials(conf.Credential.AccessKey, conf.Credential.SecretKey, "")

	svc := s3.New(&aws.Config{
		Region:      e.Region,
		Credentials: cre,
		Endpoint:    e.PublicURL,
		DisableSSL:  true,
	})

	repo := &InstanceKS3Repository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
