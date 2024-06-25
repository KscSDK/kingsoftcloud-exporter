package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	monitor "github.com/KscSDK/ksc-sdk-go/service/monitorv5"
	"github.com/google/uuid"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	Namespace_KCM = "KCM"
)

func init() {
	registerRepository(Namespace_KCM, NewInstanceKCMRepository)
}

//InstanceKCMRepository
type InstanceKCMRepository struct {
	credential config.Credential
	client     *monitor.Monitorv5
	logger     log.Logger
}

func (repo *InstanceKCMRepository) GetNamespace() string {
	return Namespace_KCM
}

func (repo *InstanceKCMRepository) GetInstanceKey() string {
	return Namespace_KCM
}

func (repo *InstanceKCMRepository) Get(id string) (instance KscInstance, err error) {
	return
}

//ListByIds
func (repo *InstanceKCMRepository) ListByIds(id []string) (instances []KscInstance, err error) {
	return
}

func (repo *InstanceKCMRepository) ListByMonitors(filters map[string]interface{}) (instances []KscInstance, err error) {
	innerURL := repo.credential.AccessInstancesURL
	if len(innerURL) <= 0 {
		err = errors.New("mock inner url is empty")
		return
	}

	apiURL := fmt.Sprintf("%s&Namespace=%s", innerURL, Namespace_KCM)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return
	}

	requestID := uuid.New().String()
	req.Header = http.Header{
		textproto.CanonicalMIMEHeaderKey("Content-Type"):     []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("Accept"):           []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("X-KSC-ACCOUNT-ID"): []string{repo.credential.AccessAccount},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Region"):     []string{repo.credential.Region},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Request-Id"): []string{requestID},
	}

	c := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 100,
		},
		Timeout: 10 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = errors.New(string(body))
		return
	}

	var response DescribeInstancesKCMResponse
	if err = json.Unmarshal(body, &response); err != nil {
		err = fmt.Errorf("parse instance list err, %+v", err)
		return
	}

	for _, v := range response.DescribeInstancesResult.Instances {
		meta := &InstancesKCMMeta{
			InstanceId: v.InstanceId,
		}
		ins := &InstanceKCM{
			InstanceBase: InstanceBase{
				InstanceID: v.InstanceId,
			},
			meta: meta,
		}
		instances = append(instances, ins)
	}

	level.Info(repo.logger).Log("msg", "KCM 资源加载完毕", "instance_num", len(instances))

	return
}

type DescribeInstancesKCMResponse struct {
	DescribeInstancesResult DescribeInstancesResult `json:"describeInstancesResult"`
	ResponseMetadata        ResponseMetadata        `json:"responseMetadata"`
}

type DescribeInstancesResult struct {
	Instances []*InstancesKCMMeta `json:"instances"`
	Total     int                 `json:"total"`
}

type ResponseMetadata struct {
	RequestId string `json:"requestId"`
}

//ListByFilters
func (repo *InstanceKCMRepository) ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) (instances []KscInstance, err error) {

	filters["Namespace"] = "KCM"

	resp, err := repo.client.DescribeInstances(&filters)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no kcm instances")
	}

	respBytes, _ := json.Marshal(resp)

	var response DescribeInstancesKCMResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse kcm instance list err, %+v", err)
	}

	for _, v := range response.DescribeInstancesResult.Instances {
		instance, err := NewInstanceKCM(v.InstanceId, v)
		if err != nil {
			level.Error(repo.logger).Log("msg", "Create kcm instance fail", "id", v.InstanceId)
			continue
		}
		level.Info(repo.logger).Log("uuid", v.InstanceId)
		instances = append(instances, instance)
	}

	level.Info(repo.logger).Log("msg", "KCM 资源加载完毕", "instance_num", len(instances))

	return
}

//NewInstanceKCMRepository
func NewInstanceKCMRepository(conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	svc := monitor.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL:                      conf.Credential.UseSSL,
			UseInternal:                 conf.Credential.UseInternal,
			CustomerDomain:              conf.Credential.CustomerDomain,
			CustomerDomainIgnoreService: conf.Credential.CustomerDomainIgnoreService,
		},
	)

	repo := &InstanceKCMRepository{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
