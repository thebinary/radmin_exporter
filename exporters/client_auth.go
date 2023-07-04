package exporters

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thebinary/radmin_exporter/libradmin"
	"github.com/thebinary/radmin_exporter/libradmin/stats"
)

type ClientAuthExporter struct {
	mutex       sync.Mutex
	radmin      *libradmin.RadminClient
	withClients []string

	requests         *prometheus.Desc
	responses        *prometheus.Desc
	accepts          *prometheus.Desc
	rejects          *prometheus.Desc
	challenges       *prometheus.Desc
	dup              *prometheus.Desc
	invalid          *prometheus.Desc
	malformed        *prometheus.Desc
	badAuthenticator *prometheus.Desc
	dropped          *prometheus.Desc
	unknownTypes     *prometheus.Desc
	elapsed          *prometheus.Desc
	//LastPacket       time.Time
}

func NewClientAuthExporter(radmin *libradmin.RadminClient, withClients ...string) (ce *ClientAuthExporter) {
	return &ClientAuthExporter{
		radmin:      radmin,
		withClients: withClients,
		requests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "requests_total"),
			"Current Total Authorization Requests",
			nil,
			nil,
		),
		responses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "responses_total"),
			"Current Total Authorization Responses",
			nil,
			nil,
		),
		accepts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "accepts_total"),
			"Current Total Authorization Accepts",
			nil,
			nil,
		),
		rejects: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "rejects_total"),
			"Current Total Authorization Rejects",
			nil,
			nil,
		),
		challenges: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "challenges_total"),
			"Current Total Authorization Challenges",
			nil,
			nil,
		),
		dup: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "dup_total"),
			"Current Total Authorization Duplicates",
			nil,
			nil,
		),
		invalid: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "invalid_total"),
			"Current Total Invalid Authorization",
			nil,
			nil,
		),
		malformed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "malformed_total"),
			"Current Total Malformed Authorization ",
			nil,
			nil,
		),
		badAuthenticator: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "bad_authenticator_total"),
			"Current Total Bad Authenticators",
			nil,
			nil,
		),
		dropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "dropped_total"),
			"Current Total Dropped Authorization",
			nil,
			nil,
		),
		unknownTypes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "unknown_types_total"),
			"Current Total Authorization with Unknown Types",
			nil,
			nil,
		),
		elapsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sub_client_auth, "elapsed_total"),
			"Authorization Response Times Bucket",
			nil,
			nil,
		),
	}
}

func (c *ClientAuthExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.requests
	ch <- c.responses
	ch <- c.accepts
	ch <- c.rejects
	ch <- c.challenges
	ch <- c.dup
	ch <- c.invalid
	ch <- c.malformed
	ch <- c.badAuthenticator
	ch <- c.dropped
	ch <- c.unknownTypes
}

func (c *ClientAuthExporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting client auth stats")

	c.mutex.Lock() // To protect metrics from concurrent collects.
	defer c.mutex.Unlock()

	err := c.radmin.Dial()
	if err != nil {
		log.Printf("error connecting to control socket: %s", err)
		return
	}
	defer c.radmin.Close()

	s, err := stats.ClientAuthStats(c.radmin)
	if err != nil {
		log.Printf("error executing stats cmd: %s", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.requests, prometheus.CounterValue, float64(s.Requests))
	ch <- prometheus.MustNewConstMetric(c.responses, prometheus.CounterValue, float64(s.Responses))
	ch <- prometheus.MustNewConstMetric(c.accepts, prometheus.CounterValue, float64(s.Accepts))
	ch <- prometheus.MustNewConstMetric(c.rejects, prometheus.CounterValue, float64(s.Rejects))
	ch <- prometheus.MustNewConstMetric(c.challenges, prometheus.CounterValue, float64(s.Challenges))
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
