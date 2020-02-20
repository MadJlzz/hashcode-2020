package realAttempt

import (
	"sort"
	"strconv"
)

var Days int
var Libraries []*Library
var Books []*Book

type Library struct {
	Index              int
	SignupTime         int
	SignupTimeTemp     int
	ParallelProcessing int
	Books              []*Book

	// Books that we want to output for the current library
	BooksOutput []*Book

	Score     int
	StartDate int
}

type Book struct {
	Index int
	Score int
	Taken bool
}

func NewLibrary(fileContent map[int][]string) {

	Books = make([]*Book, toInt(fileContent[0][0]))
	Libraries = make([]*Library, toInt(fileContent[0][1]))
	Days = toInt(fileContent[0][2])

	for index, score := range fileContent[1] {
		Books[index] = &Book{
			Index: index,
			Score: toInt(score),
			Taken: false,
		}
	}

	for i := 2; i < len(fileContent); i += 2 {
		if fileContent[i][0] == "" {
			continue
		}
		idx := i/2 - 1
		currentLibrary := &Library{
			Books:              make([]*Book, toInt(fileContent[i][0])),
			Index:              idx,
			SignupTime:         toInt(fileContent[i][1]),
			SignupTimeTemp:     toInt(fileContent[i][1]),
			ParallelProcessing: toInt(fileContent[i][2]),
		}
		for k, v := range fileContent[i+1] {
			currentLibrary.Books[k] = Books[toInt(v)]
		}

		sort.Slice(currentLibrary.Books, func(i, j int) bool { return currentLibrary.Books[i].Score > currentLibrary.Books[j].Score })
		Libraries[idx] = currentLibrary
	}

}

func toInt(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return res
}

func (l *Library) Compute() bool {
	sort.Slice(l.Books, func(i, j int) bool { return l.Books[i].Score > l.Books[j].Score })
	for _, b := range l.Books {
		if !b.Taken {
			b.Taken = true
			l.BooksOutput = append(l.BooksOutput, b)
		}
	}
	return len(l.BooksOutput) > 0
}

func Reinit(libraries []*Library) {
	for _, l := range libraries {
		l.SignupTimeTemp = l.SignupTime
		l.BooksOutput = nil
	}
}

func CalculateScore(libraries []*Library) ([]*Library, int) {
	var res []*Library
	var lastLib *Library
	lastLibCurrent := 0

	score := 0
	for i := 0; i < Days; i++ {
		if lastLib != nil && lastLib.SignupTimeTemp == 0 {
			lastLibCurrent++
			res = append(res, lastLib)
			lastLib = nil
		}
		if lastLib == nil && lastLibCurrent < len(libraries) {
			lastLib = libraries[lastLibCurrent]
		}
		if lastLib != nil {
			lastLib.SignupTimeTemp--
		}

		for _, l := range res {
			maxNb := l.ParallelProcessing
			for _, b := range l.Books {
				if !b.Taken {
					l.BooksOutput = append(l.BooksOutput, b)
					b.Taken = true
					maxNb--
					score += b.Score
					if maxNb == 0 {
						break
					}
				}
			}
		}
	}
	return res, score
}
