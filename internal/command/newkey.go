package command

import (
	"log"
	"yubienc/internal/crypto_srv"
	"yubienc/internal/file_mng"
)

type newkeyCommand struct {
	destination string
}

func NewNewkey(destination string) *newkeyCommand {

	if file_mng.IsStringFile(destination) {
		log.Fatal("[-] Specify directory not file")
	}

	return &newkeyCommand{
		destination: destination,
	}
}

func (newkeyCommand *newkeyCommand) ExecuteNewkey() error {

	destination := newkeyCommand.destination
	yk := crypto_srv.NewYPIV("")
	crypto_srv.GenerateAndSaveNewKeyset(destination, yk)
	yk.CloseYPIV()

	return nil

}
