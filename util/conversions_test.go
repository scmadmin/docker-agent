package util

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"

	"time"
	"strconv"
	"strings"
)

var maxints = "18446744073709551615"
var medints = "456789"
var maxintb = []byte(maxints)
var medintb = []byte(maxints)


func BenchmarkBytesToUInt64(b *testing.B) {
	var val = uint64(0)
	var err error
	for n := 0; n < b.N; n++ {
		val, err = BytesToUInt64(maxintb)
	}
	assert.NotZero(b, val)
	assert.Nil(b, err)
}

func BenchmarkParseUint(b *testing.B) {
	var val = uint64(0)
	var err error
	for n := 0; n < b.N; n++ {
		val, err = strconv.ParseUint(maxints, 10, 64)
	}
	assert.NotZero(b, val)
	assert.Nil(b, err)
}

func BenchmarkStringJoin(b *testing.B) {
	s := strings.Join([]string{"/sys/fs/cgroup/memory/docker",
		"61c7f6cee8bfdb783aff6687f2697d6126063bd4e29ec2092e660bd2eeb6b346",
		"memory.usage_in_bytes"}, "/")
	assert.NotNil(b, s)
}

func BenchmarkStringConcat(b *testing.B) {
	s := "/sys/fs/cgroup/memory/docker/" +
	"61c7f6cee8bfdb783aff6687f2697d6126063bd4e29ec2092e660bd2eeb6b346" +
	"/memory.usage_in_bytes"
	assert.NotNil(b, s)
}

func x(s string) {
}

func xtestSpeed(t *testing.T) {
	t.SkipNow()
	//maxints := "18446744073709551615"
	maxints := "456789"
	maxintb := []byte(maxints)

	var val = uint64(0)
	start := time.Now()
	for i:=0; i < 50000000; i++ {
		val, _ = BytesToUInt64(maxintb)
	}
	tot := time.Since(start)
	fmt.Println(tot)
	fmt.Println(val)

	start = time.Now()
	for i:=0; i < 50000000; i++ {
		val, _ = strconv.ParseUint(maxints, 10, 64)
	}
	tot = time.Since(start)
	fmt.Println(tot)
	fmt.Println(val)

}


func TestBytesToInt64(t *testing.T) {
	maxintb := []byte("9223372036854775807")
	maxint, err := BytesToInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, int64(9223372036854775807), maxint)
}


func TestBytesToInt64Two(t *testing.T) {

	maxintb := []byte("-9223372036854775807")
	maxint, err := BytesToInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, int64(-9223372036854775807), maxint)
}

func TestBytesToInt64Three(t *testing.T) {

	maxintb := []byte("0")
	maxint, err := BytesToInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, int64(0), maxint)
}

func TestBytesToInt64Four(t *testing.T) {
	maxintb := []byte("1234E678")
	_, err := BytesToInt64(maxintb);
	assert.EqualError(t, err, "unknown value in byte: E parsing 1234E678")
}

func TestBytesToInt64Five(t *testing.T) {
	maxintb := []byte("10000000000000000000")
	_, err := BytesToInt64(maxintb);
	assert.Error(t, err, "value will overflow int64: 10000000000000000000")
}

func TestBytesToInt64Six(t *testing.T) {
	maxintb := []byte("+9223372036854775807")
	maxint, err := BytesToInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, int64(9223372036854775807), maxint)
}

func TestBytesToInt64Seven(t *testing.T) {
	maxintb := []byte("9223372036854775808")
	_, err := BytesToInt64(maxintb);
	assert.EqualError(t, err, "value will overflow int64: 9223372036854775808")
}

func TestBytesToInt64Eight(t *testing.T) {
	maxintb := []byte("9223372036854775817")
	_, err := BytesToInt64(maxintb);
	assert.EqualError(t, err, "value will overflow int64: 9223372036854775817")

}

func TestBytesToInt64Nine(t *testing.T) {
	maxintb := []byte("92233X2036854775808")
	_, err := BytesToInt64(maxintb);
	assert.EqualError(t, err, "unknown value in byte: X parsing 92233X2036854775808")
}

//

func TestBytesToUInt64(t *testing.T) {
	maxintb := []byte("18446744073709551615")
	maxint, err := BytesToUInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(18446744073709551615), maxint)
}



func TestBytesToUInt64Three(t *testing.T) {

	maxintb := []byte("0")
	maxint, err := BytesToUInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(0), maxint)
}

func TestBytesToUInt64Four(t *testing.T) {
	maxintb := []byte("1234E678")
	_, err := BytesToUInt64(maxintb);
	assert.EqualError(t, err, "unknown value in byte: E parsing 1234E678")
}

func TestBytesToUInt64Five(t *testing.T) {
	maxintb := []byte("100000000000000000000")
	_, err := BytesToUInt64(maxintb);
	assert.EqualError(t, err, "value will overflow uint64: 100000000000000000000")
}

func TestBytesToUInt64Six(t *testing.T) {
	maxintb := []byte("+18446744073709551615")
	maxint, err := BytesToUInt64(maxintb);
	if (err != nil) {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(18446744073709551615), maxint)
}

func TestBytesToUInt64Seven(t *testing.T) {
	maxintb := []byte("18446744073709551616")
	_, err := BytesToUInt64(maxintb);
	assert.EqualError(t, err, "value will overflow uint64: 18446744073709551616")

}

func TestBytesToUInt64Eight(t *testing.T) {
	maxintb := []byte("18446744073709551620")
	_, err := BytesToUInt64(maxintb);
	assert.EqualError(t, err, "value will overflow uint64: 18446744073709551620")

}

func TestBytesToUInt64Nine(t *testing.T) {
	maxintb := []byte("18446744X73709551616")
	_, err := BytesToUInt64(maxintb);
	assert.EqualError(t, err, "unknown value in byte: X parsing 18446744X73709551616")

}