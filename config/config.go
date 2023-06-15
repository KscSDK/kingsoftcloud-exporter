package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/KscSDK/kingsoftcloud-exporter/constant"
	"github.com/KscSDK/kingsoftcloud-exporter/util"

	"gopkg.in/yaml.v2"
)

const (
	DefaultPeriodSeconds           = 60
	DefaultDelaySeconds            = 300
	DefaultReloadIntervalMinutes   = 30
	DefaultRateLimit               = 10
	DefaultQueryMetricBatchSize    = 2000
	DefaultKS3QueryMetricBatchSize = 60
	DefaultMaxAvailableProjects    = 100

	DefaultSupportInstances = 100

	DefaultSupportProducts = 10

	ENV_AccessKey   = "KS_CLOUD_SECRET_ID"
	ENV_SecretKey   = "KS_CLOUD_SECRET_KEY"
	ENV_ServiceRole = "KS_CLOUD_SERVICE_ROLE"
	ENV_Region      = "KS_CLOUD_REGION"

	ExporterMode_Mock = "MOCK"
)

var (
	ExporterRunningRegion = ""

	ExporterRunningMode = ""

	DebugNamespaceMetrics = map[string]bool{}

	OnlyIncludeMetrics = map[string][]string{}

	Product2Namespace = map[string]string{
		"kec":       "KEC",
		"epc":       "EPC",
		"eip":       "EIP",
		"nat":       "NAT",
		"slb":       "SLB",
		"bws":       "BWS",
		"peer":      "PEER",
		"listener":  "LISTENER",  //4层监听器
		"listener7": "LISTENER7", //7层监听器
		"krds":      "KRDS",
		"kcs":       "KCS",
		"dcgw":      "DCGW",
		"ks3":       "KS3",
	}

	SupportStatisticsTypes = map[string]bool{
		"max": true,
		"min": true,
		"avg": true,
	}

	//支持多维标签的监控项云服务产品
	SupportMultiDimensionNamespaces = map[string]bool{
		"KEC": true,
		"EPC": true,
	}

	AllProductMetricsConfig = map[string][]KscMetricConfig{
		"KEC":       AllKECMetricConfigs,
		"EPC":       AllEPCMetricConfigs,
		"EIP":       AllEIPMetricConfigs,
		"NAT":       AllNATMetricConfigs,
		"SLB":       AllSLBMetricConfigs,
		"BWS":       AllBWSMetricConfigs,
		"PEER":      AllPEERMetricConfigs,
		"LISTENER":  AllLISTENERMetricConfigs,
		"LISTENER7": AllLISTENER7MetricConfigs,
		"KRDS":      AllKRDSMetricConfigs,
		"KCS":       AllKCSMetricConfigs,
		"PGS":       AllPGSMetricConfigs,
		"DCGW":      AllDCGWMetricConfigs,
		"KS3":       AllKS3MetricConfigs,
	}
)

//Credential
type Credential struct {
	AccessKey            string `yaml:"access_key"`
	SecretKey            string `yaml:"secret_key"`
	Role                 string `yaml:"role"`
	AccessAccount        string `yaml:"access_account"`
	AccessInstancesURL   string `yaml:"access_instances_url"`
	AccessMonitorMetaURL string `yaml:"access_monitor_meta_url"`
	AccessMetricMetaURL  string `yaml:"access_metric_meta_url"`
	Region               string `yaml:"region"`
	Token                string `yaml:"token"`
	ExpiredTime          int64  `yaml:"expired_time"`
	IsInternal           bool   `yaml:"is_internal"`
}

//MonitorMetricConf
type KscMetricConfig struct {
	Namespace        string            `yaml:"namespace"`
	MetricName       string            `yaml:"metric_name"`
	MetricDesc       string            `yaml:"metric_desc"`
	MetricType       int               `yaml:"metric_type"`
	Labels           []string          `yaml:"labels"`
	Dimensions       map[string]string `yaml:"myself_dimensions"`
	Statistics       []string          `yaml:"statistics"`
	MinPeriodSeconds int64             `yaml:"min_period_seconds"`
	PeriodSeconds    int64             `yaml:"period_seconds"`
	RangeSeconds     int64             `yaml:"range_seconds"`
	DelaySeconds     int64             `yaml:"delay_seconds"`
	Unit             string            `yaml:"unit"`
}

