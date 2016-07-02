package stats

import "github.com/devadaptive/docker-agent/util"

type MemoryStats struct {
	Stats map[string]uint64 `json:"stats"`
}

func ParseMemoryStats(minorVer int, id string) (*MemoryStats, error) {
	m, err := util.ReadKeyValueLines(getMemStatPath(minorVer, id))
	if err != nil {
		return nil, err
	}
	return &MemoryStats{Stats:m}, nil
}

func getMemStatPath(minorVer int, id string) string {
	if minorVer > 8 {
		return "/sys/fs/cgroup/memory/docker/" + id + "/memory.stat"
	} else {
		return "/sys/fs/cgroup/memory/system.slice/docker-" + id + ".scope/memory.stat"
	}
}