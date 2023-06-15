package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"time"

	// "github.com/tencentyun/tencentcloud-exporter/pkg/common"
	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/google/uuid"

	"github.com/go-kit/log"
)

var (
	factoryMap = make(map[string]func(*config.KscExporterConfig, log.Logger) (InstanceRepository, error))
)

//InstanceRepository 每个产品的实例对象的Repository
type InstanceRepository interface {
	// 获取实例id
	GetInstanceKey() string

	// 根据id, 获取实例对象
	Get(id string) (KscInstance, error)

	// 根据id列表, 获取所有的实例对象
	ListByIds(ids []string) ([]KscInstance, error)

	// 根据filters, 获取符合条件的所有实例对象
	ListByFilters(filters map[string]interface{}, hasIncludeInstances bool) ([]KscInstance, error)

	ListByMonitors(filters map[string]interface{}) ([]KscInstance, error)
}

func NewInstanceRepository(namespace string, conf *config.KscExporterConfig, logger log.Logger) (InstanceRepository, error) {
	f, exists := factoryMap[namespace]
	if !exists {
		return nil, fmt.Errorf("Namespace not support, namespace=%s ", namespace)
	}
	return f(conf, logger)
}

type MonitorInstance struct {
	InstanceID   string `json:"instanceId"`
	InstanceIP   string `json:"instanceIP"`
	InstanceName string `json:"instanceName"`
	Region       string `json:"region"`
}

//DescribeMonitorInstancesResponse
type DescribeMonitorInstancesResponse struct {
	Code         string             `json:"code"`
	Message      string             `json:"message"`
	ResourceList []*MonitorInstance `json:"resourceList"`
	TotalCount   int64              `json:"totalCount"`
	RequestId    string             `json:"requestId"`
}

//DescribeMonitorInstances
func DescribeMonitorInstances(innerURL, accountId string, pType int, page, pageSize int64, region string) ([]*MonitorInstance, int64, error) {

	if len(innerURL) <= 0 {
		return nil, 0, fmt.Errorf("mock inner url is empty.")
	}

	apiURL := fmt.Sprintf("%s&PageIndex=%d&PageSize=%d&ProductType=%d&RegionKey=%s",
		innerURL,
		page,
		pageSize,
		pType,
		region,
	)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, 0, err
	}

	requestID := uuid.New().String()

	req.Header = http.Header{
		textproto.CanonicalMIMEHeaderKey("Content-Type"):     []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("Accept"):           []string{"application/json"},
		textproto.CanonicalMIMEHeaderKey("X-KSC-ACCOUNT-ID"): []string{accountId},
		textproto.CanonicalMIMEHeaderKey("X-Ksc-Region"):     []string{region},
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
		return nil, 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, 0, errors.New(string(body))
	}

	var response DescribeMonitorInstancesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, 0, fmt.Errorf("parse instance list err, %+v", err)
	}

	return response.ResourceList, response.TotalCount, nil
}

// 将TcInstanceRepository注册到factoryMap中
func registerRepository(namespace string, factory func(*config.KscExporterConfig, log.Logger) (InstanceRepository, error)) {
	factoryMap[namespace] = factory
}
