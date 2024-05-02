# Objectify

[![Go Reference](https://pkg.go.dev/badge/github.com/orme292/objectify.svg)](https://pkg.go.dev/github.com/orme292/objectify@v0.3.1)

Objectify is a Go package that reads a directory's entries and returns a slice of structs which contain information
about each directory entry like size, file mode, the symlink target, and checksums.

## Import this Module

```shell
go get github.com/orme292/objectify@v0.2.0
```

```go
import (
    objf "github.com/orme292/objectify"
)
```

## Usage

Objectify can be called by passing a path and a Sets struct.

### Sets

The Sets struct tells Objectify which fields should be populated for each directory entry.

```go
func main() {

setter := objf.Sets{
        Size: true,
        Modes: true,
        ChecksumMD5: true,
        ChecksumSHA256: true,
        LinkTarget: true,
        LinkTargetFinal: true,
    }

}
```

You can also have a Sets object returned by using a builder function:
- `setter := SetsAll()` All fields will be populated.
- `setter := SetsAllNoChecksums()` All fields except ChecksumSHA256/ChecksumMD5 will be populated.
- `setter := SetsAllMD5()` All fields except ChecksumSHA256 will be populated.
- `setter := SetsAllSHA256()` All fields except ChecksumMD5 will be populated
- `setter := SetsNone()` No optional fields will be populated.

### Call Objectify

You can call objectify by using the `Path()` or `File()` functions. 

`Path()` will walk a directory and return a `Files` slice and an error, if there is one:
```go
files, err := objf.Path("/root/path", objf.SetsAll())
```

`File()` will process a single file and return a `FileObj` struct and an error, if there is one:
```go
file, err := objf.File("/root/path/myfile.txt", objf.SetsNone())
```

Create your own `Sets` for more configuration:
```go
setter := objf.Sets{
    Size: true,
    Modes: true,
    LinkTarget: true,
    LinkTargetFinal: true,
}
files, err := objf.Path("/root/path", setter)
```

### The *Files* & *FileObj* Types

`Path()` returns a `Files` slice. The `Files` slice is made of `FileObj` structs.
`File()` returns a `FileObj` struct.

```go
type Files []*FileObj
```
```go
type FileObj struct {
    UpdatedAt time.Time

    Filename string
    Root     string

    SizeBytes int64

    ChecksumMD5    string
    MD5            []byte
    ChecksumSHA256 string
    SHA256         []byte

    Mode   entMode
    info   fs.FileMode

    Target      string
    TargetFinal string

    IsLink     bool
    IsReadable bool
    IsExists   bool

    Sets *Sets
}
```

## `FileObj` methods

- `FileObj.ChangeSets()` updates the Sets, but does not trigger an update.
- `FileObj.Force()` Forces an update on an optional field, despite Sets values.
- `FileObj.FullPath()` returns a string that joins the root directory with the entry's filename.
- `FileObj.HasChanged()` returns `true` if the file has changed since the struct was last populated.
- `FileObj.SecondsSinceUpdatedAt()` returns the number of seconds elapsed since the FileObj's fields were updated.
- `FileObj.SizeString()` returns a human-readable string representation of the directory entry's size (i.e. 500 MB)
- `FileObj.Update()` updates all fields if the actual file has been modified since the fields were originally populated.

## Example

Here's an example of basic Objectify usage:

```go
package main

import (
    "fmt"
    "os"
    
    objf "github.com/orme292/objectify"
)

func main() {

    files, err := objf.Path("/root/dir", objf.SetsAll())
    if err != nil {
        fmt.Printf("Error occurred: %s", err.Error())
        os.Exit(1)
    }
    
    for _, entry := range files {
        fmt.Printf("%s is %d BYTES", entry.FullPath(), entry.SizeBytes)
    }
    
    os.Exit(0)
}

```
