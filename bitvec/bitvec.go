// Package bitvec is bit-vector with atomic access
package bitvec

import "sync/atomic"

const (
	nbits = 5
	ws    = 1 << nbits
	mask  = ws - 1
)

// BitVec is a bitvector
type BitVec []uint32

// NewBitVec returns a new bitvector with the given size
func NewBitVec(size int) BitVec {
	return make(BitVec, uint(size+mask)>>nbits)
}

func (BitVec) isBucketBitUnset(bucket uint32, k uint32) bool {
	return bucket&(1<<(k&mask)) == 0
}

func (BitVec) offset(k uint32) (bucket, bit uint32) {
	return k >> nbits, 1 << (k & mask)
}

func (bv BitVec) GetBucket(k uint32) uint32 {
	return atomic.LoadUint32(&bv[k>>nbits])
}

func (bv BitVec) GetBuckets4(a, b, c, d uint32) (x, y, z, w uint32) {
	x = atomic.LoadUint32(&bv[a>>nbits])
	y = atomic.LoadUint32(&bv[b>>nbits])
	z = atomic.LoadUint32(&bv[c>>nbits])
	w = atomic.LoadUint32(&bv[d>>nbits])
	return
}

func (bv BitVec) TrySet(k uint32) bool {
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

func (bv BitVec) TrySetWith(old uint32, k uint32) bool {
	bucket, bit := bv.offset(k)
	if old&bit != 0 {
		return false
	}
retry:
	if atomic.CompareAndSwapUint32(&bv[bucket], old, old|bit) {
		return true
	}
	old = atomic.LoadUint32(&bv[bucket])
	if old&bit != 0 {
		return false
	}
	goto retry
}
