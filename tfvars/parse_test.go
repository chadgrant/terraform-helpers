package tfvars

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	cd, _ := os.Getwd()

	cd = "/Users/cgrant/Documents/credo/sms/sms_send_api/terraform/api"

	files, err := Parents(cd)
	if err != nil {
		t.Error(err)
	}

	if len(files) <= 0 {
		t.Error("Found no files")
	}

	_, err = Parse(files)
	if err != nil {
		t.Error(err)
	}

	// for k, v := range parsed {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
}
