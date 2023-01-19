package file_mng

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetDirFromFilePath(path string) string {

	dir, _ := filepath.Split(path)

	return dir

}

func DirectoryExists(directory string) bool {
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		return true
	}
	return false
}

func CreateDirectoryIfNotExists(destination string) bool {
	if !DirectoryExists(destination) {
		err := os.MkdirAll(destination, 0755)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
	return false
}

func DirectoriesInPath(destination string) []os.FileInfo {

	files, err := ioutil.ReadDir(destination)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func IsEmptyDirectory(destination string) bool {
	files, err := ioutil.ReadDir(destination)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		return true
	} else if len(files) == 1 && files[0].Name() == ".DS_Store" { // Mac create .DS_Store in each directory
		return true
	}
	return false
}

func IsDir(path string) bool {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return true
	}
	return false
}
