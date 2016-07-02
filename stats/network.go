package stats

import (
	"bufio"
	"os"
	"errors"
	"github.com/devadaptive/docker-agent/util"
	"strings"
	"fmt"
	"strconv"
)


//todo:  mapping between interface name and ipv6 addr can be found func init() {
//todo:    /proc/<pid>/net/if_inet6   the mac address can be extracted from that
//todo:    ipv6 address.  for instance, fe800000000000000042acfffe140002 corresponds to
//todo:    mac 02:42:AC:14:00:02

const totalValues = 16

type NetworkStats struct {
	Interface    string `json:"interface"`
	//
	RXBytes      int64 `json:"rx_bytes"`
	RXPackets    int64 `json:"rx_packets"`
	RXErrs       int64 `json:"rx_errs"`
	RXDrop       int64 `json:"rx_drop"`
	RXFifo       int64 `json:"rx_fifo"`
	RXFrame      int64 `json:"rx_frame"`
	RXCompressed int64 `json:"rx_compressed"`
	RXMulticast  int64 `json:"rx_multicast"`

	TXBytes      int64 `json:"tx_bytes"`
	TXPackets    int64 `json:"tx_packets"`
	TXErrs       int64 `json:"tx_errs"`
	TXDrop       int64 `json:"tx_drop"`
	TXFifo       int64 `json:"tx_fifo"`
	TXFrame      int64 `json:"tx_frame"`
	TXCompressed int64 `json:"tx_compressed"`
	TXMulticast  int64 `json:"tx_multicast"`
	//
	MACAddress  string `json:"mac"`
}

type Splitter struct {
	start bool
}

func (s *Splitter) SplitterFunc() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if (s.start) {
			newlineCount := 0
			for i := 0; i < len(data); i++ {
				if data[i] == '\n' {
					newlineCount++
				}
				if (newlineCount == 2) {
					s.start = false
					//fmt.Printf("chunk: %v\n", string (data[0:i + 1]))
					return i, data[0:i + 1], nil
				}
			}
			return len(data) - 1, nil, nil
		} else {
			return bufio.ScanWords(data, atEOF)
		}
	}
}

type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

func ParseNetworkData(minorVersion int, pid string) ([]NetworkStats, error) {
	netDevFilename := "/proc/" + pid + "/net/dev"
	ifInet6Filename := "/proc/" + pid + "/net/if_inet6"
	return parseNetworkData(netDevFilename, ifInet6Filename)
}

func parseNetworkData(netDevFilename string, ifInet6Filename string) ([]NetworkStats, error) {
	file, err := os.Open(netDevFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	netScanner := bufio.NewScanner(file)
	splitter := Splitter{start:true}
	netScanner.Split(splitter.SplitterFunc())

	if (!netScanner.Scan()) {
		return nil, errors.New("network file malformed: " + netDevFilename)
	}

	values := [totalValues]int64{}
	stats := make([]NetworkStats, 0)

	for netScanner.Scan() {
		name := netScanner.Text()
		//fmt.Println("interface: " + name)
		for i := 0; i < totalValues; i++ {
			values[i], err = scanAndBytesAndValue(netScanner)
			if err != nil {
				return nil, errors.New(err.Error() + " in " + netDevFilename)
			}
		}
		if (name == "lo:") {
			continue
		}

		s := NewNetworkStats(name, values)
		stats = append(stats, *s)
	}

	macs, errs := discoverMACAddress(ifInet6Filename)
	if errs != nil {
		for i, _ := range errs {
			fmt.Println(errs[i])
		}
	}

	for i, _ := range stats {
		if macAddr, ok := macs[stats[i].Interface]; ok {
			stats[i].MACAddress = macAddr
		}
	}
	return stats, nil

}

func discoverMACAddress(path string) (map[string]string, []error) {
	results := make(map[string]string)
	var errs []error
	file, err := os.Open(path)
	if err != nil {
		errs = util.AppendError(errs, errors.New("unable to open if_inet6 file.  skipping."))
		return results, errs
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		f := strings.Fields(scanner.Text())

		if len(f) < 6 || f[5] == "lo" {
			continue
		}
		ipv6Len := len(f[0])
		if (ipv6Len < 32) {
			errs = util.AppendError(errs, errors.New("malformed ipv6 address: " + f[0]))
			continue
		}
		m1 := f[0][16:18]
		m2 := f[0][18:20]
		m3 := f[0][20:22]
		m4 := f[0][26:28]
		m5 := f[0][28:30]
		m6 := f[0][30:32]

		//xor the first byte of the mac address
		m1i, err := strconv.ParseUint(m1, 16, 8)
		if (err != nil) {
			errs = util.AppendError(errs, errors.New("couldn't figure out mac address: " + err.Error()))
			continue
		}
		m1i = m1i ^ 2
		m1 = fmt.Sprintf("%02x", m1i)

		//fmt.Println(f[0])
		// + concat is fastest in golang, apparently.  see https://goo.gl/Au77tx
		macAddr := m1 + ":" + m2 + ":" + m3 + ":" + m4 + ":" + m5 + ":" + m6
		//fmt.Println(macAddr)
		results[f[5]] = macAddr
	}
	return results, errs
}

func NewNetworkStats(name string, values [totalValues]int64) *NetworkStats {
	return &NetworkStats{
		Interface:    strings.TrimSuffix(name, ":"),
		RXBytes:      values[ 0],
		RXPackets:    values[ 1],
		RXErrs:       values[ 2],
		RXDrop:       values[ 3],
		RXFifo:       values[ 4],
		RXFrame:      values[ 5],
		RXCompressed: values[ 6],
		RXMulticast:  values[ 7],
		TXBytes:      values[ 8],
		TXPackets:    values[ 9],
		TXErrs:       values[10],
		TXDrop:       values[11],
		TXFifo:       values[12],
		TXFrame:      values[13],
		TXCompressed: values[14],
		TXMulticast:  values[15],
	}
}

func scanAndBytesAndValue(scanner *bufio.Scanner) (int64, error) {
	if (!scanner.Scan()) {
		return -1, errors.New("unexpected end of network file")
	}
	b := scanner.Bytes()
	//fmt.Printf("bytes: %v\n", string(b))
	return util.BytesToInt64(b)
}
