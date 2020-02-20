package main

import (
	"flag"
	"fmt"
	"github.com/MadJlzz/hashcode-2020/realAttempt"
	"github.com/MadJlzz/hashcode-2020/tools/io"
	"path/filepath"
)

var (
	filename   = flag.String("filename", "", "the data file used for our algorithm")
	algorithm  = flag.String("solver", "", "the solver to use when trying to resolve the problem")
	skipOutput = flag.Bool("skipOutput", false, "No Output file")

	ResBasePath    = "realRes"
	TestDataFolder = "realTest"
)

// Define your solver here
var solver = realAttempt.SolveExercise

func main() {
	flag.Parse()
	resPath := io.GetOutputDir(*skipOutput, ResBasePath)

	if *filename == "" {
		testFiles := io.GetTestFiles(TestDataFolder)
		for _, testFile := range testFiles {
			handleFile(testFile, resPath, solver)
		}
		io.ZipWriter(*skipOutput, resPath)
	} else {
		handleFile(*filename, resPath, solver)
	}
}

func handleFile(filename string, resPath string, mySolver func(map[int][]string) [][]string) {
	fmt.Printf("*************\nTesting file: %s\n", filename)
	fileContent := io.ReadInput(filename)

	res := mySolver(fileContent)

	if *skipOutput {
		fmt.Printf("Result: \n%s\n*************\n", res)
	} else {
		fileOutput := fmt.Sprintf("%s/%s", resPath, filepath.Base(filename))
		io.DumpStringListToFile(fileOutput, &res)
		fmt.Printf("Result in file: %s\n*************\n", fileOutput)
	}
}
