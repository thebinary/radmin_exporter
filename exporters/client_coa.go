package exporters

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thebinary/radmin_exporter/libradmin"
	"github.com/thebinary/radmin_exporter/libradmin/stats"
)

type ClientCoaExporter struct {
	mutex    sync.Mutex
	sockAddr string

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

func NewClientCoaExporter(sockAddr string) (ce *ClientCoaExporter) {
	return &ClientCoaExporter{
		sockAddr: sockAddr,
		requests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "requests_total"),
			"Current Total COA Requests",
			nil,
			nil,
		),
		responses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "responses_total"),
			"Current Total COA Responses",
			nil,
			nil,
		),
		dup: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "dup_total"),
			"Current Total COA Duplicates",
			nil,
			nil,
		),
		invalid: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "invalid_total"),
			"Current Total Invalid COA",
			nil,
			nil,
		),
		malformed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "malformed_total"),
			"Current Total Malformed COA ",
			nil,
			nil,
		),
		badAuthenticator: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "bad_authenticator_total"),
			"Current Total Bad Authenticators",
			nil,
			nil,
		),
		dropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "dropped_total"),
			"Current Total Dropped COA",
			nil,
			nil,
		),
		unknownTypes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "unknown_types_total"),
			"Current Total COA with Unknown Types",
			nil,
			nil,
		),
		elapsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_coa, "elapsed_total"),
			"COA Response Times Bucket",
			nil,
			nil,
		),
	}
}

func (c *ClientCoaExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.requests
	ch <- c.responses
	ch <- c.dup
	ch <- c.invalid
	ch <- c.malformed
	ch <- c.badAuthenticator
	ch <- c.dropped
	ch <- c.unknownTypes
}

func (c *ClientCoaExporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting client coa stats")

	c.mutex.Lock() // To protect metrics from concurrent collects.
	defer c.mutex.Unlock()

	r, err := libradmin.NewRadminClient(c.sockAddr)
	if err != nil {
		log.Printf("error connecting to control socket: %s", err)
		return
	}
	defer r.Close()

	s, err := stats.ClientCoaStats(r)
	if err != nil {
		log.Printf("error executing stats cmd: %s", err)
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
