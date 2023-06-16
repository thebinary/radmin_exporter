package stats

import (
	"github.com/thebinary/radmin_exporter/libradmin"
)

type Queue struct {
	QueueLenInternal int
	QueueLenProxy    int
	QueueLenAuth     int
	QueueLenAcct     int
	QueueLenDetail   int
	QueuePPSIn       int
	QueuePPSOut      int
}

func QueueStats(r *libradmin.RadminClient) (stats Queue, err error) {
	var st map[string]int
	stats = Queue{}

	if st, err = getStats(r, "queue"); err != nil {
		return
	}

	if QueueLenInternal, ok := st["queue_len_internal"]; ok {
		stats.QueueLenInternal = QueueLenInternal
	}
	if QueueLenProxy, ok := st["queue_len_proxy"]; ok {
		stats.QueueLenProxy = QueueLenProxy
	}
	if QueueLenAuth, ok := st["queue_len_auth"]; ok {
		stats.QueueLenAuth = QueueLenAuth
	}
	if QueueLenAcct, ok := st["queue_len_acct"]; ok {
		stats.QueueLenAcct = QueueLenAcct
	}
	if QueueLenDetail, ok := st["queue_len_detail"]; ok {
		stats.QueueLenDetail = QueueLenDetail
	}
	if QueuePPSIn, ok := st["queue_pps_in"]; ok {
		stats.QueuePPSIn = QueuePPSIn
	}
	if QueuePPSOut, ok := st["queue_pps_out"]; ok {
		stats.QueuePPSOut = QueuePPSOut
	}

	return
}
