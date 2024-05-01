package objectify

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	EMPTY = ""
)

// attemptOpen opens a file at the specified path and returns true
// if successful, false otherwise.
func attemptOpen(path string) bool {

	var opens = true

	f, err := os.Open(path)
	defer func(f *os.File) {
		cErr := f.Close()
		if cErr != nil {
			opens = false
		}
	}(f)
	if err != nil {
		opens = false
	}

	return opens

}

// attemptStat returns the fs.FileInfo of the file at the specified path
// using os.Lstat. If the operation is successful, it returns the FileInfo and true.
// Otherwise, it returns nil and false.
func attemptStat(path string) (fs.FileInfo, bool) {

	info, err := os.Lstat(path)
	if err != nil || info == nil {
		return nil, false
	}

	return info, true

}

// calcSHA256 calculates the SHA256 hash of the content of the provided file.
// It returns nil if the file is nil or if an error occurs during the hashing process.
// Otherwise, it returns the SHA256 hash as a byte array.
func calcSHA256(f *os.File) []byte {

	if f == nil {
		return nil
	}

	hash := sha256.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		return nil
	}
	return hash.Sum(nil)

}

// calcMD5 calculates the MD5 hash of the content of the provided file.
// It returns nil if the file is nil or if an error occurs during the hashing process.
// Otherwise, it returns the MD5 hash as a byte array.
func calcMD5(f *os.File) []byte {

	if f == nil {
		return nil
	}

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil
	}
	return hash.Sum(nil)

}

// getSHA256 opens the file at the specified path and calculates
// the SHA256 hash of its content. It returns the SHA256 hash as a
// byte array, the hash as a hexadecimal string, and any error that occurs.
// If the file cannot be opened, it returns nil for the hash and an error.
// If there is an error during the hashing process, it returns nil for
// the hash and the error.
func getSHA256(path string) ([]byte, string, error) {

	f, err := os.Open(path)
	defer func(f *os.File) {
		cErr := f.Close()
		if cErr != nil {
			err = cErr
		}
	}(f)
	if err != nil {
		return nil, EMPTY, err
	}

	sum := calcSHA256(f)

	return sum, fmt.Sprintf("%x", sum), nil

}

// getMD5 opens the file at the specified path and calculates
// the MD5 hash of its content. It returns the MD5 hash as a
// byte array, the hash as a hexadecimal string, and any error that occurs.
// If the file cannot be opened, it returns nil for the hash and an error.
// If there is an error during the hashing process, it returns nil for
// the hash and the error.
func getMD5(path string) ([]byte, string, error) {

	f, err := os.Open(path)
	defer func(f *os.File) {
		cErr := f.Close()
		if cErr != nil {
			err = cErr
		}
	}(f)
	if err != nil {
		return nil, EMPTY, err
	}

	sum := calcMD5(f)

	return sum, fmt.Sprintf("%x", sum), nil

}

// getsTarget returns the target of a symbolic link at the specified path
// and a bool indicating if the retrieval was successful.
func getsTarget(path string) (string, bool) {

	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return target, false
	}

	return target, true

}

// getsFinalTarget returns the final target of a symbolic link and a bool
// indicating if the operation is successful. It takes the path of the symlink
// and the fs.FileInfo of the symlink itself. If the fs.FileInfo is nil, it will
// return an empty string and false. If the symlink is a directory, it will return
// an empty string and false. If the symlink points to another symlink, it will
// recursively evaluate the target until it reaches the final target. If any error
// occurs during evaluation, it will return an empty string and false.
func getsFinalTarget(path string, info fs.FileInfo) (string, bool) {

	if info == nil {
		return EMPTY, false
	}

	if info.Mode()&os.ModeSymlink != 0 {

		target, err := filepath.EvalSymlinks(path)
		if err != nil {
			return EMPTY, false
		}

		info, _ := attemptStat(target)

		return getsFinalTarget(target, info)

	}

	if info.IsDir() {
		return EMPTY, false
	}

	return path, true

}

// isFile checks if the specified path corresponds to a file. It uses the
// attemptStat function to get the fs.FileInfo of the path, and then returns
// true if the info is not nil and represents a non-directory file. Otherwise,
// it returns false.
func isFile(path string) bool {

	info, _ := attemptStat(path)
	return info != nil && !info.IsDir()

}

// isReadable checks if a file at the specified path is readable by attempting
// to obtain its file information using the attemptStat function. If the file
// information is successfully obtained, it then attempts to open the file using
// the attemptOpen function. If both operations are successful, it returns true,
// indicating that the file is readable. Otherwise, it returns false.
func isReadable(path string) bool {

	return attemptOpen(path)

}

// linkLeadsToDir checks if the specified symbolic link leads to a directory. It first attempts to
// retrieve the FileInfo using the attemptStat function. If the FileInfo is not found,
// it returns false. If the FileInfo represents a directory, it returns true. If the FileInfo
// represents a symbolic link, it uses filepath.EvalSymlinks to evaluate the target path.
// If an error occurs during the evaluation, it returns false. Otherwise, it recursively
// calls linkLeadsToDir on the target path. If none of the conditions are met, it returns false.
func linkLeadsToDir(path string) bool {

	info, ok := attemptStat(path)
	if !ok {
		return false
	}

	if info.IsDir() {
		return true
	}

	if info.Mode()&os.ModeSymlink != 0 {

		target, err := filepath.EvalSymlinks(path)
		if err != nil {
			return false
		}

		return linkLeadsToDir(target)

	}

	return false

}

// pathBaseSplit extracts the directory and file components from the specified path.
// If the path is empty, it returns empty strings for both directory and file.
func pathBaseSplit(path string) (dir, file string) {

	if path == EMPTY {
		return EMPTY, EMPTY
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = filepath.Join("/", filepath.Base(path))
	}

	return filepath.Dir(abs), filepath.Base(abs)

}

// pathAbsSafe returns the absolute path of the file at the specified
// path using filepath.Abs. If an error occurs during the operation, it
// returns the path joined with the root directory ("/").
func pathAbsSafe(path string) string {

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = filepath.Join("/", filepath.Base(path))
	}

	return abs

}

// sizeString returns the formatted string representation of the size in bytes.
// It converts the given size in bytes to a human-readable format (e.g., KB, MB, GB, etc.).
func sizeString(bytes int64) string {

	var unit = int64(1024)

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])

}
