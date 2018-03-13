// Package bitvec is bit-vector with atomic access
package bitvec

import "sync/atomic"

// ABitVec is a bitvector
type ABitVec []uint64

// NewABitVec returns a new bitvector with the given size
func NewABitVec(size int) ABitVec {
	return make(ABitVec, uint(size+63)/64)
}

// Get returns the given bit
func (b ABitVec) Get(bit uint32) bool {
	shift := bit % 64
	bb := b[bit/64]
	bb &= (1 << shift)
	return bb != 0
}

// Set sets the given bit
func (b ABitVec) Set(bit uint32) {
	b[bit/64] |= (1 << (bit % 64))
}

// AGet atomically returns the given bit
func (b ABitVec) AGet(bit uint32) bool {
	shift := bit % 64
	bb := atomic.LoadUint64(&b[bit/64])
	bb &= (1 << shift)
	return bb != 0
}

// ASet atomically sets the given bit
func (b ABitVec) ASet(bit uint32) {
	set := uint64(1) << (bit % 64)
	addr := &b[bit/64]
	var old uint64
	for {
		old = atomic.LoadUint64(addr)
		if (old&set != 0) || atomic.CompareAndSwapUint64(addr, old, old|set) {
			break
		}
	}
}
