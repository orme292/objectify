package objectify

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

type Files []*FileObj

// FileObj represents a directory entry object.
type FileObj struct {

	// UpdatedAt represents the last time this struct was updated.
	UpdatedAt time.Time

	// modTime represents the last time the directory entry was modified.
	modTime time.Time

	// Filename is the base name of the directory entry.
	// Root is the parent directory.
	Filename string
	Root     string

	// SizeBytes is the size of the file in Bytes.
	SizeBytes int64

	// ChecksumMD5 and ChecksumSHA256 are hash byte array string-representations.
	// MD5 and SHA256 are the hash byte arrays.
	ChecksumMD5    string
	MD5            []byte
	ChecksumSHA256 string
	SHA256         []byte

	// Mode is the EntMode of the directory entry.
	// modeFS is returned from os.Lstat
	Mode EntMode
	info fs.FileInfo

	// Target will be populated with a symlinks target path.
	Target      string
	TargetFinal string

	IsLink     bool
	IsReadable bool
	IsExists   bool

	Set *Sets
}

type Action int

const (
	F_CHECKSUM_MD5 Action = iota
	F_CHECKSUM_SHA256
	F_MODES
	F_SIZE
	F_LINKTARGET
)

// newFileObj creates a new instance of FileObj based on the provided
// path and Sets. If the path is empty, it returns nil. Otherwise,
// it splits the path into directory and file, and initializes the FileObj
// with the extracted values. The Sets field of the FileObj is set to the
// provided Sets. If the file exists and is readable, it calls the Update
// method to populate additional information. Finally, it sets the timestamp
// of the FileObj.
func newFileObj(path string, s Sets) *FileObj {

	if path == EMPTY {
		return nil
	}

	dir, file := pathBaseSplit(path)

	fo := &FileObj{
		Filename: file,
		Root:     dir,
		Set:      &s,
	}

	_ = fo.update()

	return fo

}

// hasPaths checks if the FileObj has valid values for Filename
// and Root fields. Returns true if both fields are not empty,
// otherwise returns false.
func (fo *FileObj) hasPaths() bool {

	if fo.Filename == EMPTY || fo.Root == EMPTY {
		return false
	}
	return true

}

// setChecksums calculates and sets the checksums (SHA256 and MD5) of the file specified by
// the FileObj's FullPath.
// If Sets.ChecksumSHA256 is true, it calculates and sets the SHA256 checksum.
// If Sets.ChecksumMD5 is true, it calculates and sets the MD5 checksum.
// Returns an error if there is any failure in calculating the checksums.
func (fo *FileObj) setChecksums() error {

	var err error

	if fo.IsExists && fo.IsReadable {

		if fo.Set.ChecksumSHA256 {
			fo.SHA256, fo.ChecksumSHA256, err = getSHA256(fo.FullPath())
			if err != nil {
				return err
			}
		}
		if fo.Set.ChecksumMD5 {
			fo.MD5, fo.ChecksumMD5, err = getMD5(fo.FullPath())
			if err != nil {
				return err
			}

		}
	}

	return nil

}

// setEntMode updates the Mode, info, modTime, and IsLink fields of the FileObj
// based on the values of IsExists, IsReadable, and Sets.Modes.
// If IsExists is true and IsReadable is true, it sets the Mode field by calling getEntMode
// and assigns the returned value to both the Mode and info fields.
// It also updates the modTime field by retrieving the modification time from info.
// If Sets.Modes is true and Mode is EntModeLink, it sets the IsLink field to true.
// Returns the value of the Mode field.
func (fo *FileObj) setEntMode() EntMode {

	if fo.IsExists && fo.IsReadable {

		if fo.Set.Modes {
			fo.Mode = getEntModeWithInfo(fo.info.Mode())
			fo.modTime = fo.info.ModTime()
		}

		if fo.Set.Modes && (fo.Mode == EntModeLink) {
			fo.IsLink = true
		}

	}

	return fo.Mode

}

// setPrelims updates preliminary information about the FileObj instance.
// It sets the info field with the return value of attemptStat method, and
// updates the IsExists and IsReadable fields based on the presence and
// readability of the file.
// Returns true if the FileObj has valid paths, the file exists and is readable,
// otherwise returns false.
func (fo *FileObj) setPrelims() bool {

	var ok bool

	if fo.hasPaths() {

		fo.info, ok = attemptStat(fo.FullPath())
		if !ok {
			fo.IsExists = false
			fo.IsReadable = false
		}

		if isReadable(fo.FullPath()) {
			fo.IsExists = true
			fo.IsReadable = true
		}

	} else {

		fo.IsExists = false
		fo.IsReadable = false

	}

	return fo.IsExists && fo.IsReadable

}

// setReadable sets the IsReadable field of the FileObj by calling the
// isReadable function with the FileObj's FullPath. The result is assigned
// to the IsReadable field and returned.
func (fo *FileObj) setReadable() bool {

	fo.IsReadable = isReadable(fo.FullPath())

	return fo.IsReadable

}

// setSize sets the size of the FileObj if Sets.Size is true.
// If fo.modeFS is nil, attemptStat is called to retrieve fs.FileInfo.
// If info is still nil, fo.SizeBytes is set to 0 and the function returns.
// Otherwise, fo.SizeBytes is set to the size provided by info.
func (fo *FileObj) setSize() {

	if fo.Set.Size {

		if fo.info == nil {
			fo.info, _ = attemptStat(fo.FullPath())
		}

		if fo.info == nil {
			fo.SizeBytes = 0
			return
		}

		fo.SizeBytes = fo.info.Size()

	}

}

