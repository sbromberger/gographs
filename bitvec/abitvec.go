package bitvec

import "sync/atomic"

// ABitVec is a bitvector
type ABitVec []uint32

// NewABitVec returns a new bitvector with the given size
func NewABitVec(size int) ABitVec {
	return make(ABitVec, uint(size+mask)>>nbits)
}

func (ABitVec) isBucketBitUnset(bucket uint32, k uint32) bool {
	return bucket&(1<<(k&mask)) == 0
}

func (ABitVec) offset(k uint32) (bucket, bit uint32) {
	return k >> nbits, 1 << (k & mask)
}

// GetBucket returns the int32 bucket
func (bv ABitVec) GetBucket(k uint32) uint32 {
	return atomic.LoadUint32(&bv[k>>nbits])
}

// GetBuckets4 returns buckets 4 at a time.
func (bv ABitVec) GetBuckets4(a, b, c, d uint32) (x, y, z, w uint32) {
	x = atomic.LoadUint32(&bv[a>>nbits])
	y = atomic.LoadUint32(&bv[b>>nbits])
	z = atomic.LoadUint32(&bv[c>>nbits])
	w = atomic.LoadUint32(&bv[d>>nbits])
	return
}

// TrySet will set the bit located at `k` and will
// return success
func (bv ABitVec) TrySet(k uint32) bool {
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

// TrySetWith ...
func (bv ABitVec) TrySetWith(old uint32, k uint32) bool {
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
