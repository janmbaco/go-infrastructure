package disk

import (
	"fmt"
	"io"
	"os"
)

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CreateFile(filePath string, content []byte) error {
	fc, err := os.Create(filePath)
	defer func() {
		_ = fc.Close()
	}()
	if err == nil {
		_, _ = fc.Write(content)
	}
	return err
}

func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destination.Close() }()
	_, err = io.Copy(destination, source)
	return err
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
