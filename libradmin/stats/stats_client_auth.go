package stats

import (
	"time"

	"github.com/thebinary/radmin_exporter/libradmin"
)

type ClientAuth struct {
	Requests         int
	Responses        int
	Accepts          int
	Rejects          int
	Challenges       int
	Dup              int
	Invalid          int
	Malformed        int
	BadAuthenticator int
	Dropped          int
	UnknownTypes     int
	LastPacket       time.Time
	Elapsed          map[time.Duration]int
}

func ClientAuthStats(r *libradmin.RadminClient) (stats ClientAuth, err error) {
	var st map[string]int
	stats = ClientAuth{}

	if st, err = getStats(r, "client auth"); err != nil {
		return
	}

	if Requests, ok := st["requests"]; ok {
		stats.Requests = Requests
	}
	if Responses, ok := st["responses"]; ok {
		stats.Responses = Responses
	}
	if Accepts, ok := st["accepts"]; ok {
		stats.Accepts = Accepts
	}
	if Rejects, ok := st["rejects"]; ok {
		stats.Rejects = Rejects
	}
	if Challenges, ok := st["challenges"]; ok {
		stats.Challenges = Challenges
	}
	if Dup, ok := st["dup"]; ok {
		stats.Dup = Dup
	}
	if Invalid, ok := st["invalid"]; ok {
		stats.Invalid = Invalid
	}
	if Malformed, ok := st["malformed"]; ok {
		stats.Malformed = Malformed
	}
	if BadAuthenticator, ok := st["bad_authenticator"]; ok {
		stats.BadAuthenticator = BadAuthenticator
	}
	if Dropped, ok := st["dropped"]; ok {
		stats.Dropped = Dropped
	}
	if UnknownTypes, ok := st["unknown_types"]; ok {
		stats.UnknownTypes = UnknownTypes
	}
	if LastPacket, ok := st["last_packet"]; ok {
		stats.LastPacket = time.Unix(int64(LastPacket), 0)
	}

	stats.Elapsed = elapsedStats(&st)

	return
}
