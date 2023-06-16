package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thebinary/radmin_exporter/exporters"
)

var Version = "rc"

var showVersion bool

var (
	sockAddr         string
	metricsEndpoint  string
	listeningAddress string
	enableAuth       bool
	enableAcct       bool
	enableQueue      bool
	gracefulStop     = make(chan os.Signal)
)

func main() {
	flag.StringVar(
		&sockAddr, "f",
		"/var/run/freeradius/freeradius.sock",
		"Radiusd control socket file",
	)
	flag.StringVar(
		&metricsEndpoint, "m",
		"/metrics",
		"Prometheus metrics endpoint",
	)
	flag.StringVar(
		&listeningAddress, "l",
		":9812",
		"Prometheus Listening Address",
	)

	flag.BoolVar(
		&enableAuth, "auth",
		true,
		"Enable Authentication Statistics",
	)

	flag.BoolVar(
		&enableAcct, "acct",
		true,
		"Enable Accounting Statistics",
	)

	flag.BoolVar(
		&enableQueue, "queue",
		true,
		"Enable Queue Statistics",
	)

	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Parse()

	if showVersion {
		fmt.Printf("version: %s\n", Version)
		os.Exit(0)
	}

	// listen to termination signals from the OS
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGHUP)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	//prometheus.MustRegister(version.NewCollector("radmin_exporter"))

	if enableAuth {
		log.Println("Registering Authentication statistics exporter")
		clientAuthExporter := exporters.NewClientAuthExporter(sockAddr)
		prometheus.MustRegister(clientAuthExporter)
	}
	if enableAcct {
		log.Println("Registering Accounting statistics exporter")
		clientAcctExporter := exporters.NewClientAcctExporter(sockAddr)
		prometheus.MustRegister(clientAcctExporter)
	}

	if enableQueue {
		log.Println("Registering Queue statistics exporter")
		queueExporter := exporters.NewQueueExporter(sockAddr)
		prometheus.MustRegister(queueExporter)
	}

	// listener for the termination signals from the OS
	go func() {
		log.Printf("listening and wait for graceful stop")
		sig := <-gracefulStop
		log.Printf("caught sig: %+v. Wait 2 seconds...", sig)
		time.Sleep(1 * time.Microsecond)
		log.Println("Terminate radmin-exporter on port")
		os.Exit(0)
	}()

	http.Handle(metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			 <head><title>Radmin Exporter</title></head>
			 <body>
			 <h1>Radmin Exporter</h1>
			 <p><a href='` + metricsEndpoint + `'>Metrics</a></p>
			 </body>
			 </html>
			 `))
	})
	log.Fatal(http.ListenAndServe(listeningAddress, nil))
}
