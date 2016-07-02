package stats

import (
	"testing"
	"github.com/devadaptive/docker-agent/util"
	"os"
	"github.com/stretchr/testify/assert"
)

const if_inet6DataTwoAdapters = `00000000000000000000000000000001 01 80 10 80       lo
fe800000000000000042acfffe150002 11cd 40 20 80     eth1
fe800000000000000042acfffe140002 11cb 40 20 80     eth0
`

const if_inet6DataOneAdapter = `00000000000000000000000000000001 01 80 10 80       lo
fe80000000000000bc42acfffe140002 11cb 40 20 80     eth0
`

const if_inet6DataMalformedOne = `00000000000000000000000000000001 01 80 10 80       lo
fe800000000000000042acfffe140002 11cb 40 20 80
`

const if_inet6DataMalformedTwo = `00000000000000000000000000000001 01 80 10 80       lo
fe80000000000000gg42acfffe140002 11cb 40 20 80     eth0
`

const if_inet6DataShortIPV6 = `00000000000000000000000000000001 01 80 10 80       lo
fe80000000bc42acfffe140002 11cb 40 20 80     eth0
`

const oneAdapter = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0:    1000   2000    3000    4000    5000     6000     7000    8000     9000      10000    11000    12000    13000     14000       15000    16000
    lo:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
`
const twoAdapter = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0:    1000   2000    3000    4000    5000     6000     7000    8000     9000      10000    11000    12000    13000     14000       15000    16000
  eth1:    1001   2001    3001    4001    5001     6001     7001    8001     9001      10010    11001    12001    13001     14001       15001    16001
    lo:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
`
const malformedOne = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo `

const malformedTwo = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0:    1000   BOOM    3000    4000    5000     6000     7000    8000     9000      10000    11000    12000    13000     14000       15000    16000
    lo:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
`

