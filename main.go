package main

import (
	"flag"
	"fmt"
	"github.com/madjlzz/hashcode-2020/exercise"
	"github.com/madjlzz/hashcode-2020/tools"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	filename   = flag.String("filename", "", "the data file used for our algorithm")
	algorithm  = flag.String("solver", "", "the solver to use when trying to resolve the problem")
	skipOutput = flag.Bool("skipOutput", false, "No Output file")
)

const resBasePath = "res"
const testDataFolder = "test"

// Define your solver here
var solver = exercise.SolveExercise

func main() {
	flag.Parse()
	resPath := getOutputDir()

	if *filename == "" {
		testFiles := getTestFiles()
		for _, testFile := range testFiles {
			handleFile(testFile, resPath, solver)
		}
		tools.ZipWriter(resPath)
	} else {
		handleFile(*filename, resPath, solver)
	}
}

func handleFile(filename string, resPath string, mySolver func(map[int][]string) [][]string) {
	fmt.Printf("*************\nTesting file: %s\n", filename)
	fileContent := tools.ReadInput(filename)

	res := mySolver(fileContent)

	if !*skipOutput {
		fileOutput := fmt.Sprintf("%s/%s", resPath, filepath.Base(filename))
		tools.DumpStringListToFile(fileOutput, &res)
		fmt.Printf("Result in file: %s\n*************\n", fileOutput)
	}
}

func getOutputDir() string {
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

	if !*skipOutput {
		_ = os.Mkdir(res, os.ModeDir)
	}
	return res
}

func getTestFiles() (res []string) {
	files, err := ioutil.ReadDir(testDataFolder)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(files); i++ {
		res = append(res, fmt.Sprintf("%s/%s", testDataFolder, files[i].Name()))
	}
	return res
}
