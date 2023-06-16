package stats

import (
	"github.com/thebinary/radmin_exporter/libradmin"
)

type ClientAcct struct {
	ClientCommon
}

func ClientAcctStats(r *libradmin.RadminClient) (stats ClientAcct, err error) {
	stats = ClientAcct{}
	err = fetchStats(r, "acct", "", &stats)
	return stats, err
}

func ClientAcctStatsForClient(r *libradmin.RadminClient, ipaddr string) (stats ClientAcct, err error) {
	stats = ClientAcct{}
	err = fetchStats(r, "acct", ipaddr, &stats)
	return stats, err
}
