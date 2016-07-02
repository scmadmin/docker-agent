package stats

/*
#include <unistd.h>
#include <limits.h>

int GetLongBit() {
#ifdef _SC_LONG_BIT
    int longbits;

    longbits = sysconf(_SC_LONG_BIT);
    if (longbits <  0) {
        longbits = (CHAR_BIT * sizeof(long));
    }
    return longbits;
#else
    return (CHAR_BIT * sizeof(long));
#endif
}
*/
import "C"

import 	"github.com/devadaptive/docker-agent/util"


const nanosPerSec = uint64(1000000000)

var clockTicks = uint64(GetClockTicks())


// < 1.9:  "/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-" + containerId + ".scope/cpuacct.stat"
// >= 1.9  "/sys/fs/cgroup/cpu,cpuacct/docker/" + containerId + "/cpuacct.stat"

type CPUStats struct {
	//RawValues map[string]int64 `json:"raw_data"`
	UserUsage uint64 `json:"user_usage"`
	SystemUsage uint64 `json:"system_usage"`
}

func ParseCPUData(minorVer int, id string) (*CPUStats, error) {
	m, err := util.ReadKeyValueLines(getPath(minorVer, id))
	if err != nil {
		return nil, err
	}

	return &CPUStats{
		//m,
		(m["user"] * nanosPerSec) / clockTicks,
		(m["system"] * nanosPerSec) / clockTicks,
	}, nil
}

func getPath(minorVer int, id string) string {
	if minorVer > 8 {
		return "/sys/fs/cgroup/cpu,cpuacct/docker/" + id + "/cpuacct.stat"
	} else {
		return "/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-" + id + ".scope/cpuacct.stat"
	}
}


func GetClockTicks() int {
	return int(C.sysconf(C._SC_CLK_TCK))
}
