package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func printDir(out io.Writer, path string , printFiles bool, prefix string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic("Failed to read directory")
	}
	filteredFiles := make([]os.FileInfo, 0)
	for _, f := range files {
		if printFiles || f.IsDir() {
			filteredFiles = append(filteredFiles, f)
		}
	}
	sort.Slice(filteredFiles, func(i,j int)bool {return filteredFiles[i].Name() < filteredFiles[j].Name()})
	numFiles := len(filteredFiles)
	nextPrefix := prefix + "│\t"
	for i, f := range filteredFiles {
		ch := "├───"
		if i == numFiles - 1 {
			ch = "└───"
			nextPrefix = prefix + "\t"
		}
		fmt.Fprintf(out, prefix + ch + f.Name())
		if !f.IsDir() {
			if f.Size() != 0 {
				fmt.Fprintf(out, " (%vb)\n", f.Size())
			} else {
				fmt.Fprintf(out, " (empty)\n")
			}
		} else {
			fmt.Fprintf(out,"\n")
			abs:= filepath.Join(path, f.Name())
			printDir(out, abs, printFiles, nextPrefix)
		}
	}
	return nil
}


func dirTree(out io.Writer, path string , printFiles bool) error  {
	err := printDir(out, path, printFiles, "")
	if err != nil {
		panic("Failed")
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
