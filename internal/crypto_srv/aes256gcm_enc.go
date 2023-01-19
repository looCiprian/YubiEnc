package crypto_srv

import (
	"crypto/sha256"
	"errors"
	"io"
	"strconv"

	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/mac/subtle"
	"github.com/google/tink/go/streamingaead"
)

var aadSize int = 32
var macSize int = 32
var AadMacSize int = aadSize + macSize

// create a new AES256GCMHKDF1MBKey handler
// return keyset handle
func NewKeysetHandleAES256GCMHKDF1MB() (*keyset.Handle, error) {

	return keyset.NewHandle(streamingaead.AES256GCMHKDF1MBKeyTemplate())
}

func AES256GCMEncrypt(out io.Writer, kh *keyset.Handle) (io.WriteCloser, error) {

	a, err := streamingaead.New(kh)
	if err != nil {
		return nil, err
	}

	aad, _ := generateAAD(kh)

	out.Write(aad)

	w, err := a.NewEncryptingWriter(out, aad)
	if err != nil {
		return nil, err
	}

	return w, nil

}

func AES256GCMDecrypt(in io.Reader, kh *keyset.Handle) (io.Reader, error) {

	a, err := streamingaead.New(kh)
	if err != nil {
		return nil, err
	}

	// read AAD and its HMAC at the beginning of the file
	aad := make([]byte, AadMacSize)

	io.ReadFull(in, aad)

	// verify if the AAD is correct
	if !verifyAAD(aad, kh) {
		return nil, errors.New("[-] Wrong AAD")
	}

	r, err := a.NewDecryptingReader(in, aad)
	if err != nil {
		return nil, err
	}

	return r, nil

}

// generate random AAD and perform HMAC using primary key id as key
// return 64 byte, error
func generateAAD(key *keyset.Handle) ([]byte, error) {

	// create random value
	randomAAD := generateRandomByteLenN(aadSize)

	// get primary key ID and get is hash
	keyInfo := key.KeysetInfo().PrimaryKeyId
	hashWriter := sha256.New()
	hashWriter.Write([]byte(strconv.Itoa(int(keyInfo))))
	hashKeyInfo := hashWriter.Sum(nil)

	// create new HMAC with hash of primary key ID as key
	subtleHMAC, err := subtle.NewHMAC("SHA256", hashKeyInfo, uint32(macSize))

	mac, _ := subtleHMAC.ComputeMAC(randomAAD)

	var hashTag []byte
	hashTag = append(hashTag, randomAAD...)
	hashTag = append(hashTag, mac...)

	return hashTag, err

}

// get 64 byte coposed of aad and HMAC
// return true if the mac message is verified, otherwise false
func verifyAAD(aad []byte, key *keyset.Handle) bool {

	randomAAD := aad[:aadSize]
	mac := aad[macSize:]

	// get primary key ID and get is hash
	keyInfo := key.KeysetInfo().PrimaryKeyId
	hashWriter := sha256.New()
	hashWriter.Write([]byte(strconv.Itoa(int(keyInfo))))
	hashKeyInfo := hashWriter.Sum(nil)

	subtleHMAC, err := subtle.NewHMAC("SHA256", hashKeyInfo, uint32(macSize))
	if err != nil {
		return false
	}

	err = subtleHMAC.VerifyMAC(mac, randomAAD)

	return err == nil

}
