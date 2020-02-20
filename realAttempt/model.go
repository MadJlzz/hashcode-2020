package realAttempt

import "strconv"

var Days int
var Libraries []Library
var Books []*Book

type Library struct {
	SignupTime         int
	ParallelProcessing int
	Books              []*Book

	// Books that we want to output for the current library
	BooksOutput []int

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
	Libraries = make([]Library, toInt(fileContent[0][1]))
	Days = toInt(fileContent[0][2])

	for index, score := range fileContent[1] {
		Books[index] = &Book{
			Index: index,
			Score: toInt(score),
			Taken: false,
		}
	}

	for i := 2; i < len(fileContent); i += 2 {
		currentLibrary := Library{
			Books:              make([]*Book, toInt(fileContent[i][0])),
			SignupTime:         toInt(fileContent[i][1]),
			ParallelProcessing: toInt(fileContent[i][2]),
		}
		for k, v := range fileContent[i+1] {
			currentLibrary.Books[k] = Books[toInt(v)]
		}
		Libraries[i/2-1] = currentLibrary
	}

}

func toInt(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return res
}
