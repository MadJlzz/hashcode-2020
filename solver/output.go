package solver

import (
	"bufio"
	"log"
	"os"
)

const (
	OutputEndOfLine = "\n"
	OutputDelimiter = " "
)

type Serializable interface {
	String() string
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
