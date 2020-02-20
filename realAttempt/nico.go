package realAttempt

import "strconv"

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