//MonitorProductConf
type KscProductConfig struct {
	Namespace             string                 `yaml:"namespace"`
	AllMetrics            bool                   `yaml:"all_metrics"`
	AllInstances          bool                   `yaml:"all_instances"`
	ExtraLabels           []string               `yaml:"extra_labels"`
	OnlyIncludeMetrics    []string               `yaml:"only_include_metrics"`
	ExcludeMetrics        []string               `yaml:"exclude_metrics"`
	InstanceFilters       map[string]interface{} `yaml:"instance_filters"`
	OnlyIncludeInstances  []string               `yaml:"only_include_instances"`
	IncludeInstances      map[string]bool        `yaml:"-"`
	OnlyIncludeProjects   []int64                `yaml:"only_include_projects"`
	ExcludeInstances      []string               `yaml:"exclude_instances"`
	CustomQueryDimensions []map[string]string    `yaml:"custom_query_dimensions"`
	Statistics            []string               `yaml:"statistics_types"`
	PeriodSeconds         int64                  `yaml:"period_seconds"`
	RangeSeconds          int64                  `yaml:"range_seconds"`
	DelaySeconds          int64                  `yaml:"delay_seconds"`
	MetricNameType        int32                  `yaml:"metric_name_type"` // 1=大写转下划线, 2=全小写
	ReloadIntervalMinutes int64                  `yaml:"reload_interval_minutes"`
	Metrics               []KscMetricConfig      `yaml:"metrics"`
	DebugMetrics          []string               `yaml:"debug_metrics"` //打印监控项
}

func (p *KscProductConfig) IsReloadEnable() bool {
	if util.IsStrInList(constant.NotSupportInstanceNamespaces, p.Namespace) {
		return false
	}
	return true
}

//KscExporterConfig
type KscExporterConfig struct {
	Credential           Credential         `yaml:"credential"`
	Products             []KscProductConfig `yaml:"product_conf"`
	RateLimit            float64            `yaml:"rate_limit"`
	MetricQueryBatchSize int                `yaml:"metric_query_batch_size"`
	Filename             string             `yaml:"filename"`
	CacheInterval        int64              `yaml:"cache_interval"` // 单位 s
	ExporterMode         string             `yaml:"exporter_mode"`
}

func NewConfig() *KscExporterConfig {
	return &KscExporterConfig{}
}

//LoadFile
func (c *KscExporterConfig) LoadFile(filename string) error {
	c.Filename = filename
	content, err := ioutil.ReadFile(c.Filename)
	if err != nil {
		return err
	}
	if err = yaml.UnmarshalStrict(content, c); err != nil {
		return err
	}

	if err = c.check(); err != nil {
		return err
	}

	c.fillDefault()
	return nil
}

func (c *KscExporterConfig) fillDefault() {
	if c.RateLimit <= 0 {
		c.RateLimit = DefaultRateLimit
	}

	if c.MetricQueryBatchSize <= 0 || c.MetricQueryBatchSize > 100 {
		c.MetricQueryBatchSize = DefaultQueryMetricBatchSize
	}

	for i := 0; i < len(c.Products); i++ {
		for index, metric := range c.Products[i].Metrics {
			if metric.PeriodSeconds == 0 {
				c.Products[i].Metrics[index].PeriodSeconds = DefaultPeriodSeconds
			}
			if metric.DelaySeconds == 0 {
				c.Products[i].Metrics[index].DelaySeconds = c.Products[i].Metrics[index].PeriodSeconds
			}

			if metric.RangeSeconds == 0 {
				metric.RangeSeconds = metric.PeriodSeconds
			}

			if len(c.Products[i].DebugMetrics) > 0 {
				for _, debugMetric := range c.Products[i].DebugMetrics {
					DebugNamespaceMetrics[debugMetric] = true
				}
			}
		}
	}
	for index, product := range c.Products {
		if product.ReloadIntervalMinutes <= 0 {
			c.Products[index].ReloadIntervalMinutes = DefaultReloadIntervalMinutes
		}
	}
}

