package file_mng

import (
	"os"
	"path/filepath"
)

func FileExists(destination string) bool {
	if _, err := os.Stat(destination); err == nil {
		return true
	}
	return false
}

func ReadFile(path string) []byte {

	content, err := os.ReadFile(path)

	if err != nil {
		return []byte{}
	}

	return content
}

func CreateAndWriteNewFile(path string, content []byte) error {

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	f.Write(content)
	f.Close()

	return nil
}

func IsStringFile(name string) bool {

	ext := filepath.Ext(name)

	if ext == "." {
		return false
	}

	return len(ext) != 0
}

func GetFileSize(filePath string) int64 {

	fi, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	// get the size
	size := fi.Size()

	return size
}

func GetFileNameFromPath(path string) string {

	return filepath.Base(path)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
