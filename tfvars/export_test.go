package tfvars

import (
	"os"
	"testing"
)

func TestImportTfvars(t *testing.T) {
	os.Setenv("TF_VAR_TEST", "test")

	vars := importTfvars()

	// for k, v := range vars {
	// 	fmt.Printf("%s=%s", k, v)
	// }

	if vars["TEST"] != "test" {
		t.Fail()
	}
}
