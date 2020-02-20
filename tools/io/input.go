package io

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const (
	InputEndOfLine = '\n'
	InputDelimiter = " "
)

func ReadInput(filename string) map[int][]string {
	res := make(map[int][]string)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error reading file input. got %s\n", err)
	}

	r := bufio.NewReader(f)
	for i := 0; ; i++ {
		l, err := r.ReadString(InputEndOfLine)
		if err == io.EOF {
			break
		}
		l = strings.ReplaceAll(l, "\n", "")
		l = strings.ReplaceAll(l, "\r", "")
		res[i] = strings.Split(l, InputDelimiter)
	}

	return res
}
