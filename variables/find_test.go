package variables

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParents(t *testing.T) {
	cd, _ := os.Getwd()

	files, err := Parents(filepath.Join(cd, "terraform", "service", "stack"), ".+.tfvars$")
	if err != nil {
		t.Error(err)
	}

	if len(files) <= 0 {
		t.Error("Found no files")
	}
}

func TestDescendents(t *testing.T) {
	cd, _ := os.Getwd()

	files, err := Descendents(filepath.Join(cd, "terraform"), ".+.tfvars$")
	if err != nil {
		t.Error(err)
	}

	if len(files) <= 0 {
		t.Error("Found no files")
	}
}
