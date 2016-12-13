package tfvars

import (
	"os"
	"testing"
)

func TestParents(t *testing.T) {
	cd, _ := os.Getwd()

	cd = "/Users/cgrant/Documents/credo/sms/sms_send_api/terraform/api"

	files, err := Parents(cd, "")
	if err != nil {
		t.Error(err)
	}

	if len(files) <= 0 {
		t.Error("Found no files")
	}

	// for _, f := range files {
	// 	fmt.Println(f)
	// }
}

func TestDescendents(t *testing.T) {
	cd, _ := os.Getwd()

	cd = "/Users/cgrant/Documents/credo"

	files, err := Descendents(cd, "")
	if err != nil {
		t.Error(err)
	}

	if len(files) <= 0 {
		t.Error("Found no files")
	}

	// for _, f := range files {
	// 	fmt.Println(f)
	// }
}
