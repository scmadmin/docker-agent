package util

import (
	"bufio"
	"os"
	"io/ioutil"
)


func ReadKeyValueLines(filename string) (map[string]uint64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	m := make(map[string]uint64)
	for scanner.Scan() {
		key := scanner.Text()
		if (!scanner.Scan()) {
			break
		}
		sval := scanner.Bytes()

		value, _ := BytesToUInt64(sval)
		m[key] = value
	}
	return m, nil;
}

func WriteStringToTempFile(s string, prefix string) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), prefix)
	//defer os.Remove(file.Name())
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(s);
	if err != nil {
		return "", err
	}
	file.Sync()
	return file.Name(), nil
}
