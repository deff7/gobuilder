package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deff7/gobuilder/internal/pkg/generator"
	"github.com/deff7/gobuilder/internal/pkg/parser"
)

func main() {
	var (
		file           string
		allowedStructs string
	)

	{
		flag.StringVar(&file, "f", "", "file with structure declaration")
		flag.StringVar(&allowedStructs, "s", "*", "structs list for which generate builders. * - generate for all structs")
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

	p := parser.NewParser()

	structsList := []string{}
	if allowedStructs != "*" {
		structsList = strings.Split(allowedStructs, ",")
	}

	packageName, structs, err := p.Parse(f, structsList)
	checkError(err)

	for _, s := range structs {
		g, err := generator.NewGenerator(s.Name, packageName)
		checkError(err)

		for _, f := range s.Fields {
			g.AddField(f.Name, f.TypeName)
		}

		res, err := g.Generate()
		checkError(err)

		fmt.Println(res)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
