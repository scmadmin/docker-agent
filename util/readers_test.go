package util

import (
	"testing"
	"fmt"
	"bufio"
	"strings"
	"time"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestScannerOdd(t *testing.T) {
	const slong = "9223372036854775807"
	const long = 9223372036854775807
	scanner := bufio.NewScanner(strings.NewReader(slong))
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	slb := scanner.Bytes()
	foo, err := BytesToUInt64(slb)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", foo)
	//var sb  [8]byte
	//binary.LittleEndian.PutUint64(sb, long)
	//fmt.Printf("%v", sb)

}

func TestScanner(t *testing.T) {
	const input = `cache 62128128
rss 8337182720
rss_huge 7834959872
mapped_file 18001920
swap 593104896`
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		key := scanner.Text()
		scanner.Scan()
		value := scanner.Bytes();
		fmt.Printf("k: %v, v: %v\n", key, value)
	}
}

func TestReadKeyValueLines(t *testing.T) {

	text := `cache 62128128
rss 8337182720
unevictable 0
hierarchical_memory_limit 9223372036854775807
hierarchical_memsw_limit 9223372036854775807
total_cache 62128128
`
	filename, err := WriteStringToTempFile(text, "dockeragenttest")
	defer os.Remove(filename)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()

	var m map[string]uint64
	m, err = ReadKeyValueLines(filename)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(time.Since(start))
	assert.Equal(t, m["cache"], uint64(62128128))
	assert.Equal(t, m["unevictable"], uint64(0))
	assert.Equal(t, m["hierarchical_memory_limit"], uint64(9223372036854775807))
	assert.Equal(t, m["rss"], uint64(8337182720))
	assert.Equal(t, m["total_cache"], uint64(62128128))
}

func BenchmarkReadKeyValueLines(b *testing.B) {
	s := `long_long_long_long_key01 9223372036854775807
long_long_long_long_key02 9223372036854775807
long_long_long_long_key03 9223372036854775807
long_long_long_long_key04 9223372036854775807
long_long_long_long_key05 9223372036854775807
long_long_long_long_key06 9223372036854775807
long_long_long_long_key07 9223372036854775807
long_long_long_long_key08 9223372036854775807
long_long_long_long_key09 9223372036854775807
long_long_long_long_key10 9223372036854775807
long_long_long_long_key11 9223372036854775807
long_long_long_long_key12 9223372036854775807
long_long_long_long_key13 9223372036854775807
long_long_long_long_key14 9223372036854775807
long_long_long_long_key15 9223372036854775807
long_long_long_long_key16 9223372036854775807
long_long_long_long_key17 9223372036854775807
long_long_long_long_key18 9223372036854775807
long_long_long_long_key19 9223372036854775807
long_long_long_long_key20 9223372036854775807
`
	filename, err := WriteStringToTempFile(s, "dockeragenttest")
	defer os.Remove(filename)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
}
