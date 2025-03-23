package filemanager

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type File struct {
	file *os.File
	mtx  sync.Mutex
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
		return &FileManager{
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

func validateFileName(fileName string) error {
	switch fileName {
	case "":
		return errors.New("fileName cant be empty")
	default:
		return nil
	}
}

func (fm *FileManager) getFile(fileName string) (*File, error) {

	if err := validateFileName(fileName); err != nil {
		return nil, errors.Wrapf(err, "getFile failed")
	}

	if file, ok := fm.openFiles[fileName]; ok {
		return file, nil
	}

	file, err := fm.dbDirectory.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)

	if err != nil {
		return nil, errors.Wrap(err, "getFile failed")
	}

	fm.openFiles[fileName] = &File{
		file: file,
	}

	return fm.getFile(fileName)
}

func (fm *FileManager) validateOffset(fileName string, offset int64) error {
	length, err := fm.Length(fileName)

	if err != nil {
		return errors.Wrapf(err, "File manager validateOffset")
	}

	if length <= offset || offset < 0 {
		return errors.Wrapf(err, "File manager validatOffset: Bad Offset")
	}

	return nil
}

// test reading from bad offsets ( offset more than file size, less than 0)
/*
	TODO: add documentation to this function comment
*/
func (fm *FileManager) Read(blk BlockID, p *Page) (int, error) {
	f, err := fm.getFile(blk.Filename)

	if err != nil {
		return 0, errors.Wrapf(err, "File manager Read")
	}

	var offset int64 = int64(blk.Blknum) * int64(fm.blockSize)

	err = fm.validateOffset(blk.Filename, int64(blk.Blknum))

	if err != nil {
		return 0, errors.Wrapf(err, "File manager Read")
	}

	f.mtx.Lock()
	defer f.mtx.Unlock()

	n, err := f.file.ReadAt(p.Bytes(), offset)

	if err != nil && err != io.EOF{
		return 0, errors.Wrapf(err, "File manager Read")
	}
	return n, nil
}

/*
TODO: add documentation to this function comment
*/
func (fm *FileManager) Write(blk BlockID, p *Page) (int, error) {
	f, err := fm.getFile(blk.Filename)

	if err != nil {
		return 0, errors.Wrapf(err, "File manager Write")
	}

	var offset int64 = int64(blk.Blknum) * int64(fm.blockSize)

	err = fm.validateOffset(blk.Filename, offset)

	if err != nil {
		return 0, errors.Wrapf(err, "File manager Write")
	}

	f.mtx.Lock()
	defer f.mtx.Unlock()

	n, err := f.file.WriteAt(p.Bytes(), offset)

	if err != nil {
		return 0, errors.Wrapf(err, "File manager Write")
	}
	return n, nil
}

/*
TODO: add documentation to this function comment
*/
func (fm *FileManager) Length(fileName string) (int64, error) {
	f, err := fm.getFile(fileName)

	if err != nil {
		return -1, errors.Wrapf(err, "File manager Length")
	}

	f.mtx.Lock()
	defer f.mtx.Unlock()

	info, err := f.file.Stat()

	if err != nil {
		return -1, errors.Wrapf(err, "File manager Length")
	}

	return info.Size() / int64(fm.blockSize), nil
}

func (fm *FileManager) BlockSize() int {
	return fm.blockSize
}

func (fm *FileManager) String() string {
	return fmt.Sprintf("openFiles: %+v\ndbDirectory: %+v\vblockSize: '%d'\nmtx: %v", fm.openFiles, fm.dbDirectory, fm.blockSize, fm.mtx)
}
