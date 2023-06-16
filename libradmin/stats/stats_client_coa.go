package stats

import (
	"github.com/thebinary/radmin_exporter/libradmin"
)

type ClientCoa struct {
	ClientCommon
}

func ClientCoaStats(r *libradmin.RadminClient) (stats ClientCoa, err error) {
	stats = ClientCoa{}
	err = fetchStats(r, "coa", "", &stats)
	return stats, err
}

func ClientCoaStatsForClient(r *libradmin.RadminClient, ipaddr string) (stats ClientCoa, err error) {
	stats = ClientCoa{}
	err = fetchStats(r, "coa", ipaddr, &stats)
	return stats, err
}
