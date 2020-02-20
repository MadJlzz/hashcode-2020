package realAttempt

import "strconv"

var Days int
var Libraries []Library
var Books []int

type Library struct {
	SignupTime         int
	ParallelProcessing int
	BooksIndex         []int

	Score     int
	StartDate int
}

func NewLibrary(fileContent map[int][]string) {

	Books = make([]int, toInt(fileContent[0][0]))
	Libraries = make([]Library, toInt(fileContent[0][1]))
	Days = toInt(fileContent[0][2])

	for index, score := range fileContent[1] {
		Books[index] = toInt(score)
	}

	for i := 2; i < len(fileContent); i += 2 {
		currentLibrary := Library{
			BooksIndex:         make([]int, toInt(fileContent[i][0])),
			SignupTime:         toInt(fileContent[i][1]),
			ParallelProcessing: toInt(fileContent[i][2]),
		}
		for k, v := range fileContent[i+1] {
			currentLibrary.BooksIndex[k] = toInt(v)
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
