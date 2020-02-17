package tools

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	OutputEndOfLine = "\n"
	OutputDelimiter = " "
)

type Serializable interface {
	String() string
}
type myString struct {
	o string
}

func (s myString) String() string { return s.o }

func GetOutputDir(skipOutput bool, resBasePath string) string {
	if skipOutput {
		return "dummy"
	}

	if _, err := os.Stat(resBasePath); os.IsNotExist(err) {
		_ = os.Mkdir(resBasePath, os.ModeDir)
	}

	files, err := ioutil.ReadDir(resBasePath)
	if err != nil {
		log.Fatal(err)
	}

	maxVersion := 0
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		version, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		if version > maxVersion {
			maxVersion = version
		}
	}

	res := fmt.Sprintf("%s/%d", resBasePath, maxVersion+1)
	_ = os.Mkdir(res, os.ModeDir)
	return res
}

func GetTestFiles(testDataFolder string) (res []string) {
	files, err := ioutil.ReadDir(testDataFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		res = append(res, fmt.Sprintf("%s/%s", testDataFolder, file.Name()))
	}
	return res
}

func DumpStringMapToFile(filename string, outputMap *map[int][]string) {
	var output [][]Serializable
	var maxLine int

	for n := range *outputMap {
		if n > maxLine {
			maxLine = n
		}
	}

	for i := 0; i <= maxLine; i++ {
		line, exist := (*outputMap)[i]
		newLine := []Serializable{}
		if exist {
			for j := 0; j < len(line); j++ {
				newLine = append(newLine, myString{line[j]})
			}
		}
		output = append(output, newLine)
	}
	DumpToFile(filename, &output)
}

func DumpStringListToFile(filename string, outputList *[][]string) {
	output := make([][]Serializable, len(*outputList))
	for i, v := range *outputList {
		output[i] = make([]Serializable, len(v))
		for j, w := range v {
			output[i][j] = myString{w}
		}
	}
	DumpToFile(filename, &output)
}

func DumpToFile(filename string, output *[][]Serializable) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error reading file input. got %s\n", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for i, line := range *output {
		if i > 0 {
			write(w, OutputEndOfLine)
		}

		for j, item := range line {
			if j > 0 {
				write(w, OutputDelimiter)
			}
			write(w, item.String())
		}
	}
}

func write(w *bufio.Writer, s string) {
	_, err := w.WriteString(s)
	if err != nil {
		log.Fatalf("error writing [%s] into file. got %s\n", s, err)
	}
}
