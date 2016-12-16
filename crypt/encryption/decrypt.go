package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chadgrant/terraform-helpers/variables"
)

func Decrypt(key, data []byte) ([]byte, error) {

	b, err := aes.NewCipher(getKey(key))
	if err != nil {
		return nil, err
	}

	padder := NewPkcs7Padding()
	cbc := cipher.NewCBCDecrypter(b, iv)
	cbc.CryptBlocks(data, data)

	unpadded, err := padder.Unpad(data)
	if err != nil {
		return nil, err
	}

	return unpadded, nil
}

func DecryptFiles(key []byte, path string) error {
	files, err := variables.Parents(path, ".+\\.enc$")
	if err != nil {
		return err
	}

	for _, f := range files {

		fmt.Printf("Decrypting: %s\n", f)

		dec, err := DecryptFile(key, f)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(strings.Replace(f, ".enc", "", 1), dec, 0666)
		if err != nil {
			return fmt.Errorf("Writing encrypted file %s", err.Error())
		}
	}

	return nil
}

func DecryptFile(key []byte, f string) ([]byte, error) {
	b64, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("Could not read file %s %s", f, err.Error())
	}

	data, err := base64.StdEncoding.DecodeString(string(b64))
	if err != nil {
		return nil, fmt.Errorf("Decoding base64 %s", err.Error())
	}

	return Decrypt(key, data)
}
