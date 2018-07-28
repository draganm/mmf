package mmf

import (
	"errors"
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
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return nil, err
	}

	file := &File{
		f: f,
	}

	err = file.mmap()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *File) mmap() error {
	stat, err := f.f.Stat()
	if err != nil {
		return err
	}

	nc, err := newCapacity(stat.Size())
	if err != nil {
		return err
	}

	mmap, err := gommap.MapAt(0, f.f.Fd(), 0, nc, gommap.PROT_READ, gommap.MAP_SHARED)
	if err != nil {
		return err
	}

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&mmap))
	dh.Len = int(stat.Size())

	err = mmap.Advise(gommap.MADV_RANDOM)
	if err != nil {
		return err
	}

	f.MMap = mmap
	return nil
}

// for 64bit linux
const maxMapSize = 0xFFFFFFFFFFFF

const maxMMAPStep = 1 << 30

func newCapacity(min int64) (int64, error) {
	if min < 1<<10 {
		return 1 << 10, nil
	}
	for i := uint(10); i <= 30; i++ {
		if min <= 1<<i {
			return 1 << i, nil
		}
	}

	numPages := maxMMAPStep/uint64(min) + 1

	size := numPages * maxMMAPStep
	if size > maxMapSize {
		return 0, errors.New("Too large mmap")
	}

	return int64(size), nil
}

// Append appends the data at the end of the file and extens the boundaries of the mmap
func (f *File) Append(data []byte) error {

	n, err := f.f.Write(data)

	if err != nil {
		return err
	}

	newLength := len(f.MMap) + n

	if newLength > cap(f.MMap) {
		if len(f.MMap) > 0 {
			err = f.MMap.UnsafeUnmap()
			if err != nil {
				return err
			}

			f.MMap = nil
		}

		err = f.mmap()
		if err != nil {
			return err
		}
		return nil
	}

	// really uggly, but not to be avoided :(
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&f.MMap))
	dh.Len = newLength

	return err

}

// Close closes unmmaps and closes the file
func (f *File) Close() error {

	if len(f.MMap) > 0 {
		err := f.MMap.UnsafeUnmap()
		if err != nil {
			return err
		}
	}

	err := f.f.Close()
	if err != nil {
		return err
	}

	f.f = nil
	f.MMap = nil

	return nil
}

func (f *File) Empty() error {
	if len(f.MMap) == 0 {
		return nil
	}

	err := f.MMap.UnsafeUnmap()
	if err != nil {
		return err
	}

	err = f.f.Truncate(0)
	if err != nil {
		return err
	}

	f.MMap = nil
	return nil
}
