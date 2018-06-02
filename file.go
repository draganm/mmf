package mmf

import (
	"os"

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

	mmap := gommap.MMap{}

	if stat.Size() != 0 {
		mmap, err = gommap.Map(f.Fd(), gommap.PROT_READ, gommap.MAP_SHARED)
		if err != nil {
			return nil, err
		}
		err = mmap.Advise(gommap.MADV_RANDOM)
		if err != nil {
			return nil, err
		}
	}

	return &File{
		f:    f,
		MMap: mmap,
	}, nil
}