func (c *KscExporterConfig) check() (err error) {

	if c.Credential.AccessKey == "" {
		c.Credential.AccessKey = os.Getenv(ENV_AccessKey)
	}

	if c.Credential.SecretKey == "" {
		c.Credential.SecretKey = os.Getenv(ENV_SecretKey)
	}

	if c.Credential.AccessKey != "" && c.Credential.SecretKey != "" {
		c.Credential.Role = "" // 优先使用密钥，根据 role 是否为空判断使用密钥还是 role
	} else if c.Credential.Role == "" {
		c.Credential.Role = os.Getenv(ENV_ServiceRole)
		if c.Credential.Role == "" {
			return fmt.Errorf("credential.access_key or credential.secret_key or credential.role is empty, must be set")
		}
	}

	if c.Credential.Region == "" {
		c.Credential.Region = os.Getenv(ENV_Region)
		if c.Credential.Region == "" {
			return fmt.Errorf("credential.region is empty, must be set")
		}
	}

	ExporterRunningRegion = c.Credential.Region

	if strings.ToUpper(c.ExporterMode) == ExporterMode_Mock {
		ExporterRunningMode = ExporterMode_Mock
	}

	if len(c.Products) > DefaultSupportProducts {
		return fmt.Errorf("exporter can support up to %d products at the same time.", DefaultSupportProducts)
	}

	for i := 0; i < len(c.Products); i++ {
		if _, exists := Product2Namespace[strings.ToLower(c.Products[i].Namespace)]; !exists {
			return fmt.Errorf("namespace productName not support, %s", c.Products[i].Namespace)
		}

		if len(c.Products[i].Metrics) <= 0 {
			if _, isExists := AllProductMetricsConfig[strings.ToUpper(c.Products[i].Namespace)]; isExists {
				c.Products[i].Metrics = AllProductMetricsConfig[strings.ToUpper(c.Products[i].Namespace)]
			}
		}

		if len(c.Products[i].OnlyIncludeProjects) > DefaultMaxAvailableProjects {
			return fmt.Errorf("namespace exceeds the maximum number of configurable projects, %s", c.Products[i].Namespace)
		}

		if len(c.Products[i].OnlyIncludeInstances) > 0 {
			c.Products[i].IncludeInstances = make(map[string]bool)
			for _, v := range c.Products[i].OnlyIncludeInstances {
				c.Products[i].IncludeInstances[v] = true
			}
		}

		OnlyIncludeMetrics = make(map[string][]string)
		if len(c.Products[i].OnlyIncludeMetrics) > 0 {
			OnlyIncludeMetrics[c.Products[i].Namespace] = c.Products[i].OnlyIncludeMetrics
		}
	}

	return nil
}

func (c *KscExporterConfig) GetNamespaces() (nps []string) {

	nsSet := map[string]struct{}{}

	for i := 0; i < len(c.Products); i++ {
		ns := GetStandardNamespaceFromCustomNamespace(c.Products[i].Namespace)
		nsSet[ns] = struct{}{}

		for j := 0; j < len(c.Products[i].Metrics); j++ {
			ns := GetStandardNamespaceFromCustomNamespace(c.Products[i].Metrics[j].Namespace)
			nsSet[ns] = struct{}{}
		}
	}

	for np := range nsSet {
		nps = append(nps, np)
	}
	return
}

func (c *KscExporterConfig) GetMetricConfigMap(namespace string) map[string]KscMetricConfig {
	metricLabelsMap := make(map[string]KscMetricConfig)
	for i := 0; i < len(c.Products); i++ {
		if c.Products[i].Namespace == namespace {
			for j := 0; j < len(c.Products[i].Metrics); j++ {
				ns := GetStandardNamespaceFromCustomNamespace(c.Products[i].Metrics[j].Namespace)
				if ns == namespace {
					metricLabelsMap[c.Products[i].Metrics[j].MetricName] = c.Products[i].Metrics[j]
				}
			}
		}
	}
	return metricLabelsMap
}

func (c *KscExporterConfig) GetProductConfig(namespace string) (KscProductConfig, error) {
	for _, pconf := range c.Products {
		ns := GetStandardNamespaceFromCustomNamespace(pconf.Namespace)
		if ns == namespace {
			return pconf, nil
		}
	}
	return KscProductConfig{}, fmt.Errorf("namespace config not found")
}

func GetOnlyIncludeMetrics(namespace string) map[string]struct{} {
	onlyIncludeMetricsMaps := make(map[string]struct{})
	if len(OnlyIncludeMetrics[namespace]) != 0 {
		for _, v := range OnlyIncludeMetrics[namespace] {
			onlyIncludeMetricsMaps[v] = struct{}{}
		}
	}
	return onlyIncludeMetricsMaps
}

//GetMetricConfigs
func GetMetricConfigs(namespace string) ([]KscMetricConfig, error) {
	if _, isExists := AllProductMetricsConfig[namespace]; !isExists {
		return nil, fmt.Errorf(`No support namespace="%+v" product.`, namespace)
	}

	metricsConf := AllProductMetricsConfig[namespace]

	// 导出指定指标列表
	if len(OnlyIncludeMetrics[namespace]) != 0 {
		onlyIncludeMetricsMaps := GetOnlyIncludeMetrics(namespace)
		metricsCount := len(metricsConf)
		configs := make([]KscMetricConfig, 0, metricsCount)
		for i := 0; i < metricsCount; i++ {
			if _, isExist := onlyIncludeMetricsMaps[metricsConf[i].MetricName]; isExist {
				configs = append(configs, metricsConf[i])
			}
		}

		return configs, nil
	}

	return metricsConf, nil
}

func GetStandardNamespaceFromCustomNamespace(cns string) string {
	sns, exists := Product2Namespace[strings.ToLower(cns)]
	if !exists {
		panic(fmt.Sprintf("Product not support, namespace=%s", cns))
	}

	return sns
}

func IsSupportMultiDimensionNamespace(namespace string) bool {
	if _, isExists := SupportMultiDimensionNamespaces[strings.ToUpper(namespace)]; isExists {
		return true
	}
	return false
}
