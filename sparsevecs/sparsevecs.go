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

// Insert will insert (r, c) into the sparse vector, or will
// set err if it already exists.
// func (v *UInt32SparseVec) Insert(r, c uint32) error {
// 	if v.GetIndex(r, c) {
// 		fmt.Printf("Insert: %d, %d already exists\n", r, c)
// 		return nil
// 	}
// 	// expand colptr if the column is larger than what we have now.
// 	lencol := uint32(len(v.Colptr))
// 	lastval := v.Colptr[lencol-1]
// 	fmt.Println("Insert: r = ", r)
// 	fmt.Println("Insert: v.Colptr = ", v.Colptr)

// 	if c >= uint32(lencol-1) {
// 		fmt.Println("Insert: padding")
// 		filler := make([]uint32, c-(lencol-1))
// 		for i := range filler {
// 			filler[i] = lastval
// 		}
// 		v.Colptr = append(v.Colptr, filler...)
// 		fmt.Println("appending last el")
// 		v.Colptr = append(v.Colptr, uint32(lastval+1))
// 	}

// 	// find the slice of rowidx that represents the column.
// 	fmt.Println("Insert: 2: r = ", r)
// 	fmt.Println("Insert: 2: v.Colptr = ", v.Colptr)
// 	p1 := v.Colptr[r]
// 	p2 := v.Colptr[r+1]
// 	fmt.Printf("Insert: p1 = %d, p2 = %d\n", p1, p2)
// 	x := v.Rowidx[:p1]
// 	y := insertSort(c, v.Rowidx[p1:p2])
// 	z := v.Rowidx[p2:]
// 	v.Rowidx = append(append(x, y...), z...)
// 	fmt.Println("Insert: v.Colptr = ", v.Colptr)
// 	fmt.Printf("Insert: now incrementing from v.Colptr[%d]\n", r)
// 	for i := r; i < uint32(len(v.Colptr)); i++ {
// 		(v.Colptr[i])++
// 	}

// 	return nil
// }

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
	// fmt.Println("s = ", s)
	found := (s < len(v)) && (v[s] == x)

	return uint64(s), found
}

// GetIndex returns true if the value at (r, c) is defined.
func (v UInt32SparseVec) GetIndex(r, c uint32) bool {
	// fmt.Printf("GetIndex: r = %d, c = %d\n", r, c)
	// fmt.Println("GetIndex: v.Colptr =", v.Colptr)
	// fmt.Println("GetIndex: len(v.Colptr) =", len(v.Colptr))
	// fmt.Println("GetIndex: lencheck = ", len(v.Colptr) <= int(c)+1)
	if len(v.Colptr) <= int(c)+1 {
		return false
	}

	_, found := searchsorted32(r, v.GetRange(c))
	// fmt.Println("GetIndex: colptr = ", v.Colptr)
	// fmt.Println("GetIndex: rowidx = ", v.Rowidx)
	// fmt.Println("GetIndex: found = ", found)
	return found
}

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
