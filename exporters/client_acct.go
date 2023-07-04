package exporters

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thebinary/radmin_exporter/libradmin"
	"github.com/thebinary/radmin_exporter/libradmin/stats"
)

type ClientAcctExporter struct {
	mutex       sync.Mutex
	radmin      *libradmin.RadminClient
	withClients []string

	requests         *prometheus.Desc
	responses        *prometheus.Desc
	dup              *prometheus.Desc
	invalid          *prometheus.Desc
	malformed        *prometheus.Desc
	badAuthenticator *prometheus.Desc
	dropped          *prometheus.Desc
	unknownTypes     *prometheus.Desc
	elapsed          *prometheus.Desc
	//LastPacket       time.Time
}

func NewClientAcctExporter(radmin *libradmin.RadminClient, withClients ...string) (ce *ClientAcctExporter) {
	return &ClientAcctExporter{
		radmin:      radmin,
		withClients: withClients,
		requests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "requests_total"),
			"Current Total Accounting Requests",
			nil,
			nil,
		),
		responses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "responses_total"),
			"Current Total Accounting Responses",
			nil,
			nil,
		),
		dup: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "dup_total"),
			"Current Total Accounting Duplicates",
			nil,
			nil,
		),
		invalid: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "invalid_total"),
			"Current Total Invalid Accounting",
			nil,
			nil,
		),
		malformed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "malformed_total"),
			"Current Total Malformed Accounting ",
			nil,
			nil,
		),
		badAuthenticator: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "bad_authenticator_total"),
			"Current Total Bad Authenticators",
			nil,
			nil,
		),
		dropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "dropped_total"),
			"Current Total Dropped Accounting",
			nil,
			nil,
		),
		unknownTypes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "unknown_types_total"),
			"Current Total Accounting with Unknown Types",
			nil,
			nil,
		),
		elapsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_acct, "elapsed_total"),
			"Accounting Response Times Bucket",
			nil,
			nil,
		),
	}
}

func (c *ClientAcctExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.requests
	ch <- c.responses
	ch <- c.dup
	ch <- c.invalid
	ch <- c.malformed
	ch <- c.badAuthenticator
	ch <- c.dropped
	ch <- c.unknownTypes
}

func (c *ClientAcctExporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting client acct stats")

	c.mutex.Lock() // To protect metrics from concurrent collects.
	defer c.mutex.Unlock()

	err := c.radmin.Dial()
	if err != nil {
		log.Printf("error connecting to control socket: %s", err)
		return
	}
	defer c.radmin.Close()

	s, err := stats.ClientAcctStats(c.radmin)
	if err != nil {
		log.Printf("error executing acct stats cmd: %s", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.requests, prometheus.CounterValue, float64(s.Requests))
	ch <- prometheus.MustNewConstMetric(c.responses, prometheus.CounterValue, float64(s.Responses))
	ch <- prometheus.MustNewConstMetric(c.dup, prometheus.CounterValue, float64(s.Dup))
	ch <- prometheus.MustNewConstMetric(c.invalid, prometheus.CounterValue, float64(s.Invalid))
	ch <- prometheus.MustNewConstMetric(c.malformed, prometheus.CounterValue, float64(s.Malformed))
	ch <- prometheus.MustNewConstMetric(c.badAuthenticator, prometheus.CounterValue, float64(s.BadAuthenticator))
	ch <- prometheus.MustNewConstMetric(c.dropped, prometheus.CounterValue, float64(s.Dropped))
	ch <- prometheus.MustNewConstMetric(c.unknownTypes, prometheus.CounterValue, float64(s.UnknownTypes))

	var sum float64
	elapsedHist := map[float64]uint64{}
	for k, v := range s.Elapsed {
		elapsedHist[float64(k.Seconds())] = uint64(v)
		sum += float64(v)
	}

	ch <- prometheus.MustNewConstHistogram(c.elapsed, uint64(len(elapsedHist)), sum, elapsedHist)
}
