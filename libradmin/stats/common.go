package stats

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/thebinary/radmin_exporter/libradmin"
)

var SeparatorRegex *regexp.Regexp

func init() {
	SeparatorRegex = regexp.MustCompile(
		fmt.Sprintf("%s+", string([]byte{libradmin.FR_SEPRATOR})),
	)
}

func getStats(r *libradmin.RadminClient, subcmd string) (stats map[string]int, err error) {
	command := fmt.Sprintf("stats %s", subcmd)
	stats = map[string]int{}
	var val int

	result, _, err := r.Execute([]byte(command))
	if err != nil {
		return stats, err
	}

	for _, line := range result {
		arr := SeparatorRegex.Split(string(line), -1)
		if len(arr) != 2 {
			continue
		}

		val, err = strconv.Atoi(arr[1])
		if err != nil {
			return
		}
		stats[arr[0]] = val
	}

	return
}

func elapsedStats(stats *map[string]int) (estats map[time.Duration]int) {
	estats = map[time.Duration]int{}

	if e, ok := (*stats)["elapsed.1us"]; ok {
		d, _ := time.ParseDuration("1us")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.10us"]; ok {
		d, _ := time.ParseDuration("10us")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.100us"]; ok {
		d, _ := time.ParseDuration("100us")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.1ms"]; ok {
		d, _ := time.ParseDuration("1ms")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.10ms"]; ok {
		d, _ := time.ParseDuration("10ms")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.100ms"]; ok {
		d, _ := time.ParseDuration("100ms")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.1s"]; ok {
		d, _ := time.ParseDuration("1s")
		estats[d] = e
	}
	if e, ok := (*stats)["elapsed.10s"]; ok {
		d, _ := time.ParseDuration("10s")
		estats[d] = e
	}
	return
}
