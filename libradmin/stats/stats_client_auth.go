package stats

import (
	"github.com/thebinary/radmin_exporter/libradmin"
)

type ClientAuth struct {
	ClientCommon
	Accepts    int
	Rejects    int
	Challenges int
}

func (stats *ClientAuth) PopulateWithStatsMap(st *map[string]int) {
	if Accepts, ok := (*st)["accepts"]; ok {
		stats.Accepts = Accepts
	}
	if Rejects, ok := (*st)["rejects"]; ok {
		stats.Rejects = Rejects
	}
	if Challenges, ok := (*st)["challenges"]; ok {
		stats.Challenges = Challenges
	}

	stats.ClientCommon.PopulateWithStatsMap(st)
}

func ClientAuthStats(r *libradmin.RadminClient) (stats ClientAuth, err error) {
	stats = ClientAuth{}
	err = fetchStats(r, "auth", "", &stats)
	return stats, err
}

func ClientAuthStatsForClient(r *libradmin.RadminClient, ipaddr string) (stats ClientAuth, err error) {
	stats = ClientAuth{}
	err = fetchStats(r, "auth", ipaddr, &stats)
	return stats, err
}
