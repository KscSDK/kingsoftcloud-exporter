package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KscSDK/kingsoftcloud-exporter/collector"
	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/kingsoftcloud-exporter/pkg/cachedtransactiongather"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
)

func newHandler(
	cfg *config.KscExporterConfig,
	includeExporterMetrics bool,
	maxRequests int,
	logger log.Logger,
) (*http.Handler, error) {
	exporterMetricsRegistry := prometheus.NewRegistry()
	if includeExporterMetrics {
		exporterMetricsRegistry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}

	nc, err := collector.NewKscMonitorCollector(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("ksc_exporter"))
	if err := r.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register KingSoft cloud monitor collector: %s", err)
	}

	var handler http.Handler
	gatherers := prometheus.Gatherers{exporterMetricsRegistry, r}
	opts := promhttp.HandlerOpts{
		ErrorHandling:       promhttp.ContinueOnError,
		MaxRequestsInFlight: maxRequests,
		Registry:            exporterMetricsRegistry,
	}

	if cfg.CacheInterval <= 0 {
		handler = promhttp.HandlerFor(gatherers, opts)
	} else {
		handler = promhttp.HandlerForTransactional(
			cachedtransactiongather.NewCachedTransactionGather(
				prometheus.ToTransactionalGatherer(gatherers),
				time.Duration(cfg.CacheInterval)*time.Second, logger,
			), opts,
		)
	}

	if includeExporterMetrics {
		handler = promhttp.InstrumentMetricHandler(
			exporterMetricsRegistry, handler,
		)
	}
	return &handler, nil
}

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9123").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		enableExporterMetrics = kingpin.Flag(
			"web.enable-exporter-metrics",
			"Include metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Default("false").Bool()
		maxRequests = kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("0").Int()
		configFile = kingpin.Flag(
			"config.file", "Kingsoft cloud exporter configuration file.",
		).Default("exporter.yml").String()
	)

	promlogConfig := &promlog.Config{}

	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("ksc_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting ksc_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	exporterConf := config.NewConfig()
	if err := exporterConf.LoadFile(*configFile); err != nil {
		level.Error(logger).Log("msg", "Load config error", "err", err)
		os.Exit(1)
	} else {
		level.Info(logger).Log("msg", "Load config ok")
	}

	handler, err := newHandler(exporterConf, *enableExporterMetrics, *maxRequests, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Create handler fail", "err", err)
		os.Exit(1)
	}

	http.Handle(*metricsPath, *handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Kingsoft Cloud Exporter</title></head>
			<body>
			<h1>KingSoft Cloud Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
