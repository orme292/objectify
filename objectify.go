// Package objectify reads a provided directory and turns each directory entry / file into an
// object. Each object (FileObj) contains calculated checksums, symlink target paths, size and
// file mode information, and provides functions to check for recent modifications.
package objectify

import (
	"fmt"
	"os"
	"path/filepath"
)

// Path is a function that takes a rootPath and a Sets struct as parameters.
// It creates a worker instance and runs it using the run function to collect file information.
// It returns a slice of FileObj structs and an error if any.
// The Sets struct is used to specify which fields of the FileObj struct need to be populated.
func Path(rootPath string, s Sets) (files Files, err error) {

	return run(newPathWorker(rootPath, s))

}

// File is a function that accepts a path and a Sets struct as parameters.
// It creates a worker instance and runs it using the run function to collect file information.
// It returns a slice of FileObj structs and an error if any. The returned slice should contain
// a single FileObj, if not, nil and an error is returned.
// As long as the files slice contains a single FileObj, it is returned.
func File(path string, s Sets) (file *FileObj, err error) {

	files, err := run(newFileWorker(path, s))
	if err != nil || len(files) == 0 || len(files) > 1 {
		return nil, err
	}

	return files[0], nil

}

// run is a function that takes a worker pointer w as a parameter. It first validates
// the worker by calling its validate method. If the validation fails, it returns
// an error indicating that the StartingPath is inaccessible. If the worker has no
// non-directory entries, it returns an error indicating that the StartingPath has
// no non-directory entries. It then initializes an empty slice of FileObj structs.
// It reads the directory entries using os.ReadDir and iterates over each entry.
// If the entry is a directory, it continues to the next one. If the entry is a symlink
// and it leads to a directory, it continues to the next one. Otherwise, it creates
// a new FileObj using the newFileObj function and appends it to the files slice.
// Finally, it returns the files slice and any error that occurred during the process.
func run(w *worker) (Files, error) {

	// validate checks if there is a valid path provided.
	if !w.validate() {
		return nil, fmt.Errorf("StartingPath is not correct: %s", w.RootPath)
	}

	// checks to see that the provided path contains actual file entries.
	// may be removed in the future.
	if !w.singleFileMode {
		if !w.hasEntries() {
			return nil, fmt.Errorf("StartingPath has no non-directory entries: %s", w.RootPath)
		}
	}

	files := Files{}

	if w.singleFileMode {

		file := newFileObj(w.RootPath, w.setter)
		files = append(files, file)

		return files, nil

	}

	dirents, err := os.ReadDir(w.RootPath)
	if err != nil {
		return nil, err
	}

	for _, ent := range dirents {

		path := filepath.Join(w.RootPath, ent.Name())

		if ent.IsDir() {
			continue
		}
		if ent.Type()&os.ModeSymlink != 0 {
			if linkLeadsToDir(path) {
				continue
			}
		}

		file := newFileObj(path, w.setter)
		files = append(files, file)

	}

	return files, err

}