func (fo *FileObj) setTargets() {

	if fo.IsExists && fo.IsReadable && fo.IsLink {

		if fo.Set.LinkTarget || fo.Set.LinkTargetFinal {

			if !fo.Set.Modes {
				fo.Mode, fo.info = getEntMode(fo.FullPath())
			}

		}

		if fo.Set.LinkTarget {
			fo.Target, _ = getsTarget(fo.FullPath())
		}

		if fo.Set.LinkTargetFinal {
			fo.TargetFinal, _ = getsFinalTarget(fo.FullPath(), fo.info)
		}

	}

}

// timestamp sets the UpdatedAt field of the FileObj to the current
// time and returns it.
func (fo *FileObj) timestamp() time.Time {

	fo.UpdatedAt = time.Now()
	return fo.UpdatedAt

}

// update updates the FileObj by performing the following actions:
//   - If setPrelims (which sets info, checks exists and readability) passes, then:
//   - Calls setEntMode to update the Mode, modTime, and IsLink fields
//   - Calls setSize to update the SizeBytes field based on the file size
//   - Calls setTargets to update the Target and/or TargetFinal fields if
//     Sets.LinkTarget/Sets.LinkTargetFinal is true
//   - Calls setChecksums to update the checksums (SHA256 and MD5) if file exists and
//     is readable
//   - Calls timestamp to update the UpdatedAt field to the current time
//
// Returns nil (for now).
func (fo *FileObj) update() error {

	if fo.setPrelims() {

		_ = fo.setEntMode()
		fo.setSize()
		fo.setTargets()
		_ = fo.setChecksums()
		fo.timestamp()

	}

	return nil

}

// ChangeSets overwrites the Set field with a new Sets object.
func (fo *FileObj) ChangeSets(s Sets) {

	fo.Set = &s

}

// Force applies the specified action to the FileObj by changing its sets and
// calling the corresponding helper methods. The original sets of the FileObj are
// stored temporarily and restored after applying the action.
// The available actions are:
//   - F_CHECKSUM_MD5: Changes the sets to enable checksum MD5 calculation and
//     calls the setChecksums() method.
//   - F_CHECKSUM_SHA256: Changes the sets to enable checksum SHA256 calculation
//     and calls the setChecksums() method.
//   - F_MODES: Changes the sets to enable mode and file info retrieval and calls
//     the setEntMode() method.
//   - F_SIZE: Changes the sets to enable size calculation and calls the setSize()
//     method.
//   - F_LINKTARGET: Changes the sets to enable link target retrieval and calls
//     the setTargets() method.
func (fo *FileObj) Force(a Action) {

	originalSets := fo.Set

	switch a {
	case F_CHECKSUM_MD5:

		fo.ChangeSets(Sets{ChecksumMD5: true})
		_ = fo.setChecksums()

	case F_CHECKSUM_SHA256:

		fo.ChangeSets(Sets{ChecksumSHA256: true})
		_ = fo.setChecksums()

	case F_MODES:

		fo.ChangeSets(Sets{Modes: true})
		_ = fo.setEntMode()

	case F_SIZE:

		fo.ChangeSets(Sets{Size: true})
		fo.setSize()

	case F_LINKTARGET:

		fo.ChangeSets(Sets{LinkTarget: true})
		fo.setTargets()

	}

	fo.Set = originalSets

}

// FullPath returns the full path of the FileObj by joining the Root and Filename.
// Utilizes filepath.Join to combine the two components.
func (fo *FileObj) FullPath() string {
	return filepath.Join(fo.Root, fo.Filename)
}

// HasChanged checks if the file specified by FileObj has been modified since
// its last update. It returns true if the file exists, is readable, and its
// modification time is after the last update time. Otherwise, it returns false.
func (fo *FileObj) HasChanged() bool {

	if fo.IsExists && fo.IsReadable {

		info, ok := attemptStat(fo.FullPath())
		if !ok {
			return false
		}

		return info.ModTime().After(fo.modTime)

	}

	return false

}

// SecondsSinceUpdatedAt returns the number of seconds since the UpdatedAt time of
// the FileObj.
func (fo *FileObj) SecondsSinceUpdatedAt() int64 {
	return int64(time.Now().Sub(fo.UpdatedAt).Seconds())
}

// SizeString returns the formatted string representation of the size in bytes.
// Converts the given size in bytes to a human-readable format (e.g., KB, MB, GB, etc.).
func (fo *FileObj) SizeString() string {
	return sizeString(fo.SizeBytes)
}

// Update checks if the file specified by FileObj has been
// modified since its last update. If it has changed, and
// the file exists, is readable, and its modification time
// is after the last update time, Update calls update.
func (fo *FileObj) Update() *FileObj {

	if fo.HasChanged() {

		_ = fo.update()

	}

	return fo

}

/* DEBUG */

func (fo *FileObj) DebugOut() {
	fmt.Println("=========")
	fmt.Printf("Filename: %s\nRoot: %s\n", fo.Filename, fo.Root)
	fmt.Printf("Size: %s\n", fo.SizeString())
	fmt.Printf("ChecksumMD5: %s\nChecksumSHA256: %s\n", fo.ChecksumMD5, fo.ChecksumSHA256)
	fmt.Printf("EntMode: %s\n", fo.Mode.String())
	fmt.Printf("Target: %s\n", fo.Target)
	fmt.Printf("IsExists: %t\nIsReadable: %t\nIsLink: %t\n", fo.IsExists, fo.IsReadable, fo.IsLink)
	fmt.Printf("Sets: %v\n", fo.Set)
	fmt.Printf("modTime: %s\n", fo.modTime.Format("Mon Jan 2 15:04:05 MST 2006"))
}
