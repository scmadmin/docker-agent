package docker_agent

import (
	"github.com/devadaptive/docker-agent/stats"
)

type ContainerData struct {
	Id           string `json:"container_id"`
	NetworkStats []stats.NetworkStats `json:"network"`
	CPUStats     *stats.CPUStats
	MemoryStats  *stats.MemoryStats
}
