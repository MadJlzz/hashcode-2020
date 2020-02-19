package internal

var SliceChildrenSize int
var SliceChildrenCanAdd = func(*ArrayScored, int) (bool, float64) { return true, 0 }

type ArrayScored struct {
	Data  []int
	Score float64
}

func (a *ArrayScored) Add(i int) bool {
	canAdd, newScore := SliceChildrenCanAdd(a, i)
	if canAdd {
		a.Data = append(a.Data, i)
		a.Score = newScore
	}
	return canAdd
}

func (a *ArrayScored) Clone() *ArrayScored {
	newData := make([]int, len(a.Data))
	copy(newData, a.Data)
	return &ArrayScored{newData, a.Score}
}

type SliceChildren struct {
	Base    *ArrayScored
	Current *ArrayScored

	ChildIndex   int
	CanHaveChild bool
}

func (s *SliceChildren) Complete() {
	var start int
	if s.Base.Data != nil && len(s.Base.Data) >= 0 {
		start = s.Base.Data[len(s.Base.Data)-1] + 1
	}

	for i := start; i < SliceChildrenSize; i++ {
		s.Current.Add(i)
	}

	if len(s.Base.Data) >= len(s.Current.Data) {
		s.CanHaveChild = false
		return
	}
	s.ChildIndex = s.Current.Data[len(s.Base.Data)]
	s.CanHaveChild = s.ChildIndex < SliceChildrenSize
}

func (s *SliceChildren) GetChild() *SliceChildren {
	var res *SliceChildren = nil
	temp := &SliceChildren{s.Base.Clone(), nil, 0, true}

	for j := s.ChildIndex + 1; j < SliceChildrenSize; j++ {
		isAdded := temp.Base.Add(j)
		if isAdded {
			temp.Current = temp.Base.Clone()
			s.ChildIndex = j
			s.CanHaveChild = s.ChildIndex < SliceChildrenSize
			return temp
		}
	}
	s.CanHaveChild = false
	return res
}
