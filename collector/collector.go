package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/metric"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const exporterNamespace = "ksc"

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, "scrape", "collector_duration_seconds"),
		"ksc_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(exporterNamespace, "scrape", "collector_success"),
		"ksc_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

//KscMonitorCollector 指标采集器, 包含多个产品的采集器
type KscMonitorCollector struct {
	Collectors map[string]*KscProductCollector
	Reloaders  map[string]*KscProductCollectorReloader
	config     *config.KscExporterConfig
	logger     log.Logger
	lock       sync.Mutex
}

const (
	defaultHandlerEnabled = true
)

var (
	collectorState = make(map[string]int)
)

func (n *KscMonitorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

func (n *KscMonitorCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))
	for name, c := range n.Collectors {
		go func(name string, c *KscProductCollector) {
			defer wg.Done()
			collect(name, c, ch, n.logger)
		}(name, c)
	}
	wg.Wait()
}

func collect(name string, c *KscProductCollector, ch chan<- prometheus.Metric, logger log.Logger) {
	begin := time.Now()
	level.Info(logger).Log("msg", "Start collect......", "name", name)

	err := c.Collect(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		level.Error(logger).Log("msg", "Collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		success = 0
	} else {
		level.Info(logger).Log("msg", "Collect done", "name", name, "duration_seconds", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

//NewKscMonitorCollector
func NewKscMonitorCollector(
	exporterConf *config.KscExporterConfig,
	logger log.Logger,
) (*KscMonitorCollector, error) {

	collectors := make(map[string]*KscProductCollector)
	reloaders := make(map[string]*KscProductCollectorReloader)

	metricRepo, err := metric.NewMetricRepository(exporterConf, logger)
	if err != nil {
		return nil, err
	}

	metricRepoCache := metric.NewMetricCache(metricRepo, logger)

	for _, namespace := range exporterConf.GetNamespaces() {

		state, exists := collectorState[namespace]
		if exists && state == 1 {
			continue
		}

		productConf, err := exporterConf.GetProductConfig(namespace)
		if err != nil {
			return nil, err
		}

		collector, err := NewKscProductCollector(
			namespace,
			metricRepoCache,
			exporterConf,
			&productConf,
			logger,
		)

		if err != nil {
			panic(fmt.Sprintf("Create product collector fail, err=%s, Namespace=%s", err, namespace))
		}

		collectors[namespace] = collector
		collectorState[namespace] = 1
		level.Info(logger).Log("msg", "Create product collecter ok", "Namespace", namespace)

		if productConf.IsReloadEnable() {
			reloadInterval := time.Duration(productConf.ReloadIntervalMinutes * int64(time.Minute))
			reloader := NewKscProductCollectorReloader(context.TODO(), collector, reloadInterval, logger)
			reloaders[namespace] = reloader

			go reloader.Run()

			level.Info(logger).Log("msg", fmt.Sprintf("reload %s instances every %d minutes", namespace, productConf.ReloadIntervalMinutes))
		}
	}

	level.Info(logger).Log("msg", "Create all product collector ok", "num", len(collectors))

	return &KscMonitorCollector{
		Collectors: collectors,
		Reloaders:  reloaders,
		config:     exporterConf,
		logger:     logger,
	}, nil
}
