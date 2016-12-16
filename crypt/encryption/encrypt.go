package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/chadgrant/terraform-helpers/variables"
)

//Using a static IV for deterministic outputs for git - hooks,
//not as secure
var (
	iv                 = []byte("81asqwuu786asdas")
	encryptFileFilters = []*regexp.Regexp{
		regexp.MustCompile("^terraform\\.tfvars$"),
		regexp.MustCompile("^private\\.tfvars$"),
		regexp.MustCompile(".+-private\\.tfvars$"),
		regexp.MustCompile(".+\\.pem$"),
	}
)

func getKey(key []byte) []byte {
	for len(key) < 32 {
		key = append(key, key...)
	}
	return key[:32]
}

func Encrypt(key, data []byte) ([]byte, error) {

	b, err := aes.NewCipher(getKey(key))
	if err != nil {
		return nil, err
	}

	padded, err := NewPkcs7Padding().Pad(data)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(b, iv)
	dest := make([]byte, len(padded))
	cbc.CryptBlocks(dest, padded)

	return dest, nil
}

func EncryptFiles(key []byte, path string) error {
	files, err := variables.Descendents(path, ".+\\.tfvars$|.+\\.pem$")
	if err != nil {
		return err
	}

	for _, f := range files {
		if !shouldEncrypt(f) {
			continue
		}

		data, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			continue
		}

		file := fmt.Sprintf("%s.enc", f)

		if _, sterr := os.Stat(file); sterr == nil {
			dec, derr := DecryptFile(key, file)
			if derr == nil {
				if areEqual(dec, data) {
					continue
				}
			}
		}

		fmt.Printf("Encrypting: %s\n", file)

		enc, err := Encrypt(key, data)
		if err != nil {
			return err
		}

		os.Remove(file)

		b64 := base64.StdEncoding.EncodeToString(enc)

		err = ioutil.WriteFile(file, []byte(b64), 0666)
		if err != nil {
			return fmt.Errorf("Error writing file %s.enc %s", f, err.Error())
		}
	}
	return nil
}

func areEqual(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func shouldEncrypt(path string) bool {
	_, file := filepath.Split(path)

	for _, re := range encryptFileFilters {
		if re.MatchString(file) {
			return true
		}
	}
	return false
}
