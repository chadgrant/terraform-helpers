package variables

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTerraformFile(t *testing.T) {

	cd, _ := os.Getwd()

	parsed, err := ParseTerraformFiles(filepath.Join(cd, "test.tf"))
	if err != nil {
		t.Error(err)
	}

	assertValue(t, parsed, "this_is_a_var", "Default value")
	assertValue(t, parsed, "this_is_a_number", "2")
	assertValue(t, parsed, "this_is_a_bool", "false")
	assertValue(t, parsed, "this_is_a_list", "value1,value1")
	assertValue(t, parsed, "this_is_a_map.key1", "val1")
	assertValue(t, parsed, "this_is_a_map.key2", "val2")
}

func TestParseKeyValueFile(t *testing.T) {

	cd, _ := os.Getwd()

	parsed, err := ParseVarFiles(filepath.Join(cd, "test.tfvars"))
	if err != nil {
		t.Error(err)
	}

	assertValue(t, parsed, "aws_region", "us-west-2")
}

func assertValue(t *testing.T, m map[string]string, k, v string) {
	if m[k] != v {
		t.Errorf("Expected %s=%s but got %s", k, v, m[k])
	}
}
