package variables

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/terraform/config"
)

func ParseVarFiles(files ...string) (map[string]string, error) {
	vars := make(map[string]string, 0)

	for _, path := range files {
		// Read the HCL file and prepare for parsing
		d, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("Error reading %s: %s", path, err)
		}

		// Parse it
		obj, err := hcl.Parse(string(d))
		if err != nil {
			return nil, fmt.Errorf("Error parsing %s: %s", path, err)
		}

		var result map[string]interface{}
		if err := hcl.DecodeObject(&result, obj); err != nil {
			return nil, fmt.Errorf(
				"Error decoding Terraform vars file: %s\n\n"+
					"The vars file should be in the format of `key = \"value\"`.\n"+
					"Decoding errors are usually caused by an invalid format.",
				err)
		}

		for k, v := range result {
			vars[k] = fmt.Sprintf("%s", v)
		}

	}

	return vars, nil
}

type variable struct {
	Key   string
	Value string
}

func ParseTerraformFiles(files ...string) (map[string]string, error) {
	vars := make(map[string]string, 0)

	for _, f := range files {
		config, err := config.LoadFile(f)
		if err != nil {
			return vars, fmt.Errorf("Unable to read file: %s\n%s", f, err.Error())
		}

		for _, tvar := range config.Variables {
			parsed := parseVariable(tvar)
			for _, v := range parsed {
				vars[v.Key] = v.Value
			}
		}
	}

	return vars, nil
}

func parseVariable(v *config.Variable) []*variable {
	vars := make([]*variable, 0)

	switch v.Type() {
	case config.VariableTypeString:
		vars = append(vars, &variable{v.Name, fmt.Sprintf("%s", v.Default)})

	case config.VariableTypeList:
		if list, ok := v.Default.([]interface{}); ok {
			vars = append(vars, &variable{v.Name, listToString(list)})
		}

	case config.VariableTypeMap:
		if m, ok := v.Default.(map[string]interface{}); ok {
			vars = append(vars, mapToStrings(v.Name, m)...)
		}

	}
	return vars
}

func listToString(list []interface{}) string {
	vals := make([]string, 0)
	for _, v := range list {
		vals = append(vals, fmt.Sprintf("%s", v))
	}
	return strings.Join(vals, ",")
}

func mapToStrings(prefix string, m map[string]interface{}) []*variable {
	vars := make([]*variable, 0)

	for k, v := range m {
		vars = append(vars, &variable{fmt.Sprintf("%s.%s", prefix, k), fmt.Sprintf("%s", v)})
	}

	return vars
}
