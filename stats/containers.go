package stats

import (
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/client"
	"os"
	"bufio"
	"io"
)

type ContainerData struct {
	Id           string `json:"container_id"`
	NetworkStats []NetworkStats `json:"network"`
	CPUStats     *CPUStats
	MemoryStats  *MemoryStats
}

func GatherStatsAllContainers(containers []types.Container, client *client.Client) map[string]ContainerData {

	containerDataMap := make(map[string]ContainerData)
	containerDataChan := make(chan ContainerData)

	for _, c := range containers {
		go gatherContainerData(c, client, containerDataChan)
	}

	for range containers {
		containerData := <-containerDataChan
		containerDataMap[containerData.Id] = containerData
	}
	return containerDataMap
}

func gatherContainerData(container types.Container, cli *client.Client, dataChannel chan ContainerData) {
	pid, err := GetPid(container.ID)
	if (err != nil) {
		//todo:  use inspect API to get the pid the slow way for this container.
		panic(err)
	}
	netstats, err := ParseNetworkData(11, pid)
	cpuData, err := ParseCPUData(11, container.ID)
	memoryStats, err := ParseMemoryStats(11, container.ID)
	if err != nil {
		panic(err)
	}
	containerData := ContainerData{
		container.ID,
		netstats,
		cpuData,
		memoryStats}
	dataChannel <- containerData


	//getInspectStats(container, cli, nil)
}

func GetPid(contId string) (string, error) {
	//todo:  there could be multiple pids in here and some of them disappear.  There is one pid that should
	//todo:  last as long as the container lasts, though.  All the other pids should have the main pid as
	//todo:  their parent pid.   the main pid will have the docker daemon as the parent pid.
	//todo:  Potentially this could be leveraged to figure out the main pid without the expense of calling
	//todo:  docker inspect api on every container.
	filename := "/sys/fs/cgroup/systemd/docker/" + contId + "/cgroup.procs";

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		//todo:  unexpected no tokens error.
		panic(io.ErrUnexpectedEOF)
	}
	//r := bufio.NewReader(f)
	//s, err := Readln(r)
	//if (err != nil) {
	//	return -1, err
	//}
	pid := scanner.Text()//strconv.Atoi(s)
	return pid, err;
}
