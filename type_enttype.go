package objectify

import (
	"io/fs"
	"os"
)

// EntMode is a simplified representation of fs.FileInfo.
type EntMode string

var (
	EntModeDir       EntMode = "dir"
	EntModeLink      EntMode = "link"
	EntModeRegular   EntMode = "regular_file"
	EntModeTemp      EntMode = "temp_file"
	EntModePipe      EntMode = "fifo_pipe"
	EntModeSocket    EntMode = "unix_socket"
	EntModeDevice    EntMode = "device_file"
	EntModeIrregular EntMode = "irregular_file"
	EntModeOther     EntMode = "other"
	EntModeErrored   EntMode = "unknown"
)

// String returns the string representation of the EntMode.
func (e EntMode) String() string {
	return string(e)
}

// getEntMode returns the EntMode and fs.FileInfo for the given path.
// If there is an error in retrieving fs.FileInfo, the function returns
// EntModeErrored and nil.
func getEntMode(path string) (EntMode, fs.FileInfo) {

	info, err := os.Lstat(path)
	if err != nil {
		return EntModeErrored, nil
	}
	return getEntModeWithInfo(info.Mode()), info

}

// getEntModeWithInfo returns the EntMode based on the given fs.FileMode.
// It checks the various flags of fs.FileMode and returns the corresponding
// EntMode value. If none of the flags match, it returns EntModeOther.
// The flags checked are os.ModeDir, os.ModeType, os.ModeSymlink,
// os.ModeTemporary, os.ModeNamedPipe, os.ModeSocket, os.ModeDevice,
// and os.ModeIrregular.
func getEntModeWithInfo(info fs.FileMode) EntMode {

	if info&os.ModeDir != 0 {
		return EntModeDir
	}
	if info&os.ModeType == 0 {
		return EntModeRegular
	}
	if info&os.ModeSymlink != 0 {
		return EntModeLink
	}
	if info&os.ModeTemporary != 0 {
		return EntModeTemp
	}
	if info&os.ModeNamedPipe != 0 {
		return EntModePipe
	}
	if info&os.ModeSocket != 0 {
		return EntModeSocket
	}
	if info&os.ModeDevice != 0 {
		return EntModeDevice
	}
	if info&os.ModeIrregular != 0 {
		return EntModeIrregular
	}

	return EntModeOther

}
