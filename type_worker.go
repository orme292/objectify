package objectify

import (
	"os"
)

// worker represents a worker that performs operations on files and directories.
type worker struct {
	RootPath       string
	singleFileMode bool
	setter         Sets
}

// newPathWorker creates a new instance of the worker struct with the provided startPath and Sets.
// It returns a pointer to the worker.
func newPathWorker(startPath string, s Sets) *worker {
	return &worker{
		RootPath:       startPath,
		singleFileMode: false,
		setter:         s,
	}
}

func newFileWorker(path string, s Sets) *worker {
	return &worker{
		RootPath:       path,
		singleFileMode: true,
		setter:         s,
	}
}

// validate checks if the worker's RootPath is non-empty.
// If it is empty, it returns false.
// If the worker is in "single" file mode, it checks if the RootPath
// points to a file. If so, it returns true.
// Otherwise, it returns true.
func (w *worker) validate() bool {

	if w.RootPath == EMPTY {
		return false
	}

	if w.singleFileMode {
		return isFile(w.RootPath)
	}

	return true

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
