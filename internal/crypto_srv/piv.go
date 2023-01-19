package crypto_srv

import (
	"crypto"
	"log"
	"strings"
	"yubienc/internal/input_mng"

	"crypto/rand"
	"crypto/rsa"

	"github.com/go-piv/piv-go/piv"
)

type YubikeyMng struct {
	yk   *piv.YubiKey
	pub  crypto.PublicKey
	priv crypto.PrivateKey
}

func NewYPIV(pin string) *YubikeyMng {

	myYubikeyMng := YubikeyMng{}
	yk := getYubikey()
	if yk == nil {
		log.Fatal("[-] No Yubikey detected")
	}

	myYubikeyMng.yk = yk
	priv, pub := getYubikeyPrivPubKey(yk, pin)
	myYubikeyMng.priv = priv
	myYubikeyMng.pub = pub

	return &myYubikeyMng

}

func (ym *YubikeyMng) CloseYPIV() {

	ym.yk.Close()
	ym.priv = nil
	ym.pub = nil

}

func getYubikey() *piv.YubiKey {
	cards, err := piv.Cards()
	if err != nil {
		return nil
	}

	// Find a YubiKey and open the reader.
	var yk *piv.YubiKey
	for _, card := range cards {
		if strings.Contains(strings.ToLower(card), "yubikey") {
			if yk, err = piv.Open(card); err != nil {
				return nil
			}
			break
		}
	}

	if yk == nil {
		return nil
	}

	return yk
}

// get Yubikey priv and pub key
// if pin == "" ask for stdin else use pin
func getYubikeyPrivPubKey(yk *piv.YubiKey, pin string) (crypto.PrivateKey, crypto.PublicKey) {

	pubCert, err := yk.Certificate(piv.SlotKeyManagement)

	if err != nil {
		log.Fatal("[-] Cannot get Yubikey certificate")
	}

	mypublicKey := pubCert.PublicKey.(crypto.PublicKey)

	yKeyAuth := piv.KeyAuth{

		PINPolicy: piv.PINPolicyOnce,
	}

	if pin == "" {
		yKeyAuth.PINPrompt = func() (pin string, err error) {
			inputPin, err := input_mng.InputRequest("Enter the secret PIN: ", true)
			if err != nil {
				log.Fatal("[-] Cannot read Yubikey PIN")
			}
			return inputPin, nil
		}
	} else {
		yKeyAuth.PIN = pin
	}

	priv, err := yk.PrivateKey(piv.SlotKeyManagement, mypublicKey, yKeyAuth)

	if err != nil {
		log.Fatal("[-] Cannot get private and public key from Yubikey")
	}

	return priv, mypublicKey

}

func (ym YubikeyMng) PivDecrypt(encryptedMsg []byte) []byte {

	priv := ym.priv

	decrypter, ok := priv.(crypto.Decrypter)
	if !ok {
		log.Fatal("Cannot create decrypter")
	}

	rng := rand.Reader
	plainText, err := decrypter.Decrypt(rng, encryptedMsg, nil)

	if err != nil {
		log.Fatal("[-] Cannot decrypt the symmetric key using Yubikey. Check your Yubikey PIN")
	}

	return plainText

}

func (ym YubikeyMng) PivEncrypt(plainText []byte) []byte {

	pub := ym.pub

	rng := rand.Reader
	encrypted, err := rsa.EncryptPKCS1v15(rng, pub.(*rsa.PublicKey), plainText)

	if err != nil {
		log.Fatal("[-] Cannot encrypt the symmetric key using Yubikey")
	}

	return encrypted

}
