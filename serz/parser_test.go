package serz

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func extractName(name string) string {
	name = filepath.Base(name)
	blocks := strings.Split(name, ".")
	return blocks[0]
}

func extractPackage(name string) string {
	name = filepath.Dir(name)
	name = filepath.Base(name)
	return name
}

func TestYsoserial(t *testing.T) {
	files, err := filepath.Glob("../testcases/ysoserial/*.ser")
	require.Nil(t, err)
	require.NotZero(t, len(files))

	for _, name := range files {
		data, err := ioutil.ReadFile(name)
		require.Nil(t, err)

		ser, err := FromBytes(data)
		require.Nilf(t, err, "an error is occurred in file %v", name)
		require.Truef(t, bytes.Equal(data, ser.ToBytes()), "original serz data is different from generation data in file %v", name)
	}
}

func TestJDK8u20(t *testing.T) {
	// current skipped
	t.SkipNow()

	var filename = "../testcases/pwntester/JDK8u20.ser"
	data, err := ioutil.ReadFile(filename)
	require.Nil(t, err)

	ser, err := FromBytes(data)
	require.Nilf(t, err, "an error is occurred in file %v", filename)
	require.Truef(t, bytes.Equal(data, ser.ToBytes()), "original serz data is different from generation data in file %v", filename)
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	var (
		ysosers []string
		ptsers  []string
		files   []string
	)
	var err error

	ysosers, err = filepath.Glob("../testcases/ysoserial/*.ser")
	if err != nil {
		exitCode = exitCode | 1
		goto cleanup
	}

	ptsers, err = filepath.Glob("../testcases/pwntester/*.ser")
	if err != nil {
		exitCode = exitCode | 1
		goto cleanup
	}

	files = append(ysosers, ptsers...)
	fmt.Println("| Gadget | Package | Parsed | Rebuild | Parse Time |")
	fmt.Println("|--------|--------|--------|--------|--------|")
	for _, name := range files {
		data, err := ioutil.ReadFile(name)
		if err != nil {
			exitCode = exitCode | 1
			goto cleanup
		}

		parseFlag := "❌"
		rebuildFlag := "❌"
		start := time.Now()
		serialization, err := FromBytes(data)
		duration := time.Since(start)

		if err == nil {
			parseFlag = "✅"

			if bytes.Equal(serialization.ToBytes(), data) {
				rebuildFlag = "✅"
			}
		}

		fmt.Printf("| %s | %s | %s | %s | %s |\n", extractName(name), extractPackage(name), parseFlag, rebuildFlag, duration)
	}

cleanup:
	os.Exit(exitCode)
}
