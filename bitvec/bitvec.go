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
var WS = 64
var NBITS = uint(6)
var ONES = uint(WS - 1)
var SZ = 1

type BitVec []int

func NewBitVec(n int) BitVec {
	if n > WS {
		SZ = (n / WS) + 1
	}
	log.Printf("Bit vector of base size %d (%d max bits)\n", SZ, WS*SZ)
	return make([]int, SZ, SZ)
}

func (bv BitVec) IsSet(k int) bool {
	// return (bv[k/WS] & (1 << (uint(k % WS)))) != 0
	return bv[k>>NBITS]&(1<<(uint(k)&ONES)) != 0
}

func (bv BitVec) Set(k int) {
	bv[k>>NBITS] |= (1 << uint(uint(k)&ONES))
	// bv[k/WS] |= (1 << uint(k%WS))
}

func (bv BitVec) Clear(k int) {
	bv[k>>NBITS] &= ^(1 << uint(uint(k)&ONES))
	// bv[k/WS] &= ^(1 << uint(k%WS))
}
