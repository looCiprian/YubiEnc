# ***YubiEnc***

YubiEnc allows you to encrypt/decrypt files and directories using your YubiKey.

## Available Commands

```
  decrypt     Decrypt command will decrypt the <source> file/directory to the <destination> directory
  encrypt     Encrypt command will encrypt the <source> file/directory to the <destination> directory, will create <destination> if not exists
  help        Help about any command
  newkey      Create a new key and save it encrypted in <destination>
  verifykey   Check if the encrypted keyset located in <source> can be decrypted using YubiKey
  version     Print the version number
```

### Build

```
go build yubienc.go
```

### Download binaries from release page

[Download here](https://github.com/looCiprian/YubiEnc/releases)

### Example

- Create new key
```
go run yubienc.go newkey -d test/
```

- Verify the key
```
go run yubienc.go verifykey -s test/key.yk.enc
```

- Encrypt a directory, using 20 threads (threads are used to perform multiple encryption or decryption operations at the same time)
```
go run yubienc.go encrypt -d test/encrypted -s test/toEncrypt -k test/key.yk.enc -t 20
```

- Decrypt a directory (single thread)
```
go run yubienc.go decrypt -s test/encrypted -d test/decrypted -k test/key.yk.enc
```

## Yubikey configuration

To configure your YubiKey download [Yubikey Manager](https://www.yubico.com/support/download/yubikey-manager/) and create a new self-signed certificate (RSA 2048) on slot 9d.

It is ***highly*** suggested to create a new certificate (outside YubiKey) and then import it to have a certificate backup, otherwise, if it breaks you won't be able to decrypt your files.

#### Tested on

Tested on USB-A YubiKey 5 NFC. Should work with all YubiKeys that support PIV specifications.

## Encryption model

The encryption model is based on [Hybrid cryptosystem](https://en.wikipedia.org/wiki/Hybrid_cryptosystem), specifically the [YubiKey PIV Key Management slot](https://developers.yubico.com/PIV/Introduction/Certificate_slots.html) is used as KEK to encrypt the symmetric key used as DEK to encrypt files. All files inside the source directory will be encrypted with the specified symmetric key.

The algorithms used to perform encryption/decryption operations are:
 - symmetric: [AES256GCMHKDF1MB](https://pkg.go.dev/github.com/google/tink/go/streamingaead#AES256GCMHKDF1MBKeyTemplate)
 - asymmetric: [RSA 2048 and the padding scheme from PKCS #1 v1.5](https://pkg.go.dev/crypto/rsa#EncryptPKCS1v15)

To generate the symmetric key and perform Streaming AEAD operations [Google's Tink](https://developers.google.com/tink) library has been used. Associated Data (64 bytes) is saved at the beginning of each encrypted file.

## Know bug

Progress bar does not support [multi-line output](https://github.com/schollz/progressbar/issues/6) so the progress bar will not be reliable if you use more than one thread.

