package realAttempt

import (
	"fmt"
	"sort"
	"strconv"
)

func DumpRes(libs []*Library) [][]string {
	s := func(i int) string { return strconv.Itoa(i) }
	res := make([][]string, len(libs)*2+1)
	res[0] = []string{s(len(libs))}

	for k, v := range libs {
		res[k*2+1] = []string{s(v.Index), s(len(v.BooksOutput))}
		temp := make([]string, len(v.BooksOutput))
		for k2, v2 := range v.BooksOutput {
			temp[k2] = s(v2.Index)
		}
		res[k*2+2] = temp
	}
	return res
}

func (l *Library) Scoring() {
	l.Score = l.ParallelProcessing * len(l.Books) / l.SignupTime
}

type HandleLibs func([]*Library) []*Library

func BasicSolve(solve func([]*Library) []*Library) [][]string {
	var res []*Library
	tempLibs := solve(Libraries)

	for _, l := range tempLibs {
		add := l.Compute()
		if add {
			res = append(res, l)
		}
	}
	return DumpRes(res)
}

// SignupTime
func basicSignupTime(l []*Library) []*Library {
	sort.Slice(l, func(i, j int) bool {
		if l[i].SignupTime == l[j].SignupTime {
			return l[i].ParallelProcessing < l[j].ParallelProcessing
		}
		return l[i].SignupTime < l[j].SignupTime
	})
	return l
}

// parrallel * scoreBooks / signupTime
func basicScore(l []*Library) []*Library {
	for _, l := range Libraries {
		score := 0
		for _, b := range l.Books {
			score += b.Score
		}
		l.Score = score * l.ParallelProcessing / l.SignupTime
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i].Score > l[j].Score
	})
	return l
}

func ScoreSolve() [][]string {
	maxCount := 1000

	fmt.Printf("Start Solve%v\n", len(Libraries))
	var batch [][]*Library

	fmt.Printf("Permut %v\n", len(batch))
	score := 0
	var dump [][]string
	for i := 0; i < maxCount && i < len(batch); i++ {
		lib, tempScore := CalculateScore(batch[i])
		if tempScore > score {
			score = tempScore
			dump = DumpRes(lib)
		}
		fmt.Printf("Score ok %v %v\n", i, tempScore)
	}
	fmt.Printf("End Solve %v\n", score)
	return dump
}

// Perm calls f with each permutation of a.
func Perm(a []*Library, f func([]*Library)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []*Library, f func([]*Library), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

// Thx to https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go !
func Permutations(arr []*Library) [][]*Library {
	var helper func([]*Library, int)
	res := [][]*Library{}

	helper = func(arr []*Library, n int) {
		if n == 1 {
			tmp := make([]*Library, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
