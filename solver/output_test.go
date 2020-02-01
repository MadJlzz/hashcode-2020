package solver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

const chunkSize = 64000

type dummy struct {
	a int
	b string
}

func (d dummy) String() string {
	return fmt.Sprintf("a=%d-b=%s", d.a, d.b)
}

func TestEmptyDumpToFile(t *testing.T) {
	empty := make([][]Serializable, 0)
	compare(t, &empty, "output_empty.txt", "../test/output_empty.txt")
}

func TestDumpToFile(t *testing.T) {
	someItem := [][]Serializable{{dummy{1, "a"}, dummy{2, "b"}}, {dummy{3, "c"}}}
	compare(t, &someItem, "output_filled.txt", "../test/output_filled.txt")
}

func compare(t *testing.T, input *[][]Serializable, resFilePath string, expectedFilePath string) {
	DumpToFile(resFilePath, input)
	deepCompare(t, expectedFilePath, resFilePath)
	os.Remove(resFilePath)
}

func deepCompare(t *testing.T, expected, result string) {
	// Check file size ...

	f1, err := os.Open(expected)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(result)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return
			} else if err1 == io.EOF || err2 == io.EOF {
				return
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			t.Errorf("Unexpected file: expected=[%s], res=[%s]", b1, b2)
		}
	}
}
