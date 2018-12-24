package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var (
		file    string
		structs string
	)

	{
		flag.StringVar(&file, "f", "", "file with structure declaration")
		flag.StringVar(&structs, "s", "*", "structs list for which generate builders. * - generate for all structs")
		flag.Parse()
	}

	if file == "" {
		fmt.Println("file is empty")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	checkError(err)

	f, err := os.Open(filepath.Join(wd, file))
	checkError(err)
	defer f.Close()

	p := NewParser()
	_ = f
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
