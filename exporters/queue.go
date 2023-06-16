package exporters

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thebinary/radmin_exporter/libradmin"
	"github.com/thebinary/radmin_exporter/libradmin/stats"
)

type QueueExporter struct {
	mutex    sync.Mutex
	sockAddr string

	queueLenInternal *prometheus.Desc
	queueLenProxy    *prometheus.Desc
	queueLenAuth     *prometheus.Desc
	queueLenAcct     *prometheus.Desc
	queueLenDetail   *prometheus.Desc
	queuePPSIn       *prometheus.Desc
	queuePPSOut      *prometheus.Desc
}

func NewQueueExporter(sockAddr string) (ce *QueueExporter) {
	return &QueueExporter{
		sockAddr: sockAddr,
		queueLenInternal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_len_internal"),
			"Internal Queue Length",
			nil,
			nil,
		),
		queueLenProxy: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_len_proxy"),
			"Proxy Queue Length",
			nil,
			nil,
		),
		queueLenAuth: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_len_auth"),
			"Authorization Queue Length",
			nil,
			nil,
		),
		queueLenAcct: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_len_acct"),
			"Accounting Queue Length",
			nil,
			nil,
		),
		queueLenDetail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_len_detail"),
			"Detail Queue Length",
			nil,
			nil,
		),
		queuePPSIn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_pps_in"),
			"Queue Packets Per Second In",
			nil,
			nil,
		),
		queuePPSOut: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_queue, "queue_pps_out"),
			"Queue Packets Per Second Out",
			nil,
			nil,
		),
	}
}

func (c *QueueExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.queueLenInternal
	ch <- c.queueLenProxy
	ch <- c.queueLenAuth
	ch <- c.queueLenAcct
	ch <- c.queueLenDetail
	ch <- c.queuePPSIn
	ch <- c.queuePPSOut
}

func (c *QueueExporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting queue stats")

	c.mutex.Lock() // To protect metrics from concurrent collects.
	defer c.mutex.Unlock()

	r, err := libradmin.NewRadminClient(c.sockAddr)
	if err != nil {
		log.Printf("error connecting to control socket: %s", err)
		return
	}
	defer r.Close()

	s, err := stats.QueueStats(r)
	if err != nil {
		log.Printf("error executing stats cmd: %s", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.queueLenInternal, prometheus.GaugeValue, float64(s.QueueLenInternal))
	ch <- prometheus.MustNewConstMetric(c.queueLenProxy, prometheus.GaugeValue, float64(s.QueueLenProxy))
	ch <- prometheus.MustNewConstMetric(c.queueLenAuth, prometheus.GaugeValue, float64(s.QueueLenAuth))
	ch <- prometheus.MustNewConstMetric(c.queueLenAcct, prometheus.GaugeValue, float64(s.QueueLenAcct))
	ch <- prometheus.MustNewConstMetric(c.queueLenDetail, prometheus.GaugeValue, float64(s.QueueLenDetail))
	ch <- prometheus.MustNewConstMetric(c.queuePPSIn, prometheus.GaugeValue, float64(s.QueuePPSIn))
	ch <- prometheus.MustNewConstMetric(c.queuePPSOut, prometheus.GaugeValue, float64(s.QueuePPSOut))
}
