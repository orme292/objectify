package main

import (
	"log"
	"os"

	objf "github.com/orme292/objectify"
)

func main() {

	path := "/Users/aorme/github/objectify/cmd/test/testdir"
	log.Println("Starting: ", path)
	files, err := objf.Path(path, objf.SetsAll())
	if err != nil {
		log.Printf("%v", err)
	}

	path = "/Users/aorme/github/objectify/cmd/test/testdir/testsubdir"
	log.Println("Starting: ", path)
	files, err = objf.Path(path, objf.SetsAll())
	if err != nil {
		log.Printf("%v", err)
	}

	for _, file := range files {
		file.DebugOut()
	}

	path = "/Users/aorme/github/objectify/cmd/test/testdir/testsubdir/file1"
	log.Println("Starting: ", path)
	file, err := objf.File(path, objf.SetsAll())
	if err != nil || file == nil {
		log.Printf("%v", err)
	} else {
		file.DebugOut()
	}

	os.Exit(0)

}
