package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func tree(w *bufio.Writer, path string, printFiles bool, prefix string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		onlyDirs := make([]os.DirEntry, 0)
		for _, d := range dir {
			if d.IsDir() {
				onlyDirs = append(onlyDirs, d)
			}
		}
		dir = nil
		dir = onlyDirs
	}
	var lastIdx = len(dir) - 1
	for i, entry := range dir {
		if _, err := w.WriteString(prefix); err != nil {
			break
		}

		var ch rune = '├'
		var newPrefix string = prefix
		if i == lastIdx {
			ch = '└'
			newPrefix += "\t"
		} else {
			newPrefix += "│\t"
		}
		if _, err := w.WriteRune(ch); err != nil {
			break
		}

		var info string

		if !entry.IsDir() {
			f, _ := entry.Info()
			size := f.Size()
			if size == 0 {
				info = " (empty)"
			} else {
				info = fmt.Sprintf(" (%db)", size)
			}
		}

		if _, err := w.WriteString("───" + entry.Name() + info + "\n"); err != nil {
			break
		}

		if entry.IsDir() {
			err = tree(w, path+string(os.PathSeparator)+entry.Name(), printFiles, newPrefix)
			if err != nil {
				break
			}
		}

		if err = w.Flush(); err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func dirTree(w io.Writer, path string, printFiles bool) error {
	bufW := bufio.NewWriter(w)
	err := tree(bufW, path, printFiles, "")
	if err == nil {
		err = bufW.Flush()
	}
	return err
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
