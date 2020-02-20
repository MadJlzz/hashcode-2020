package realAttempt

import "fmt"

func SolveExercise(fileContent map[int][]string) (res [][]string) {

	NewLibrary(fileContent)

	fmt.Println(Libraries)
	fmt.Println(Days)
	fmt.Println(Books)

	return nil
}
