package sparsevecs

import (
	"sort"
)

//UInt32SparseVec holds a row index and column pointer. If a point is defined
//at a particular row i and column j, an edge exists between vertex i and vertex j.
type UInt32SparseVec struct {
	Rowidx []uint32
	Colptr []uint64 // indexes into rowidx - must be uint64
}

// inserts value x into sorted vector v, preserving ordering.
func insertSort(x uint32, v []uint32) []uint32 {
	index := sort.Search(len(v), func(i int) bool { return v[i] > x })
	v = append(v, uint32(0))
	copy(v[index+1:], v[index:])
	v[index] = x
	return v
}

// find value x in sorted vector v. Return index and true/false indicating found.
func searchsorted32(x uint32, v []uint32) (uint32, bool) {
	s := sort.Search(len(v), func(i int) bool { return v[i] >= x })
	// fmt.Println("s = ", s)
	found := (s < len(v)) && (v[s] == x)

	return uint32(s), found
}

func searchsorted64(x uint64, v []uint64) (uint64, bool) {
	s := sort.Search(len(v), func(i int) bool { return v[i] >= x })
	found := (s < len(v)) && (v[s] == x)

	return uint64(s), found
}

// GetIndex returns true if the value at (r, c) is defined.
func (v UInt32SparseVec) GetIndex(r, c uint32) bool {
	if len(v.Colptr) <= int(c)+1 {
		return false
	}

	_, found := searchsorted32(r, v.GetRange(c))
	return found
}

// GetIndexInt returns true if the value at (r, c) is defined.
func (v UInt32SparseVec) GetIndexInt(x, y int) bool {
	return v.GetIndex(uint32(x), uint32(y))
}

// GetRange returns the row slice for a given column
func (v UInt32SparseVec) GetRange(c uint32) []uint32 {
	p1 := v.Colptr[c]
	p2 := v.Colptr[c+1]
	return v.Rowidx[p1:p2]
}

// GetRangeInt returns the row slice for a given column expressed as an integer
func (v UInt32SparseVec) GetRangeInt(c int) []uint32 {
	return v.GetRange(uint32(c))
}
