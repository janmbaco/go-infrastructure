package disk

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ExistsPath return if path exits
func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateFile create a file
func CreateFile(filePath string, content []byte) {
	fc, err := os.Create(filePath)
	defer func() {
		_ = fc.Close()
	}()

	if err == nil {
		_, _ = fc.Write(content)
	} else {
		panic(err)
	}
}

// Copy a file
func Copy(src, dst string) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		panic(err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		panic(errors.New(fmt.Sprintf("%s is not a regular file", src)))
	}

	source, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer func() { _ = destination.Close() }()
	_, err = io.Copy(destination, source)
	if err != nil {
		panic(err)
	}
}

// DeleteFile deletes a file
func DeleteFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}
