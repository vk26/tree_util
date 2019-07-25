package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := printDir(out, path, printFiles, "")
	return err
}

func printDir(out io.Writer, path string, printFiles bool, prefix string) error {

	files, err := getFilesToShow(path, printFiles)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	var filePrefix string
	var newPrefix string
	var sizeStr string

	for i, file := range files {
		if i == len(files)-1 {
			filePrefix = "└───"
			newPrefix = prefix + "\t"

		} else {
			filePrefix = "├───"
			newPrefix = prefix + "│\t"
		}
		line := prefix + filePrefix + file.Name()
		if file.IsDir() == false {
			size := file.Size()
			if size == 0 {
				sizeStr = "empty"
			} else {
				sizeStr = strconv.FormatInt(int64(size), 10) + "b"
			}
			fmt.Fprintf(out, "%v (%v)\n", line, sizeStr)
		} else {
			fmt.Fprint(out, line+"\n")
		}

		if file.IsDir() == true {
			printDir(out, path+"/"+file.Name(), printFiles, newPrefix)
		}
	}
	return nil
}

func getFilesToShow(path string, printFiles bool) ([]os.FileInfo, error) {
	var filesToShow []os.FileInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() == false && printFiles == false {
			continue
		}
		filesToShow = append(filesToShow, file)
	}
	return filesToShow, nil
}
