package graph

import (
	"os"
	"unsafe"

	mmap "github.com/edsrzf/mmap-go"
)

// raw defines a struct for binary SimpleGraphs format.
type raw struct {
	file *os.File
	data mmap.MMap

	FIndicesLen uint64
	FIndPtrLen  uint64
	BIndicesLen uint64
	BIndPtrLen  uint64

	FIndPtr  []uint64
	FIndices []uint32

	BIndPtr  []uint64
	BIndices []uint32
}

func (raw *raw) close() error {
	var err1, err2 error
	if raw.data != nil {
		err1 = raw.data.Unmap()
	}
	if raw.file != nil {
		err2 = raw.file.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}

func loadRaw(filename string) (*raw, error) {
	var err error
	raw := &raw{}

	raw.file, err = os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		raw.close()
		return nil, err
	}

	raw.data, err = mmap.Map(raw.file, mmap.RDONLY, 0)
	if err != nil {
		raw.close()
		return nil, err
	}

	x := 0
	copy((*[8]byte)(unsafe.Pointer(&raw.FIndPtrLen))[:], raw.data[x:x+8])
	x += 8
	copy((*[8]byte)(unsafe.Pointer(&raw.FIndicesLen))[:], raw.data[x:x+8])
	x += 8

	copy((*[8]byte)(unsafe.Pointer(&raw.BIndPtrLen))[:], raw.data[x:x+8])
	x += 8
	copy((*[8]byte)(unsafe.Pointer(&raw.BIndicesLen))[:], raw.data[x:x+8])
	x += 8

	raw.FIndPtr = ((*[1 << 40]uint64)(unsafe.Pointer(&raw.data[x])))[0:int(raw.FIndPtrLen)]
	x += 8 * int(raw.FIndPtrLen)
	raw.FIndices = ((*[1 << 40]uint32)(unsafe.Pointer(&raw.data[x])))[0:int(raw.FIndicesLen)]
	x += 4 * int(raw.FIndicesLen)

	raw.BIndPtr = ((*[1 << 40]uint64)(unsafe.Pointer(&raw.data[x])))[0:int(raw.BIndPtrLen)]
	x += 8 * int(raw.BIndPtrLen)
	raw.BIndices = ((*[1 << 40]uint32)(unsafe.Pointer(&raw.data[x])))[0:int(raw.BIndicesLen)]

	return raw, nil
}

// Save saves a SimpleGraph to a file in raw (binary) format.
func (g SimpleGraph) Save(filename string) error {
	FIndPtr := g.FMat().IndPtr
	FIndices := g.FMat().Indices
	BIndPtr := g.BMat().IndPtr
	BIndices := g.BMat().Indices

	output, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	FIndPtrLen := int64(len(FIndPtr))
	FIndicesLen := int64(len(FIndices))
	BIndPtrLen := int64(len(BIndPtr))
	BIndicesLen := int64(len(BIndices))

	FIndPtrBytes := 8 * len(FIndPtr)
	FIndicesBytes := 4 * len(FIndices)
	BIndPtrBytes := 8 * len(BIndPtr)
	BIndicesBytes := 4 * len(BIndices)

	err = output.Truncate(int64(8 + 8 + 8 + 8 + FIndPtrBytes + FIndicesBytes + BIndPtrBytes + BIndicesBytes))
	if err != nil {
		return err
	}

	data, err := mmap.Map(output, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	defer data.Unmap()

	x := 0

	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&FIndPtrLen))[:]))
	x += 8
	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&FIndicesLen))[:]))
	x += 8

	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&BIndPtrLen))[:]))
	x += 8
	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&BIndicesLen))[:]))
	x += 8

	if len(FIndPtr) > 0 {
		copy(data[x:x+FIndPtrBytes],
			((*[1 << 40]byte)(unsafe.Pointer(&FIndPtr[0]))[:FIndPtrBytes]))
		x += FIndPtrBytes
	}
	if len(FIndices) > 0 {
		copy(data[x:x+FIndicesBytes],
			((*[1 << 40]byte)(unsafe.Pointer(&FIndices[0]))[:FIndicesBytes]))
		x += FIndicesBytes
	}

	if len(FIndPtr) > 0 {
		copy(data[x:x+BIndPtrBytes],
			((*[1 << 40]byte)(unsafe.Pointer(&BIndPtr[0]))[:BIndPtrBytes]))
		x += BIndPtrBytes
	}
	if len(BIndices) > 0 {
		copy(data[x:x+BIndicesBytes],
			((*[1 << 40]byte)(unsafe.Pointer(&BIndices[0]))[:BIndicesBytes]))
		x += BIndicesBytes
	}

	return nil
}

// Load loads a raw (binary) SimpleGraph file and returns a SimpleGraph.
func Load(fn string) (SimpleGraph, error) {
	raw, err := loadRaw(fn)
	if err != nil {
		return SimpleGraph{}, err
	}

	find := make([]uint32, raw.FIndicesLen)
	findptr := make([]uint64, raw.FIndPtrLen)
	bind := make([]uint32, raw.BIndicesLen)
	bindptr := make([]uint64, raw.BIndPtrLen)
	copy(find, raw.FIndices)
	copy(findptr, raw.FIndPtr)
	copy(bind, raw.BIndices)
	copy(bindptr, raw.BIndPtr)
	return FromRaw(findptr, find, bindptr, bind)
}
