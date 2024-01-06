/*
Программа выводит дерево каталогов и файлов (если указан аргумент -f)
*/
package main

import (
	"fmt"
	"io"
	"os"
)

func getFormatRawPath(dir os.DirEntry, rawPath, indentsRoot, prefixDir string) (string, error) {
	if !dir.IsDir() {
		di, err := dir.Info()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		fs := di.Size()
		if fs == 0 {
			rawPath = fmt.Sprintf("%s%s%s (empty)\n", indentsRoot, prefixDir, dir.Name())
		} else {
			rawPath = fmt.Sprintf("%s%s%s (%db)\n", indentsRoot, prefixDir, dir.Name(), di.Size())
		}
	} else {
		rawPath = fmt.Sprintf("%s%s%s\n", indentsRoot, prefixDir, dir.Name())
	}
	return rawPath, nil
}

func drawDirTree(out io.Writer, path string, printFiles bool, iParentDir, lenParentDirs int, indentsRoot string) error {
	dirsEntry, err := os.ReadDir(path)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	dirs := new([]os.DirEntry)

	for _, file := range dirsEntry {
		if file.IsDir() || printFiles {
			*dirs = append(*dirs, file)
		}
	}

	if iParentDir == 0 && lenParentDirs == 0 {
		lenParentDirs = len(*dirs)
	}

	for i, dir := range *dirs {
		pathDir := path + string(os.PathSeparator) + dir.Name()
		var indents string
		var rawPath string

		if i < len(*dirs)-1 {
			indents += indentsRoot + "│\t"
			rawPath, err = getFormatRawPath(dir, rawPath, indentsRoot, "├───")
			if err != nil {
				return err
			}
		} else {
			indents += indentsRoot + "\t"
			rawPath, err = getFormatRawPath(dir, rawPath, indentsRoot, "└───")
			if err != nil {
				return err
			}
		}
		fmt.Fprintf(out, rawPath)

		drawDirTree(out, pathDir, printFiles, i, len(*dirs), indents)
	}

	return nil
}

const bud = `├───project
├───static
│	├───a_lorem
│	│	└───ipsum
│	├───css
│	├───html
│	├───js
│	└───z_lorem
│		└───ipsum
└───zline
	└───lorem
		└───ipsum
`

func dirTree(out io.Writer, path string, printFiles bool) error {
	var iParentDir, lenParentDirs int
	var indentsRoot string

	err := drawDirTree(out, path, printFiles, iParentDir, lenParentDirs, indentsRoot)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
