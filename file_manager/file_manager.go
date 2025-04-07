package filemanager

import (
	stdErrors "errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type OpenFile struct {
	file     *os.File
	fmtx     sync.Mutex

	// need more work and refactoring
	cmtx     sync.Mutex
	refCount int
}

/*
	size returns the current size of the file in bytes.
	Returns an error if the file cannot be accessed.
*/
func (of *OpenFile) size() (int64, error) {
	
	of.fmtx.Lock()
	defer of.fmtx.Unlock()
	
	info, err := of.file.Stat()
	
	if err != nil {
		return -1, errors.Wrapf(err, "Open file size")
	}

	return info.Size(), nil
}

func (of *OpenFile) Acquire() {
	of.cmtx.Lock()
	defer of.cmtx.Unlock()
	of.refCount++
}

func (of *OpenFile) Release() {
	of.cmtx.Lock()
	defer of.cmtx.Unlock()
	of.refCount--
}

type FileManager struct {
	openFiles   map[string]*OpenFile
	dbDirectory *os.Root
	blockSize   int
	mtx         sync.Mutex
}

/*
	- create a file manager instence in the dbPath. 
		dbPath need to be a full path (can be absolute or relative)
		blockSize need to be bigger than 0 and prefferably align with the
		OS block size (4096 is most of the time).
TODO - add the removal of temp files created by *materialized* opration.
*/
func NewFileManager(dbPath string, blockSize int) (*FileManager, error) {

	if blockSize <= 0 {
		return nil, errors.New("block size cant be less than 1")
	}

	dbDirectory, err := os.OpenRoot(dbPath)

	if err == nil {
		return &FileManager{
			openFiles:   make(map[string]*OpenFile),
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

/*
any function calling getFile should call the release function of
the requested file to signal not working with it anymore.

this nead more work in the future, currently i dont know some
pieces of the puzzel to write this syncronization logic in a reliable way.
*/
func (fm *FileManager) getFile(fileName string) (*OpenFile, error) {
	if err := validateFileName(fileName); err != nil {
		return nil, errors.Wrapf(err, "getFile failed")
	}

	fm.mtx.Lock()
	defer fm.mtx.Unlock()

	if file, ok := fm.openFiles[fileName]; ok {
		return file, nil
	}

	file, err := fm.dbDirectory.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)

	if err != nil {
		return nil, errors.Wrap(err, "getFile failed")
	}

	openFile := &OpenFile{
		file: file,
	}

	fm.openFiles[fileName] = openFile

	return openFile, nil
}

func (fm *FileManager) validateOffset(fileName string, offset int64) error {
	length, err := fm.Length(fileName)

	if err != nil {
		return errors.Wrapf(err, "File manager validateOffset")
	}

	if length <= offset || offset < 0 {
		return errors.New("File manager validatOffset: Bad Offset")
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

	f.fmtx.Lock()
	defer f.fmtx.Unlock()

	n, err := f.file.ReadAt(p.Bytes(), offset)

	if err != nil && err != io.EOF {
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

	f.fmtx.Lock()
	defer f.fmtx.Unlock()

	n, err := f.file.WriteAt(p.Bytes(), offset)

	if err != nil {
		return 0, errors.Wrapf(err, "File manager Write")
	}
	return n, nil
}

/*
	Length returns the number of blocks in the specified file. 
	A "block" is a fixed-size unit of data defined by the FileManager's 
	block size.
*/
func (fm *FileManager) Length(fileName string) (int64, error) {
	
	f, err := fm.getFile(fileName)

	if err != nil {
		return -1, errors.Wrapf(err, "File manager Length") 
	}
	
	size, err := f.size()

	if err != nil {
		return -1, errors.Wrapf(err, "File manager Length")
	}
	
	return size / int64(fm.blockSize), nil
}

/*
	Extend increases the file's size by one block. 
	Returns an error if the file cannot be accessed or resized.
*/
func (fm *FileManager) Extend(fileName string) error {	
	f, err := fm.getFile(fileName)

	if err != nil {
		return errors.Wrapf(err, "File manager Extend")
	}

	size, err := f.size()

	if err != nil {
		return errors.Wrapf(err, "File manager Extend")
	}

	err = f.file.Truncate(size + int64(fm.blockSize))

	if err != nil {
		return errors.Wrapf(err, "File manager Extend")
	}
	return nil
}

/*
CloseAll closes all open files, following a best-effort approach.
Meaning it attempts to close every file, collects all errors encountered, and returns them as joined error using errors.Join.

After calling CloseAll, the FileManager openFiles map will be empty.

there is a possibility that files will be accessed concurrently so
locking the files before closing is naccessary.

I think open/closed files that are already in use by other
components concurrently and waiting for the CloseAll mtx to unlock
need to be managed in a better way, to discard the files that are
not in use curently safely, this case gonna happen if another
component got the file but close all locked the file for closing
this gonna result in use of file after close and will throw
an error that could be avoided by better managing how the
syncronisation is done between components that use the files.

TODO: i would love to make this function concurrent using go routines

after doing some research on that, here is what i found:

	making this function concurrent is just foolish and absurd, Disk I/O is inherently sequential; concurrent Close() calls on the program level won’t speed up the process or make it concurrent on the disk side. in fact it is less effitiant to do it concurrently because the Goroutine overhead and error aggregation complexity outweigh any theoretical gains.
*/
func (fm *FileManager) CloseAll() error {
	var errs []error

	for path, openF := range fm.openFiles {
		openF.fmtx.Lock()
		err := openF.file.Close()
		openF.fmtx.Unlock()

		if err != nil {
			errs = append(errs, errors.Wrapf(err, "File manager CloseAll"))
		}

		delete(fm.openFiles, path)
	}

	return stdErrors.Join(errs...)
}

func (fm *FileManager) BlockSize() int {
	return fm.blockSize
}

func (fm *FileManager) String() string {
	return fmt.Sprintf("openFiles: %+v\ndbDirectory: %+v\vblockSize: '%d'\nmtx: %v", fm.openFiles, fm.dbDirectory, fm.blockSize, fm.mtx)
}
