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

type encryptionCommand struct {
	destination   string
	source        string
	keyPath       string
	threatEncrypt int
	keySet        *keyset.Handle
	isDir         bool
}

func NewEncryption(source string, destination string, keyPath string, threatEncrypt int) *encryptionCommand {

	if file_mng.IsStringFile(destination) {
		log.Fatal("[-] Specify directory not file in <destination>")
	}

	destination = file_mng.AddSlashIfNotPresent(destination) // /tmp/encryption/

	isDir := file_mng.IsDir(source)
	if isDir {
		source = file_mng.AddSlashIfNotPresent(source)
	}

	ym := crypto_srv.NewYPIV("")
	keysetHandle := crypto_srv.LoadEncryptedKeyset(keyPath, ym)
	ym.CloseYPIV()

	if file_mng.CreateDirectoryIfNotExists(destination) {
		fmt.Println("Created a new encryption dir:", destination)
	}

	return &encryptionCommand{
		source:        source,
		destination:   destination,
		keyPath:       keyPath,
		keySet:        keysetHandle,
		threatEncrypt: threatEncrypt,
		isDir:         isDir,
	}

}

// Execute encryption
func (encCommand *encryptionCommand) ExecuteEncrypt() error {

	destination := encCommand.destination
	source := encCommand.source
	keysetHandle := encCommand.keySet
	theadPool := encCommand.threatEncrypt

	if !encCommand.isDir {
		destination = destination + file_mng.GetFileNameFromPath(source) + internal.Encryption_extension
		alreadyDecrypted := file_mng.FileExists(destination)
		if alreadyDecrypted {
			fmt.Println("Skipping:", source, "already encrypted in", destination)
			return nil
		}
		err := fileEncrypt(source, destination, keysetHandle)
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
				fullSourcePath := path                                                    // /tmp/source/1/1.txt
				relativePath := fullSourcePath[len(source):]                              // 1/1.txt
				encFilePath := destination + relativePath + internal.Encryption_extension // /tmp/dest/1/1.txt.enc
				fileExists := file_mng.FileExists(encFilePath)

				// No error no file present in encryption
				if !fileExists {

					// create dir for file inside dir in source
					dir := file_mng.GetDirFromFilePath(encFilePath)
					file_mng.CreateDirectoryIfNotExists(dir)

					err := fileEncrypt(fullSourcePath, encFilePath, keysetHandle)

					if err != nil {
						fmt.Println("[-] Error encrypting file:", info.Name(), "with error:", err.Error())
					}
				} else if fileExists { // No error but file already exists

					fmt.Println("File", info.Name(), "already present in", encFilePath)
				}
			}
			<-sem
		}()

		return nil

	})

	for n := theadPool; n > 0; n-- {
		sem <- 0
	}

	fmt.Println("Encryption Done!")
	return nil
}

func fileEncrypt(fileToEncryptPath string, encryptedFilePath string, keysetHandle *keyset.Handle) error {

	infoFileName := file_mng.GetFileNameFromPath(fileToEncryptPath)
	size := file_mng.GetFileSize(fileToEncryptPath)

	srcFile, err := os.Open(fileToEncryptPath)
	if err != nil {
		return err
	}

	ctFile, err := os.Create(encryptedFilePath)
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(size, "Encrypting: "+infoFileName)

	w, err := crypto_srv.AES256GCMEncrypt(ctFile, keysetHandle)
	if err != nil {
		file_mng.DeleteFile(encryptedFilePath)
		return err
	}

	if _, err := io.Copy(io.MultiWriter(w, bar), srcFile); err != nil {
		return err
	}

	bar.Finish()

	if err := w.Close(); err != nil {
		return err
	}

	if err := ctFile.Close(); err != nil {
		return err
	}
	if err := srcFile.Close(); err != nil {
		return err
	}

	return nil

}
