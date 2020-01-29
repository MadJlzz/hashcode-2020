package solver

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const (
	inputEndOfLine = '\n'
	inputDelimiter = " "
)

func ReadInput(filename string) map[int][]string {
	res := make(map[int][]string)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error reading file input. got %s\n", err)
	}

	r := bufio.NewReader(f)
	for i := 0; ; i++ {
		l, err := r.ReadString(inputEndOfLine)
		if err == io.EOF {
			break
		}
		res[i] = strings.Split(l, inputDelimiter)
	}

	return res
}
