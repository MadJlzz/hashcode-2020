package realAttempt

import (
	"fmt"
	"sort"
)

type Res struct {
	name   string
	score  int
	output [][]string
}

func SolveExercise(fileContent map[int][]string) (res [][]string) {

	NewLibrary(fileContent)

	var resList []Res

	resList = append(resList, test("basicScore1", Scoring1))
	resList = append(resList, test("basicScore2", Scoring2))
	resList = append(resList, test("basicScore3", Scoring3))
	resList = append(resList, test("basicScore4", Scoring4))

	sort.Slice(resList, func(i, j int) bool { return resList[i].score > resList[j].score })
	max := resList[0]

	fmt.Printf("Best score: %d with %s", max.score, max.name)

	return max.output
}

func Scoring1(bookScore int, l *Library) int {
	return bookScore * l.ParallelProcessing / (l.SignupTime * l.SignupTime)
}
func Scoring2(bookScore int, l *Library) int {
	return bookScore * bookScore * l.ParallelProcessing / (l.SignupTime * l.SignupTime)
}
func Scoring3(bookScore int, l *Library) int { return bookScore / (l.SignupTime * l.SignupTime) }
func Scoring4(bookScore int, l *Library) int {
	return bookScore * len(l.Books) / (l.SignupTime * l.SignupTime)
}

func test(name string, attempt func(bookScore int, l *Library) int) Res {
	l := basicScore(Libraries, attempt)
	temp, score := CalculateScore(l, attempt)
	res := DumpRes(temp)
	fmt.Printf("Attempt %s score %d\n", name, score)
	Reinit()
	return Res{name, score, res}
}
