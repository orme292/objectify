package objectify

import (
	"os"
)

// worker represents a worker that performs operations on files and directories.
type worker struct {
	RootPath string
	setter   Sets
}

// newWorker creates a new instance of the worker struct with the provided startPath and Sets.
// It returns a pointer to the worker.
func newWorker(startPath string, s Sets) *worker {
	return &worker{
		RootPath: startPath,
		setter:   s,
	}
}

// validate checks if the RootPath field of the worker struct is empty.
// If it is empty, it returns false. Otherwise, it calls the isReadable
// function passing the absolute path of the RootPath as an argument.
// If the isReadable function returns true, indicating that the file is
// readable, validate returns true. Otherwise, it returns false.
func (w *worker) validate() bool {

	if w.RootPath == EMPTY {
		return false
	}

	return isReadable(pathAbsUnsafe(w.RootPath))

}

// hasEntries checks if the worker's RootPath directory has any non-directory entries.
// If reading the directory fails (due to an error), it returns false.
// If the directory is empty (no entries), it returns false.
// It iterates over each directory entry and if any of them is not a directory,
// it returns true.
// If all directory entries are directories, it returns false.
func (w *worker) hasEntries() bool {

	dirents, err := os.ReadDir(w.RootPath)
	if err != nil {
		return false
	}

	if len(dirents) == 0 {
		return false
	}

	for _, ent := range dirents {
		if !ent.IsDir() {
			return true
		}
	}

	return false

}