const malformedThree = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0:    1000   2000    3000    4000    5000  `

func assertDevNetEth0Stats(t *testing.T, stats *NetworkStats) {
	assert.Equal(t, stats.Interface, "eth0", "interface name")
	assert.Equal(t, stats.RXBytes, int64(1000))
	assert.Equal(t, stats.RXPackets, int64(2000))
	assert.Equal(t, stats.RXErrs, int64(3000))
	assert.Equal(t, stats.RXDrop, int64(4000))
	assert.Equal(t, stats.RXFifo, int64(5000))
	assert.Equal(t, stats.RXFrame, int64(6000))
	assert.Equal(t, stats.RXCompressed, int64(7000))
	assert.Equal(t, stats.RXMulticast, int64(8000))
	assert.Equal(t, stats.TXBytes, int64(9000))
	assert.Equal(t, stats.TXPackets, int64(10000))
	assert.Equal(t, stats.TXErrs, int64(11000))
	assert.Equal(t, stats.TXDrop, int64(12000))
	assert.Equal(t, stats.TXFifo, int64(13000))
	assert.Equal(t, stats.TXFrame, int64(14000))
	assert.Equal(t, stats.TXCompressed, int64(15000))
	assert.Equal(t, stats.TXMulticast, int64(16000))
}

func assertDevNetEth1Stats(t *testing.T, stats *NetworkStats) {
	assert.Equal(t, stats.Interface, "eth1", "interface name")
	assert.Equal(t, stats.RXBytes, int64(1001))
	assert.Equal(t, stats.RXPackets, int64(2001))
	assert.Equal(t, stats.RXErrs, int64(3001))
	assert.Equal(t, stats.RXDrop, int64(4001))
	assert.Equal(t, stats.RXFifo, int64(5001))
	assert.Equal(t, stats.RXFrame, int64(6001))
	assert.Equal(t, stats.RXCompressed, int64(7001))
	assert.Equal(t, stats.RXMulticast, int64(8001))
	assert.Equal(t, stats.TXBytes, int64(9001))
	assert.Equal(t, stats.TXPackets, int64(10010))
	assert.Equal(t, stats.TXErrs, int64(11001))
	assert.Equal(t, stats.TXDrop, int64(12001))
	assert.Equal(t, stats.TXFifo, int64(13001))
	assert.Equal(t, stats.TXFrame, int64(14001))
	assert.Equal(t, stats.TXCompressed, int64(15001))
	assert.Equal(t, stats.TXMulticast, int64(16001))
}

func BenchmarkParseNetworkData(b *testing.B) {
	filename, ifInet6, err := writeTempFiles(oneAdapter, if_inet6DataOneAdapter)
	if (err != nil) {
		b.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	b.ResetTimer()
	stats, err := parseNetworkData(filename, ifInet6)
	b.StopTimer()
	if err != nil {
		b.Fatal(err)
	}

	if (len(stats) != 1 ) {
		b.Fatal("expected 1")
	}
}

func TestParseNetworkData(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(oneAdapter, if_inet6DataOneAdapter)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	stats, err := parseNetworkData(filename, ifInet6)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(stats); l != 1 {
		t.Errorf("expected 1 network stats, got %d", l)
	}
	assertDevNetEth0Stats(t, &stats[0])
	assert.Equal(t, stats[0].MACAddress, "be:42:ac:14:00:02")

}

func TestParseNetworkDataMalformedIfInet6(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(oneAdapter, if_inet6DataMalformedOne)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	stats, err := parseNetworkData(filename, ifInet6)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(stats); l != 1 {
		t.Errorf("expected 1 network stats, got %d", l)
	}
	assertDevNetEth0Stats(t, &stats[0])
	assert.Equal(t, stats[0].MACAddress, "")
}

func TestParseNetworkDataTwoAdapters(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(twoAdapter, if_inet6DataTwoAdapters)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	stats, err := parseNetworkData(filename, ifInet6)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(stats); l != 2 {
		t.Errorf("expected 2 network stats, got %d", l)
	}
	assertDevNetEth0Stats(t, &stats[0])
	assert.Equal(t, stats[0].MACAddress, "02:42:ac:14:00:02")
	//
	assertDevNetEth1Stats(t, &stats[1])
	assert.Equal(t, stats[1].MACAddress, "02:42:ac:15:00:02")

}

//this won't work on non-linux systems or systems without network adapters !
func TestParseNetworkDataLinux(t *testing.T) {
	stats, err := ParseNetworkData(11, "1")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotZero(t, len(stats))
}

func TestParseNetworkDataMalformed1(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(malformedOne, if_inet6DataMalformedOne)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	_, err = parseNetworkData(filename, ifInet6)
	if err == nil {
		t.Error("was expecting an error")
	}
}

func TestParseNetworkDataMalformed2(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(malformedTwo, if_inet6DataTwoAdapters)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	_, err = parseNetworkData(filename, ifInet6)
	if err == nil {
		t.Error("was expecting an error")
	}
}

func TestParseNetworkDataMalformed3(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(malformedThree, if_inet6DataTwoAdapters)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	_, err = parseNetworkData(filename, ifInet6)
	if err == nil {
		t.Error("was expecting an error")
	}
}

func TestParseNetworkDataFileDoesntExist(t *testing.T) {
	_, err := parseNetworkData("/this/file/does/nt/exist", "this/doesnt/exist/either")
	if err == nil {
		t.Error("was expecting an error")
	}
}

func TestParseNetworkDataOkDevNetMalformedIPV6(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(oneAdapter, if_inet6DataShortIPV6)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)

	stats, err := parseNetworkData(filename, ifInet6)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(stats); l != 1 {
		t.Errorf("expected 1 network stats, got %d", l)
	}
	assertDevNetEth0Stats(t, &stats[0])
	assert.Equal(t, stats[0].MACAddress, "")

}

func TestDiscoverMACAddress(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(oneAdapter, if_inet6DataTwoAdapters)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)
	m, errs := discoverMACAddress(ifInet6)

	assert.Nil(t, errs)

	assert.Contains(t, m, "eth0")
	assert.Contains(t, m, "eth1")

}

func TestDiscoverMACAddressMalformed2(t *testing.T) {
	filename, ifInet6, err := writeTempFiles(oneAdapter, if_inet6DataMalformedTwo)
	if (err != nil) {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	defer os.Remove(ifInet6)
	m, errs := discoverMACAddress(ifInet6)
	assert.EqualError(t, errs[0], "couldn't figure out mac address: strconv.ParseUint: parsing \"gg\": invalid syntax")
	assert.Empty(t, m)

}

func TestDiscoverMACAddressBadFile(t *testing.T) {

	m, errs := discoverMACAddress("/no/real/file")
	assert.EqualError(t, errs[0], "unable to open if_inet6 file.  skipping.")

	assert.Empty(t, m)

}

func writeTempFiles(netDevContent string, ifInet6Data string) (string, string, error) {
	filename, err := util.WriteStringToTempFile(netDevContent, "datest")
	if (err != nil) {
		return "", "", err
	}

	ifInet6File, err := util.WriteStringToTempFile(ifInet6Data, "datest")
	if (err != nil) {
		return "", "", err
	}
	return filename, ifInet6File, nil
}