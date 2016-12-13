package encryption

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	key      = "somerandokeywhocares"
	testText = "This is a secret message!"
)

func TestEncrypt(t *testing.T) {

	data, err := Encrypt([]byte(key), []byte(testText))
	if err != nil {
		t.Error(err)
	}

	if data[0] != 231 || data[len(data)-1] != 77 {
		t.Error("encryption is not deterministic, will screw up git diffs")
	}

	if len(data) <= 0 {
		t.Error("no data")
	}
}

func TestEncryptFiles(t *testing.T) {

	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "terraform")

	wd = "/Users/cgrant/Documents/credo"
	err := EncryptFiles([]byte(os.Getenv("TERRAFORM_DECRYPT")), wd)
	if err != nil {
		t.Error(err)
	}
}
