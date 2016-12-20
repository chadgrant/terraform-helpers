package variables

import (
	"os"
	"testing"
)

func TestReplace(t *testing.T) {

	wd, _ := os.Getwd()

	err := Replace(wd, "*.txt", "\\${\\s*var\\.tag_prefix\\s*}", "TAGPREFIX")
	if err != nil {
		t.Error(err)
	}

}
