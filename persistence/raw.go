package persistence

import (
	"log"
	"os"
	"unsafe"

	mmap "github.com/edsrzf/mmap-go"
	"github.com/sbromberger/gographs"
)

type Raw struct {
	file *os.File
	data mmap.MMap

	rowidxlen uint64
	colptrlen uint64

	rowidx []uint32
	colptr []uint64
}

func (raw *Raw) Rowidx() []uint32 { return raw.rowidx }
func (raw *Raw) Colptr() []uint64 { return raw.colptr }

func SaveRaw(filename string, g gographs.Graph) error {
	spVecPtr := g.ToSparseVec()
	rowidx := spVecPtr.Rowidx
	colptr := spVecPtr.Colptr
	output, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	rowidxlen := int64(len(rowidx))
	colptrlen := int64(len(colptr))

	rowidxbytes := 4 * len(rowidx)
	colptrbytes := 8 * len(colptr)

	err = output.Truncate(int64(8 + 8 + rowidxbytes + colptrbytes))
	if err != nil {
		return err
	}

	data, err := mmap.Map(output, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	defer data.Unmap()

	x := 0

	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&rowidxlen))[:]))
	x += 8

	copy(data[x:x+8], ((*[8]byte)(unsafe.Pointer(&colptrlen))[:]))
	x += 8

	if len(rowidx) > 0 {
		copy(data[x:x+rowidxbytes],
			((*[1 << 40]byte)(unsafe.Pointer(&rowidx[0]))[:rowidxbytes]))
		x += rowidxbytes
	}

	if len(colptr) > 0 {
		copy(data[x:x+colptrbytes],
			((*[1 << 40]byte)(unsafe.Pointer(&colptr[0]))[:colptrbytes]))
		x += colptrbytes
	}

	return nil
}

func LoadRaw(filename string) (*Raw, error) {
	var err error
	raw := &Raw{}

	raw.file, err = os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		raw.Close()
		return nil, err
	}

	raw.data, err = mmap.Map(raw.file, mmap.RDONLY, 0)
	if err != nil {
		raw.Close()
		return nil, err
	}

	x := 0
	copy((*[8]byte)(unsafe.Pointer(&raw.rowidxlen))[:], raw.data[x:x+8])
	x += 8

	copy((*[8]byte)(unsafe.Pointer(&raw.colptrlen))[:], raw.data[x:x+8])
	x += 8

	raw.rowidx = ((*[1 << 40]uint32)(unsafe.Pointer(&raw.data[x])))[0:int(raw.rowidxlen)]
	x += 4 * int(raw.rowidxlen)

	raw.colptr = ((*[1 << 40]uint64)(unsafe.Pointer(&raw.data[x])))[0:int(raw.colptrlen)]

	return raw, nil
}

func GraphFromRaw(fn string) gographs.Graph {
	raw, err := LoadRaw(fn)
	if err != nil {
		log.Fatal("error: ", err)
	}

	ri := make([]uint32, raw.rowidxlen)
	cp := make([]uint64, raw.colptrlen)
	copy(ri, raw.rowidx)
	copy(cp, raw.colptr)
	return gographs.MakeGraph(ri, cp)
}

func (raw *Raw) Close() error {
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
