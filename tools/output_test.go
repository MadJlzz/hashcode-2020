package tools

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

func Check(t *testing.T, want interface{}, got interface{}) {
	if want != got {
		t.Errorf("Error: expected %v - got %v", want, got)
	}
}

func CheckExist(t *testing.T, path string, exist bool) {
	found := false
	if _, err := os.Stat(path); os.IsExist(err) {
		found = true
	}
	if (exist && !found) || (!exist && found) {
		t.Errorf("%s unexpected result -> should exist=%v", path, exist)
	}
}

func TestGetOutputDir(t *testing.T) {
	dummy := GetOutputDir(true, "test")
	Check(t, "dummy", dummy)

	path1 := GetOutputDir(false, "test")
	Check(t, "test/1", path1)

	path2 := GetOutputDir(false, "test")
	Check(t, "test/2", path2)

	_ = os.Remove(path1)
	_ = os.Remove(path2)
}

func TestGetTestFiles(t *testing.T) {
	_ = os.Mkdir("abc", 0777)
	f1, _ := os.Create("abc/test1.txt")
	f2, _ := os.Create("abc/test2.txt")
	f3, _ := os.Create("abc/test3.txt")
	_ = f1.Close()
	_ = f2.Close()
	_ = f3.Close()
	files := GetTestFiles("abc")
	Check(t, 3, len(files))
	Check(t, "abc/test1.txt", files[0])
	Check(t, "abc/test2.txt", files[1])
	Check(t, "abc/test3.txt", files[2])
	_ = os.RemoveAll("abc/")
}

func TestEmptyDumpToFile(t *testing.T) {
	empty := make([][]Serializable, 0)
	compare(t, &empty, "output_empty.txt", "test/output_empty.txt")
}

func TestDumpToFile(t *testing.T) {
	someItem := [][]Serializable{{dummy{1, "a"}, dummy{2, "b"}}, {dummy{3, "c"}}}
	compare(t, &someItem, "output_filled.txt", "test/output_filled.txt")
}

func TestDumpMapToFile(t *testing.T) {
	someItem := make(map[int][]string)
	someItem[0] = []string{"a=1-b=a", "a=2-b=b"}
	someItem[2] = []string{"a=3-b=c"}
	compareWithMap(t, &someItem, "output_filled.txt", "test/output_filled2.txt")
}

func TestDumpListToFile(t *testing.T) {
	someItem := make([][]string, 3)
	someItem[0] = []string{"a=1-b=a", "a=2-b=b"}
	someItem[2] = []string{"a=3-b=c"}
	compareWithList(t, &someItem, "output_filled.txt", "test/output_filled2.txt")
}

func compareWithMap(t *testing.T, input *map[int][]string, resFilePath string, expectedFilePath string) {
	DumpStringMapToFile(resFilePath, input)
	deepCompare(t, expectedFilePath, resFilePath)
	os.Remove(resFilePath)
}

func compareWithList(t *testing.T, input *[][]string, resFilePath string, expectedFilePath string) {
	DumpStringListToFile(resFilePath, input)
	deepCompare(t, expectedFilePath, resFilePath)
	os.Remove(resFilePath)
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
