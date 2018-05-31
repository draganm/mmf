package mmf

import "os"

type File struct {
	f *os.File
}

func Open(fileName string) (*File, error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_SYNC|os.O_RDWR, 0700)
	if err != nil {
		return nil, err
	}

	return &File{
		f: f,
	}, nil
}
