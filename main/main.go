package main

import (
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
	"os"
	"bufio"
	"io"
	"time"
	"github.com/devadaptive/docker-agent"
	"github.com/devadaptive/docker-agent/stats"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Temp struct {
	Id   string `json:"id"`
	Json string `json:"Json"`
	Lost string `json:"lost"`
}

func mainCrap() {
	t := time.Now()

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	client, err := client.NewClient("http://localhost:4243", "v1.22", nil, defaultHeaders)
	//r, err := http.Get("http://localhost:4243/containers/json")
	if (err != nil) {
		panic(err)
	}
	options := types.ContainerListOptions{All: false}
	containers, err := client.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	containerDataMap := make(map[string]string)
	containerDataChan := make(chan Temp)

	for _, c := range containers {
		go getStats(c, client, containerDataChan)
	}

	for range containers {
		containerData := <-containerDataChan
		containerDataMap[containerData.Id] = containerData.Json
	}

	totalTime := time.Since(t)
	for k, v := range containerDataMap {
		fmt.Printf("cont %v:   %v\n", k, v)
	}

	fmt.Println(totalTime)
}



func main() {
	t := time.Now()

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	client, err := client.NewClient("http://localhost:4243", "v1.22", nil, defaultHeaders)
	//r, err := http.Get("http://localhost:4243/containers/json")
	if (err != nil) {
		panic(err)
	}
	options := types.ContainerListOptions{All: false}
	containers, err := client.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	containerDataMap := make(map[string]docker_agent.ContainerData)
	containerDataChan := make(chan docker_agent.ContainerData)

	for _, c := range containers {
		go gatherContainerData(c, client, containerDataChan)
	}

	for range containers {
		containerData := <-containerDataChan
		containerDataMap[containerData.Id] = containerData
	}

	totalTime := time.Since(t)

	for k, _ := range containerDataMap {
		jsonString, err := json.Marshal(containerDataMap[k])
		if err != nil {
			panic(err)
		}

		fmt.Printf("cont %v:   %v\n", k, string(jsonString))
	}

	fmt.Println(totalTime)
}

func getInspectStats(container types.Container, client *client.Client, dc chan Temp) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 60)
	defer cancel()
	_, err := client.ContainerInspect(ctx, container.ID)
	if err != nil {
		panic(err)
	}
	//contJSON.AppArmorProfile
}

func getStats(container types.Container, client *client.Client, dc chan Temp) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 60)
	defer cancel()
	rc, err := client.ContainerStats(ctx, container.ID, false)
	if err != nil {
		panic(err)
	}
	defer rc.Close()
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	result := Temp{container.ID, string(b), ""}
	dc <- result
}

func gatherContainerData(container types.Container, cli *client.Client, dataChannel chan docker_agent.ContainerData) {
	pid, err := getPid(container.ID)
	if (err != nil) {
		//todo:  use inspect API to get the pid the slow way for this container.
		panic(err)
	}
	netstats, err := stats.ParseNetworkData(11, pid)
	cpuData, err := stats.ParseCPUData(11, container.ID)
	memoryStats, err := stats.ParseMemoryStats(11, container.ID)
	containerData := docker_agent.ContainerData{container.ID, netstats, cpuData, memoryStats}
	dataChannel <- containerData


	//getInspectStats(container, cli, nil)
}

func doDockerStats(contId string, cli client.Client) {
	//cli.ContainerStats()
}

func collectNetworkData(pid string) []stats.NetworkStats {
	counters, err := stats.ParseNetworkData(9, pid)

	if err != nil {
		panic(err)
	}

	return counters
	/*
	originalNS, err := netns.Get()
	defer originalNS.Close()

	if (err != nil) {
		panic(err)
	}

	ns, err:=netns.GetFromPid(pid);
	defer ns.Close();
	if (err != nil) {
		panic(err)
	}
	netns.Set(ns);
	//inside namespace
	*/

}

func getPid(contId string) (string, error) {
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


// xtoi2 converts the next two hex digits of s into a byte.
// If s is longer than 2 bytes then the third byte must be e.
// If the first two bytes of s are not hex digits or the third byte
// does not match e, false is returned.
func xtoi2(s string, e byte) (byte, bool) {
	if len(s) > 2 && s[2] != e {
		return 0, false
	}
	n, ei, ok := xtoi(s[:2], 0)
	return byte(n), ok && ei == 2
}
// Hexadecimal to integer starting at &s[i0].
// Returns number, new offset, success.
func xtoi(s string, i0 int) (n int, i int, ok bool) {
	n = 0
	for i = i0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i] - 'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i] - 'A') + 10
		} else {
			break
		}
		if n >= 0xFFFFFF {
			return 0, i, false
		}
	}
	if i == i0 {
		return 0, i, false
	}
	return n, i, true
}
// Count occurrences in s of any bytes in t.
func countAnyByte(s string, t string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(t, s[i]) >= 0 {
			n++
		}
	}
	return n
}

// Split s at any bytes in t.
func splitAtBytes(s string, t string) []string {
	a := make([]string, 1 + countAnyByte(s, t))
	n := 0
	last := 0
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(t, s[i]) >= 0 {
			if last < i {
				a[n] = string(s[last:i])
				n++
			}
			last = i + 1
		}
	}
	if last < len(s) {
		a[n] = string(s[last:])
		n++
	}
	return a[0:n]
}
