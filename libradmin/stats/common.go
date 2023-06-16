package stats

import (
	"regexp"
	"time"
)

var SeparatorRegex *regexp.Regexp

type ClientCommon struct {
	Requests         int
	Responses        int
	Dup              int
	Invalid          int
	Malformed        int
	BadAuthenticator int
	Dropped          int
	UnknownTypes     int
	LastPacket       time.Time
	Elapsed          map[time.Duration]int
}

func (stats *ClientCommon) PopulateWithStatsMap(st *map[string]int) {
	if Requests, ok := (*st)["requests"]; ok {
		stats.Requests = Requests
	}
	if Responses, ok := (*st)["responses"]; ok {
		stats.Responses = Responses
	}
	if Dup, ok := (*st)["dup"]; ok {
		stats.Dup = Dup
	}
	if Invalid, ok := (*st)["invalid"]; ok {
		stats.Invalid = Invalid
	}
	if Malformed, ok := (*st)["malformed"]; ok {
		stats.Malformed = Malformed
	}
	if BadAuthenticator, ok := (*st)["bad_authenticator"]; ok {
		stats.BadAuthenticator = BadAuthenticator
	}
	if Dropped, ok := (*st)["dropped"]; ok {
		stats.Dropped = Dropped
	}
	if UnknownTypes, ok := (*st)["unknown_types"]; ok {
		stats.UnknownTypes = UnknownTypes
	}
	if LastPacket, ok := (*st)["last_packet"]; ok {
		stats.LastPacket = time.Unix(int64(LastPacket), 0)
	}

	stats.Elapsed = elapsedStats(st)
}
