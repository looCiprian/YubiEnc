package command

import (
	"fmt"
	"log"
	"yubienc/internal/crypto_srv"
	"yubienc/internal/file_mng"
)

type keyVerifyCommand struct {
	source string
}

func NewKeyVerify(source string) *keyVerifyCommand {

	if !file_mng.FileExists(source) {
		log.Fatal("[-] Key path not found")
	}

	return &keyVerifyCommand{
		source: source,
	}
}

func (kVCommand *keyVerifyCommand) ExecuteKeyVerify() error {

	source := kVCommand.source
	yk := crypto_srv.NewYPIV("")
	crypto_srv.LoadEncryptedKeyset(source, yk)
	yk.CloseYPIV()

	fmt.Println("[+] Verified OK")

	return nil

}
