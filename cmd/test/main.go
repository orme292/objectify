package main

import (
	"log"
	"os"

	objf "github.com/orme292/objectify"
)

func main() {

	log.Println("Starting: /Users/aorme/github/objectify/cmd/test/testdir")
	files, err := objf.Path("/Users/aorme/github/objectify/cmd/test/testdir", objf.SetsAll())
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("Starting: /Users/aorme/github/objectify/cmd/test/testdir/testsubdir")
	files, err = objf.Path("/Users/aorme/github/objectify/cmd/test/testdir/testsubdir", objf.SetsAll())
	if err != nil {
		log.Printf("%v", err)
	}

	for _, file := range files {
		file.DebugOut()
	}

	log.Println("Starting: /Users/aorme/github/objectify/cmd/test/testdir/testsubdir/file1")
	file, err := objf.File("/Users/aorme/github/objectify/cmd/test/testdir/testsubdir/file1", objf.SetsAll())
	if err != nil || file == nil {
		log.Printf("%v", err)
	} else {
		file.DebugOut()
	}

	os.Exit(0)

}
