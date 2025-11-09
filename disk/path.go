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
func CreateFile(filePath string, content []byte) error {
	fc, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fc.Close()
	}()

	if _, err := fc.Write(content); err != nil {
		return err
	}
	return nil
}

// Copy a file
func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return errors.New(fmt.Sprintf("%s is not a regular file", src))
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
	
	if _, err = io.Copy(destination, source); err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes a file
func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
