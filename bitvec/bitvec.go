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
	nbits = 5
	ws    = 1 << nbits
	mask  = ws - 1
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

func (bv BitVec) offset(k uint32) (bucket, bit uint32) {
	return k >> nbits, 1 << (k & mask)
}

func (bv BitVec) TrySet(k uint32) bool {
	bucket, bit := bv.offset(k)
	unset := bv[bucket]&bit == 0
	bv[bucket] |= bit
	return unset
}

func (bv BitVec) IsSet(k uint32) bool {
	bucket, bit := bv.offset(k)
	return bv[bucket]&bit != 0
}

func (bv BitVec) Set(k uint32) {
	bucket, bit := bv.offset(k)
	bv[bucket] |= bit
}

func (bv BitVec) Clear(k uint32) {
	bucket, bit := bv.offset(k)
	bv[bucket] &= ^bit
}
