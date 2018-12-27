package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/deff7/gobuilder/internal/pkg/generator"
	"github.com/deff7/gobuilder/internal/pkg/parser"
)

func parseDirTree(p *parser.Parser, root string, allowedStructs []string) (map[string][]parser.StructDecl, error) {
	result := map[string][]parser.StructDecl{}

	var parse func(dir string) error
	parse = func(dir string) error {
		packages, err := p.ParseDir(dir, allowedStructs)
		if err != nil {
			return err
		}

		for packageName, structs := range packages {
			if s, ok := result[packageName]; ok {
				result[packageName] = append(s, structs...)
				continue
			}

			result[packageName] = structs
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() {
				continue
			}

			err := parse(filepath.Join(dir, f.Name()))
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := parse(root)
	return result, err
}

func main() {
	var (
		dir            string
		file           string
		allowedStructs string
		fieldsFilter   string
		recursive      bool
	)

	{
		flag.StringVar(&dir, "d", "", "directory with go files")
		flag.BoolVar(&recursive, "r", false, "parse directories recursively")
		flag.StringVar(&file, "f", "", "file with structure declaration")
		flag.StringVar(&allowedStructs, "s", "*", "structs list for which generate builders. * - generate for all structs")
		flag.StringVar(&fieldsFilter, "fields-filter", "", "specify regexp for skipping field names")
		flag.Parse()
	}

	if dir == "" && file == "" {
		fmt.Println("file is empty")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	checkError(err)

	p := parser.NewParser()

	structsList := []string{}
	if allowedStructs != "*" {
		structsList = strings.Split(allowedStructs, ",")
	}

	var packages map[string][]parser.StructDecl

	// TODO: refactor
	if dir != "" {
		dir = filepath.Join(wd, dir)
		if recursive {
			packages, err = parseDirTree(p, dir, structsList)
		} else {
			packages, err = p.ParseDir(dir, structsList)
		}
	} else {

		f, err := os.Open(filepath.Join(wd, file))
		checkError(err)
		defer f.Close()

		packages, err = p.Parse(f, structsList)
	}
	checkError(err)

	for packageName, structs := range packages {
		for _, s := range structs {
			g, err := generator.NewGenerator(s.Name, packageName, fieldsFilter)
			checkError(err)

			for _, f := range s.Fields {
				g.AddField(f.Name, f.TypeName)
			}

			res, err := g.Generate()
			checkError(err)

			fmt.Println(res)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
