package filemanager_test

import (
	filemanager "Clotho/file_manager"
	"Clotho/file_manager/block"
	"Clotho/file_manager/page"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestNewFileManager(t *testing.T) {
	tests := []struct {
		name        string
		dbPath      string
		blockSize   int
		expectError bool
		setup       func(testPath string) error
	}{
		{
			name:        "dbPath exist",
			dbPath:      "exist",
			blockSize:   512,
			expectError: false,
			setup: func(testPath string) error {
				return os.Mkdir(testPath, 0755)
			},
		},
		{
			name:        "dbPath doesnt exist",
			dbPath:      "doesntExist",
			blockSize:   512,
			expectError: false,
			setup:       nil,
		},
		{
			name:        "dbPath exist but is a file",
			dbPath:      "existbutfile",
			blockSize:   512,
			expectError: true,
			setup: func(testPath string) error {
				_, err := os.Create(testPath)

				return err
			},
		},
		{
			name:        "dbPath exist but bad permission",
			dbPath:      "exitbutbadperm",
			blockSize:   512,
			expectError: true,
			setup: func(testPath string) error {
				return os.Mkdir(testPath, 0111)
			},
		},
		{
			name:        "block size negative",
			dbPath:      "doesntExist",
			blockSize:   -512,
			expectError: true,
			setup:       nil,
		},
	}

	tempDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := filepath.Join(tempDir, tt.dbPath)

			if tt.setup != nil {
				if err := tt.setup(testPath); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			fm, err := filemanager.NewFileManager(testPath, tt.blockSize)

			if (err != nil) != tt.expectError {
				t.Errorf("expectError %v but got: %v", tt.expectError, err)
			}

			if fm != nil && !tt.expectError {
				if _, err := os.Stat(testPath); err != nil {
					t.Errorf("db path wasnt created")
				}

				if fm.BlockSize() != tt.blockSize {
					t.Errorf("wrong block size got %v but expected %v", fm.BlockSize(), tt.blockSize)
				}
			}
		})
	}
}

// TODO: i need to write tests, but i dont see a point of doing it now

func TestReadWrite(t *testing.T) {
	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, "testDB")

	fm, err := filemanager.NewFileManager(dbPath, 32)

	if err != nil {
		t.Fatalf("failed to create the file manager: %s", err.Error())
	}

	data := "wsuup"
	
	dataBlock := block.BlockID {
		Filename: "testFile",
		Blknum: 0,
	}

	fm.Extend(dataBlock.Filename)
	expectedPage := page.NewPage(uint(fm.BlockSize()))
	expectedPage.SetString(data, 10)


	_, err = fm.Write(dataBlock, expectedPage)

	if err != nil {
		t.Fatalf("failed to write to the file manager test file: %s", err.Error())
	}

	gottenPage := page.NewPage(uint(fm.BlockSize()))

	_, err = fm.Read(dataBlock, gottenPage)

	if err != nil {
		t.Fatalf("failed to read the data into memory: %s", err.Error())
	}

	if !slices.Equal(gottenPage.Bytes(), expectedPage.Bytes()) {
		t.Errorf("got %v but expected %v", gottenPage.Bytes(), expectedPage.Bytes())
	}
}
