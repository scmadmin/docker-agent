package main

import (
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
	"time"
	"github.com/devadaptive/docker-agent/stats"
	"encoding/json"
	"io/ioutil"
	"strings"
	"github.com/devadaptive/docker-agent/server"
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

	containerDataMap := stats.GatherStatsAllContainers(containers, client)

	totalTime := time.Since(t)

	for k, _ := range containerDataMap {
		jsonString, err := json.MarshalIndent(containerDataMap[k], "", "    ")
		if err != nil {
			panic(err)
		}

		fmt.Printf("cont %v:   %v\n", k, string(jsonString))
	}

	fmt.Println(totalTime)

	server.Serve()
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
