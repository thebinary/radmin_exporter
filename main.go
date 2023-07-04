package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thebinary/radmin_exporter/exporters"
	"github.com/thebinary/radmin_exporter/libradmin"
)

var Version = "rc"

var showVersion bool

var (
	sockAddr         string
	sockTcp          string
	metricsEndpoint  string
	listeningAddress string
	enableAuth       bool
	enableAcct       bool
	enableCoa        bool
	enableQueue      bool
	clients          string
	gracefulStop     = make(chan os.Signal)
)

func main() {
	flag.StringVar(
		&sockAddr, "f",
		"/var/run/freeradius/freeradius.sock",
		"Radiusd control socket file",
	)

	flag.StringVar(
		&sockTcp, "t",
		"127.0.0.1:9912",
		"Radiusd control TCP socket exposed using tools like socat",
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
		&enableCoa, "coa",
		true,
		"Enable COA Statistics",
	)

	flag.BoolVar(
		&enableQueue, "queue",
		true,
		"Enable Queue Statistics",
	)

	flag.StringVar(
		&clients, "clients",
		"",
		"Comma separated list of clients for per client statistics as well",
	)

	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Parse()

	if showVersion {
		fmt.Printf("version: %s\n", Version)
		os.Exit(0)
	}

	// socket type and address
	socketType := "unix"
	if sockTcp != "" {
		socketType = "tcp"
		sockAddr = sockTcp
	}

	// listen to termination signals from the OS
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGHUP)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	//prometheus.MustRegister(version.NewCollector("radmin_exporter"))

	poll_clients := strings.Split(clients, ",")
	// remove first empty element resulting from split above
	poll_clients = poll_clients[1:]

	if enableAuth {
		log.Println("Registering Authentication statistics exporter")
		auth_radmin := libradmin.NewRadminClientWithConn(socketType, sockAddr)
		clientAuthExporter := exporters.NewClientAuthExporter(auth_radmin, poll_clients...)
		prometheus.MustRegister(clientAuthExporter)
	}

	if enableAcct {
		log.Println("Registering Accounting statistics exporter")
		acct_radmin := libradmin.NewRadminClientWithConn(socketType, sockAddr)
		clientAcctExporter := exporters.NewClientAcctExporter(acct_radmin, poll_clients...)
		prometheus.MustRegister(clientAcctExporter)
	}

	if enableCoa {
		log.Println("Registering COA statistics exporter")
		coa_radmin := libradmin.NewRadminClientWithConn(socketType, sockAddr)
		clientCoaExporter := exporters.NewClientCoaExporter(coa_radmin, poll_clients...)
		prometheus.MustRegister(clientCoaExporter)
	}

	if enableQueue {
		log.Println("Registering Queue statistics exporter")
		queue_radmin := libradmin.NewRadminClientWithConn(socketType, sockAddr)
		queueExporter := exporters.NewQueueExporter(queue_radmin, poll_clients...)
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
