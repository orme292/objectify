package objectify

import (
	"io/fs"
	"os"
)

// entMode is a simplified representation of fs.FileInfo.
type entMode string

var (
	entModeDir       entMode = "dir"
	entModeLink      entMode = "link"
	entModeRegular   entMode = "regular_file"
	entModeTemp      entMode = "temp_file"
	entModePipe      entMode = "fifo_pipe"
	entModeSocket    entMode = "unix_socket"
	entModeDevice    entMode = "device_file"
	entModeIrregular entMode = "irregular_file"
	entModeOther     entMode = "other"
	entModeErrored   entMode = "unknown"
)

// String returns the string representation of the entMode.
func (e entMode) String() string {
	return string(e)
}

// getEntMode returns the entMode and fs.FileInfo for the given path.
// If there is an error in retrieving fs.FileInfo, the function returns
// entModeErrored and nil.
func getEntMode(path string) (entMode, fs.FileInfo) {

	info, err := os.Lstat(path)
	if err != nil {
		return entModeErrored, nil
	}
	return getEntModeWithInfo(info.Mode()), info

}

// getEntModeWithInfo returns the entMode based on the given fs.FileMode.
// It checks the various flags of fs.FileMode and returns the corresponding
// entMode value. If none of the flags match, it returns entModeOther.
// The flags checked are os.ModeDir, os.ModeType, os.ModeSymlink,
// os.ModeTemporary, os.ModeNamedPipe, os.ModeSocket, os.ModeDevice,
// and os.ModeIrregular.
func getEntModeWithInfo(info fs.FileMode) entMode {

	if info&os.ModeDir != 0 {
		return entModeDir
	}
	if info&os.ModeType == 0 {
		return entModeRegular
	}
	if info&os.ModeSymlink != 0 {
		return entModeLink
	}
	if info&os.ModeTemporary != 0 {
		return entModeTemp
	}
	if info&os.ModeNamedPipe != 0 {
		return entModePipe
	}
	if info&os.ModeSocket != 0 {
		return entModeSocket
	}
	if info&os.ModeDevice != 0 {
		return entModeDevice
	}
	if info&os.ModeIrregular != 0 {
		return entModeIrregular
	}

	return entModeOther

}
