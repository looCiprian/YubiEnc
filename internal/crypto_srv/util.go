package crypto_srv

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"reflect"
	"time"
	internal "yubienc/internal"
	"yubienc/internal/file_mng"
	"yubienc/internal/input_mng"

	"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/keyset"
)

// Create safe random byte lenght n
func generateRandomByteLenN(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Cannot generate random byte")
	}
	return b
}

// return JSON keyset string from keyset handle
func getKeysetByte(kh *keyset.Handle) []byte {

	var handleByte []byte
	handleBuffer := bytes.NewBuffer(handleByte)
	//keysetWriter := keyset.NewJSONWriter(handleBuffer)
	keysetWriter := keyset.NewBinaryWriter(handleBuffer)
	insecurecleartextkeyset.Write(kh, keysetWriter)

	return handleBuffer.Bytes()
}

// return the keyset handler from JSON string
func getKeysetHandleFromByte(keyJSON []byte) *keyset.Handle {

	//stringReader := strings.NewReader(keyJSON)
	binaryReader := bytes.NewReader(keyJSON)
	keysetReader := keyset.NewBinaryReader(binaryReader)

	kr, _ := insecurecleartextkeyset.Read(keysetReader)

	return kr
}

// Load and decrypt the encrypted keyset
// return keyset handle
func LoadEncryptedKeyset(keyPathInput string, ym *YubikeyMng) *keyset.Handle {

	keyContent := file_mng.ReadFile(keyPathInput)
	decryptedKey := ym.PivDecrypt(keyContent)
	keysetHandle := getKeysetHandleFromByte(decryptedKey)
	if keysetHandle == nil {
		log.Fatal("[-] Cannot load the keyset")
	}

	return keysetHandle
}

// Generate new keyset and save it encrypted
// return keyset handle
func GenerateAndSaveNewKeyset(destinationKeyPath string, ym *YubikeyMng) *keyset.Handle {

	// If no key has been specified create a new one
	// Generate keyset e key for AES256GCMHKDF1MBKeyTemplate
	keysetHandle, err := NewKeysetHandleAES256GCMHKDF1MB()
	if err != nil {
		log.Fatal("[-] Cannot create new keyset" + err.Error())
	}

	// get yk serial number and timestamp to create a unique name key
	ykSerial, err := ym.yk.Serial()
	if err != nil {
		log.Fatal("[-] Cannot get the yubikey serial number" + err.Error())
	}
	now := time.Now()
	timeStamp := now.Unix()

	destinationKeyPath = file_mng.AddSlashIfNotPresent(destinationKeyPath)
	destinationKeyPath = destinationKeyPath + fmt.Sprintf(internal.Key_name, ykSerial, timeStamp)

	if file_mng.FileExists(destinationKeyPath) {
		inputUser, _ := input_mng.InputRequest("A key already exists in: "+destinationKeyPath+" override it? [y/N]", false)

		if inputUser != "y" {
			log.Fatal("[-] Remove the existing key and try again")
		}

	}

	keyByte := getKeysetByte(keysetHandle)

	encryptedKey := ym.PivEncrypt(keyByte)
	file_mng.CreateAndWriteNewFile(destinationKeyPath, encryptedKey)

	// Check if the keyset has been correctly saved by loading and comparing it
	keyContent := file_mng.ReadFile(destinationKeyPath)
	decryptedKey := ym.PivDecrypt(keyContent)
	keySetHandleDecrypted := getKeysetHandleFromByte(decryptedKey)

	if reflect.DeepEqual(keysetHandle, keySetHandleDecrypted) {
		log.Fatal("[-] Key saved and loaded are different!!!")
	}

	return keysetHandle
}
