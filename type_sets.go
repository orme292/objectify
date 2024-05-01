package objectify

// Sets fields are flags for FileObj fields which can be optionally populated.
type Sets struct {
	Size            bool
	Modes           bool
	ChecksumMD5     bool
	ChecksumSHA256  bool
	LinkTarget      bool
	LinkTargetFinal bool
}

// SetsAll returns a Sets object with all fields set to true.
func SetsAll() Sets {
	return Sets{
		Size:            true,
		Modes:           true,
		ChecksumMD5:     true,
		ChecksumSHA256:  true,
		LinkTarget:      true,
		LinkTargetFinal: true,
	}
}

// SetsAllNoChecksums sets all fields of the Sets object to true,
// except for the ChecksumMD5 and ChecksumSHA256.
func SetsAllNoChecksums() Sets {
	s := SetsAll()
	s.ChecksumMD5 = false
	s.ChecksumSHA256 = false
	return s
}

// SetsAllMD5 returns a Sets object with all fields set to true,
// except for ChecksumSHA256.
func SetsAllMD5() Sets {
	s := SetsAll()
	s.ChecksumSHA256 = false
	return s
}

// SetsAllSHA256 returns a Sets object with all fields set to true,
// except for ChecksumMD5.
func SetsAllSHA256() Sets {
	s := SetsAll()
	s.ChecksumMD5 = false
	return s
}

// SetsNone returns a Sets object with all fields set to false.
func SetsNone() Sets {
	return Sets{}
}
