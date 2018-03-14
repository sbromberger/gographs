// Package bitvec is bit-vector with atomic access
package bitvec

import "sync/atomic"

// BBitVec is a bitvector
type BBitVec []uint32

// NewBBitVec returns a new bitvector with the given size
func NewBBitVec(size int) BBitVec {
	return make(BBitVec, uint(size+mask)>>nbits)
}

func (BBitVec) offset(k uint32) (bucket, bit uint32) {
	return k >> nbits, 1 << (k & mask)
}

func (bv BBitVec) TrySet(k uint32) bool {
	bucket, bit := bv.offset(k)
retry:
	old := atomic.LoadUint32(&bv[bucket])
	if old&bit != 0 {
		return false
	}
	if atomic.CompareAndSwapUint32(&bv[bucket], old, old|bit) {
		return true
	}
	goto retry
}
