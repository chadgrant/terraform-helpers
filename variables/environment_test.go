package variables

import (
	"os"
	"testing"
)

func TestImportTfvars(t *testing.T) {
	os.Setenv("TF_VAR_TEST", "test")

	vars := ImportEnvironmentVariables()

	if vars["TEST"] != "test" {
		t.Fail()
	}
}
