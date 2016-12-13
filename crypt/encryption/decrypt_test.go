package encryption

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecrypt(t *testing.T) {

	data, err := Encrypt([]byte(key), []byte(testText))
	if err != nil {
		t.Error(err)
	}

	if len(data) <= 0 {
		t.Error("no data")
	}

	dec, err := Decrypt([]byte(key), data)
	if err != nil {
		t.Error(err)
	}

	s := string(dec)

	if len(s) != len(testText) {
		t.Errorf("string length did not match got %d expected %d", len(s), len(testText))
	}

	if s != testText {
		t.Errorf("Strings did not match : %s", dec)
	}
}

func TestDecryptFiles(t *testing.T) {

	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "terraform", "service", "stack")

	err := DecryptFiles([]byte(key), wd)
	if err != nil {
		t.Error(err)
	}
}
