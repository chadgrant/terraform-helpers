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

	"bitbucket.org/credomobile/terraform/tfvars"
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

	padded, err := NewPkcs7Padding(b.BlockSize()).Pad(data)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(b, iv)
	dest := make([]byte, len(padded))
	cbc.CryptBlocks(dest, padded)

	return dest, nil
}

func EncryptFiles(key []byte, path string) error {
	files, err := tfvars.Descendents(path, ".+\\.tfvars$|.+\\.pem$")
	if err != nil {
		return err
	}

	for _, f := range files {
		if !shouldEncrypt(f) {
			continue
		}
		fmt.Printf("Encrypting: %s\n", f)

		data, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			continue
		}

		enc, err := Encrypt(key, data)
		if err != nil {
			return err
		}

		nf := fmt.Sprintf("%s.enc", f)

		os.Remove(nf)

		b64 := base64.StdEncoding.EncodeToString(enc)

		werr := ioutil.WriteFile(nf, []byte(b64), 0666)
		if werr != nil {
			return werr
		}
	}
	return nil
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
