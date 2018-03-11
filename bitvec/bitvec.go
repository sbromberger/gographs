package bitvec

import (
	"log"
)

// this is the word size of an int(32)
// in bits. an int(32) requires a min-
// imum of 4 bytes, each of which are
// made up of 8 bits, therefore 4x8=32
// this same notion applies for int(64)
// such that 8 bytes * 8 bits/byte = 64
//
const (
	ws    = 32
	nbits = 5
	ones  = ws - 1
)

type BitVec []uint32

func NewBitVec(n int) BitVec {
	sz := 1
	if n > ws {
		sz = (n / ws) + 1
	}
	log.Printf("Bit vector of base size %d (%d max bits)\n", sz, ws*sz)
	return make([]uint32, sz, sz)
}

func (bv BitVec) IsSet(k uint32) bool {
	// return (bv[k/ws] & (1 << (uint(k % ws)))) != 0
	return bv[k>>nbits]&(1<<(k&ones)) != 0
}

func (bv BitVec) Set(k uint32) {
	bv[k>>nbits] |= (1 << (k & ones))
	// bv[k/ws] |= (1 << uint(k%ws))
}

func (bv BitVec) Clear(k uint32) {
	bv[k>>nbits] &= ^(1 << (k & ones))
	// bv[k/ws] &= ^(1 << uint(k%ws))
}
