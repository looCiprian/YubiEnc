package command

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"yubienc/internal"
	"yubienc/internal/crypto_srv"
	"yubienc/internal/file_mng"

	"github.com/google/tink/go/keyset"
	"github.com/schollz/progressbar/v3"
)

type decryptCommand struct {
	destination   string
	source        string
	keyPath       string
	threatDecrypt int
	keySet        *keyset.Handle
	isDir         bool
}

func NewDecryption(source string, destination string, keyPath string, threatDecrypt int) *decryptCommand {

	if file_mng.IsStringFile(destination) {
		log.Fatal("[-] Specify directory not file in <destination>")
	}

	destination = file_mng.AddSlashIfNotPresent(destination)

	isDir := file_mng.IsDir(source)

	if isDir {
		source = file_mng.AddSlashIfNotPresent(source)
	}

	if !file_mng.FileExists(keyPath) {
		log.Fatal("Key file not exists")
	}

	ym := crypto_srv.NewYPIV("")
	keySet := crypto_srv.LoadEncryptedKeyset(keyPath, ym)
	ym.CloseYPIV()

	if file_mng.CreateDirectoryIfNotExists(destination) {
		fmt.Println("Created a new decryption dir:", destination)
	}

	return &decryptCommand{
		destination:   destination,
		source:        source,
		keyPath:       keyPath,
		keySet:        keySet,
		threatDecrypt: threatDecrypt,
		isDir:         isDir,
	}

}

func (decCommand *decryptCommand) ExecuteDecrypt() error {

	destination := decCommand.destination
	source := decCommand.source
	keySet := decCommand.keySet
	theadPool := decCommand.threatDecrypt

	if !decCommand.isDir {
		destination = destination + file_mng.GetFileNameFromPath(source)[:len(file_mng.GetFileNameFromPath(source))-len(internal.Encryption_extension)]
		alreadyDecrypted := file_mng.FileExists(destination)
		if alreadyDecrypted {
			fmt.Println("Skipping:", source, "already decrypted in", destination)
			return nil
		}
		err := fileDecrypt(source, destination, keySet)
		return err
	}

	sem := make(chan int, theadPool)

	// Loop all files from source
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		sem <- 0
		go func() {
			if !info.IsDir() {
				filePathSource := path // /tmp/source/1/1.txt.enc

				relativeFilePath := filePathSource[len(source):] // 1/1.txt.enc

				destinationPath := destination + relativeFilePath[:len(relativeFilePath)-len(internal.Encryption_extension)] // /tmp/dest/1/1.txt
				alreadyDecrypted := file_mng.FileExists(destinationPath)

				// check if file has been already decrypted or if needs to be override (like .DS_Store)
				if !alreadyDecrypted || toOverride(destinationPath) {

					dir := file_mng.GetDirFromFilePath(destinationPath)
					file_mng.CreateDirectoryIfNotExists(dir)

					err := fileDecrypt(filePathSource, destinationPath, keySet)

					if err != nil {
						fmt.Println("[-] ERROR decrypting file:", info.Name(), "with error:", err.Error())
					}
				} else {
					fmt.Println("Skipping:", info.Name(), "already decrypted")
				}
			}
			<-sem
		}()

		return nil
	})

	for n := theadPool; n > 0; n-- {
		sem <- 0
	}

	fmt.Println("Decryption Done!")
	return nil
}

func fileDecrypt(fileToDecryptPath string, decryptedFilePath string, keySet *keyset.Handle) error {

	infoFileName := file_mng.GetFileNameFromPath(fileToDecryptPath)
	size := file_mng.GetFileSize(fileToDecryptPath)

	ctFile, err := os.Open(fileToDecryptPath)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(decryptedFilePath)
	if err != nil {
		return err
	}

	r, err := crypto_srv.AES256GCMDecrypt(ctFile, keySet)
	if err != nil {
		file_mng.DeleteFile(decryptedFilePath)
		return err
	}

	bar := progressbar.DefaultBytes(size, "Decrypting: "+infoFileName)

	if _, err := io.Copy(io.MultiWriter(dstFile, bar), r); err != nil {
		return err
	}

	bar.Finish()

	if err := dstFile.Close(); err != nil {
		return err
	}
	if err := ctFile.Close(); err != nil {
		return err
	}

	return nil
}

func toOverride(filePath string) bool {

	for _, fileName := range internal.OverrideFiles {

		if filepath.Base(filePath) == fileName {
			return true
		}
	}

	return false

}
