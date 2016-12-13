package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chadgrant/terraform-helpers/tfvars"
)

func Decrypt(key, data []byte) ([]byte, error) {

	b, err := aes.NewCipher(getKey(key))
	if err != nil {
		return nil, err
	}

	padder := NewPkcs7Padding(b.BlockSize())
	cbc := cipher.NewCBCDecrypter(b, iv)
	cbc.CryptBlocks(data, data)

	unpadded, err := padder.Unpad(data)
	if err != nil {
		return nil, err
	}

	return unpadded, nil
}

func DecryptFiles(key []byte, path string) error {
	files, err := tfvars.Parents(path, ".+\\.enc$")
	if err != nil {
		return err
	}

	for _, f := range files {

		fmt.Printf("Decrypting: %s\n", f)

		b64, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		data, err := base64.StdEncoding.DecodeString(string(b64))
		if err != nil {
			return err
		}

		dec, err := Decrypt(key, data)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(strings.Replace(f, ".enc", "", 1), dec, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}
