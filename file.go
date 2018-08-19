package mmf

import (
	"errors"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

// File represents state of the memory mapped file
type File struct {
	f            *os.File
	MMap         []byte
	originalData uintptr
}

// Open opens the memory mapped file
func Open(fileName string) (*File, error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return nil, err
	}

	file := &File{
		f:    f,
		MMap: []byte{},
	}

	err = file.mmap()
	if err != nil {
		return nil, err
	}

	return file, nil
}

const (
	PROT_NONE  uint = 0x0
	PROT_READ  uint = 0x1
	PROT_WRITE uint = 0x2
	PROT_EXEC  uint = 0x4

	MAP_SHARED    uint = 0x1
	MAP_PRIVATE   uint = 0x2
	MAP_FIXED     uint = 0x10
	MAP_ANONYMOUS uint = 0x20
	MAP_GROWSDOWN uint = 0x100
	MAP_LOCKED    uint = 0x2000
	MAP_NONBLOCK  uint = 0x10000
	MAP_NORESERVE uint = 0x4000
	MAP_POPULATE  uint = 0x8000

	MADV_NORMAL     uint = 0x0
	MADV_RANDOM     uint = 0x1
	MADV_SEQUENTIAL uint = 0x2
	MADV_WILLNEED   uint = 0x3
	MADV_DONTNEED   uint = 0x4
	MADV_REMOVE     uint = 0x9
	MADV_DONTFORK   uint = 0xa
	MADV_DOFORK     uint = 0xb

	MREMAP_MAYMOVE uint = 0x1
)

func (f *File) mmap() error {
	stat, err := f.f.Stat()
	if err != nil {
		return err
	}

	nc, err := newCapacity(stat.Size())
	if err != nil {
		return err
	}

	addr, err := mmap_syscall(uintptr(int64(0)), uintptr(nc), uintptr(PROT_READ), uintptr(MAP_SHARED), f.f.Fd(), 0)
	if err != syscall.Errno(0) {
		return err
	}

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&(f.MMap)))

	f.originalData = dh.Data
	dh.Data = addr

	dh.Len = int(stat.Size())
	dh.Cap = int(nc)

	_, _, err = syscall.Syscall(syscall.SYS_MADVISE, uintptr(dh.Data), uintptr(dh.Len), uintptr(MADV_RANDOM))
	if err != syscall.Errno(0) {
		return err
	}

	return nil
}

func mmap_syscall(addr, length, prot, flags, fd uintptr, offset int64) (uintptr, error) {
	addr, _, err := syscall.Syscall6(syscall.SYS_MMAP, addr, length, prot, flags, fd, uintptr(offset))
	return addr, err
}

// for 64bit linux
const maxMapSize = 0xFFFFFFFFFFFFFFFF

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

	numPages := uint64(min)/maxMMAPStep + 1

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

		dh := (*reflect.SliceHeader)(unsafe.Pointer(&(f.MMap)))

		// MREMAP_MAYMOVE

		var nc int64
		nc, err = newCapacity(int64(newLength))
		if err != nil {
			return err
		}

		addr, _, err := syscall.Syscall6(syscall.SYS_MREMAP, uintptr(dh.Data), uintptr(dh.Len), uintptr(nc), uintptr(MREMAP_MAYMOVE), uintptr(0), uintptr(0))
		if err != syscall.Errno(0) {
			return err
		}

		dh.Data = addr
		dh.Cap = int(nc)
		dh.Len = newLength
		return nil
	}

	// really uggly, but not to be avoided :(
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&f.MMap))
	dh.Len = newLength

	return err

}

// Close closes unmmaps and closes the file
func (f *File) Close() error {

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&(f.MMap)))

	_, _, err := syscall.Syscall(syscall.SYS_MUNMAP, uintptr(dh.Data), uintptr(dh.Cap), uintptr(0))
	if err != syscall.Errno(0) {
		return err
	}

	dh.Data = f.originalData
	dh.Cap = 0
	dh.Len = 0

	e := f.f.Close()
	if e != nil {
		return e
	}

	f.f = nil
	f.MMap = nil

	return nil
}

func (f *File) Empty() error {
	err := f.f.Truncate(0)
	if err != nil {
		return err
	}

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&(f.MMap)))
	dh.Len = 0

	return nil
}
