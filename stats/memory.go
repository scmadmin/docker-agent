package stats

import "github.com/devadaptive/docker-agent/util"

type MemoryStats struct {
	Usage uint64 `json:"usage"`
	Stats map[string]uint64 `json:"stats"`
	MaxUsage uint64 `json:"max_usage"`
	FailCount uint64 `json:"failcnt"`
	Limit uint64 `json:"limit"`
}

func ParseMemoryStats(minorVer int, id string) (*MemoryStats, error) {
	pathPrefix := getMemPathPrefix(minorVer, id)
	m, err := util.ReadKeyValueLines(pathPrefix + "/memory.stat")
	if err != nil {
		return nil, err
	}
	//
	u, err := util.ReadSingleUInt64ValueFile(pathPrefix + "/memory.usage_in_bytes")
	if err != nil {
		return nil, err
	}
	//
	mu, err := util.ReadSingleUInt64ValueFile(pathPrefix + "/memory.max_usage_in_bytes")
	if err != nil {
		return nil, err
	}
	//
	fc, err := util.ReadSingleUInt64ValueFile(pathPrefix + "/memory.failcnt")
	if err != nil {
		return nil, err
	}
	//
	lim, err := util.ReadSingleUInt64ValueFile(pathPrefix + "/memory.limit_in_bytes")
	if err != nil {
		return nil, err
	}

	return &MemoryStats{
		Usage: u,
		Stats:m,
		MaxUsage: mu,
		FailCount: fc,
		Limit: lim,
	}, nil
}

func getMemPathPrefix(minorVer int, id string) string {
	if minorVer > 8 {
		return "/sys/fs/cgroup/memory/docker/" + id
	} else {
		return "/sys/fs/cgroup/memory/system.slice/docker-" + id
	}
}