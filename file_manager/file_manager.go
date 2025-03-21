package filemanager

import (
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type File struct {
	file *os.File
	mtx sync.Mutex
}

type FileManager struct {
	openFiles   map[string]*File
	dbDirectory *os.Root
	blockSize   int
	mtx         sync.Mutex
}

/*
	it need to check if the dbPath exist, if not create one.
	TODO - add the removal of temp files created by *materialized* opration.
*/
func NewFileManager(dbPath string, blockSize int) (*FileManager, error) {

	if blockSize <= 0 {
		return nil, errors.New("block size cant be less than 1")
	}

	dbDirectory, err := os.OpenRoot(dbPath)

	if err == nil {
		return &FileManager {
			openFiles:   make(map[string]*File),
			dbDirectory: dbDirectory,
			blockSize:   blockSize,
		}, nil
	}

	err = os.Mkdir(dbPath, 0755)

	if err != nil {
		return nil, errors.Wrapf(err, "File Manager creation failed (probably is a file or bad permissions), path: %q", dbPath)
	}

	return NewFileManager(dbPath, blockSize)
}


func (fm *FileManager) Read(blk BlockID, p *Page) {

}



func (fm *FileManager) BlockSize() int {
	return fm.blockSize
}

func (fm *FileManager) String() string {
	return fmt.Sprintf("openFiles: %+v\ndbDirectory: %+v\vblockSize: '%d'\nmtx: %v", fm.openFiles, fm.dbDirectory, fm.blockSize, fm.mtx)
}
