package mmf

import (
	"os"
	"reflect"
	"unsafe"

	"github.com/tysontate/gommap"
)

// File represents state of the memory mapped file
type File struct {
	f    *os.File
	MMap gommap.MMap
}

// Open opens the memory mapped file
func Open(fileName string) (*File, error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_SYNC|os.O_RDWR, 0700)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	file := &File{
		f: f,
	}

	if stat.Size() != 0 {
		err = file.mmap()
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (f *File) mmap() error {
	mmap, err := gommap.Map(f.f.Fd(), gommap.PROT_READ, gommap.MAP_SHARED)
	if err != nil {
		return err
	}
	err = mmap.Advise(gommap.MADV_RANDOM)
	if err != nil {
		return err
	}
	f.MMap = mmap
	return nil
}

// Append appends the data at the end of the file and extens the boundaries of the mmap
func (f *File) Append(data []byte) error {

	shouldMMAP := len(f.MMap) == 0

	n, err := f.f.Write(data)

	if err != nil {
		return err
	}

	newLength := len(f.MMap) + n

	// really uggly, but not to be avoided :(
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&f.MMap))
	dh.Len = newLength
	dh.Cap = newLength

	if shouldMMAP {
		err = f.mmap()
	}

	return err

}
