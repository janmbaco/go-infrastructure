package disk

import (
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

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
