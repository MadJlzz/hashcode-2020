package tools

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// thx to https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang
func ZipWriter(skipOutput bool, baseOutputDir string) {
	if skipOutput {
		return
	}

	baseFolder, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	baseFolder += "/"

	// Get a Buffer to Write To
	outFile, err := os.Create(fmt.Sprintf("%s/project.zip", baseOutputDir))
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}
