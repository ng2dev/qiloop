package main

import (
	"io"
	"log"
	"os"
	"qiloop/meta/object"
	"qiloop/meta/proxy"
)

func main() {
	var input io.Reader
	var output io.Writer

	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		input = file
		defer file.Close()
	} else {
		input = os.Stdin
	}

	if len(os.Args) > 2 {
		filename := os.Args[2]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		output = file
		defer file.Close()
	} else {
		output = os.Stdout
	}

	metaObj, err := object.ReadMetaObject(input)
	if err != nil {
		log.Fatalf("failed to parse MetaObject: %s", err)
	}

	err = proxy.GenerateProxy(metaObj, "services", "Directory", output)
	if err != nil {
		log.Fatalf("proxy generation failed: %s\n", err)
	}
	return
}
