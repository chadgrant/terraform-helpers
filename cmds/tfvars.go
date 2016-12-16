package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/chadgrant/terraform-helpers/crypt/encryption"
	"github.com/chadgrant/terraform-helpers/variables"
)

func TFVars(dir, environment string) (map[string]string, []string, error) {

	err := encryption.DecryptFiles([]byte(os.Getenv("TERRAFORM_DECRYPT")), dir)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not decrypt terraform files %s", err.Error())
	}

	varFiles, err := variables.Parents(dir, ".+\\.tfvars$")
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting var files %s", err.Error())
	}

	varFiles = filter(varFiles, func(f string) bool {
		return shouldInclude(environment, f)
	})

	defaultFiles, err := variables.Parents(dir, ".+\\.tf$")
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting default vars %s", err.Error())
	}

	defaults, err := variables.ParseTerraformFiles(defaultFiles...)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing tf files %s", err.Error())
	}

	vars, err := variables.ParseVarFiles(varFiles...)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing varfiles %s", err.Error())
	}

	envs := variables.ImportEnvironmentVariables()
	for k, v := range envs {
		defaults[k] = v
	}

	for k, v := range vars {
		defaults[k] = v
	}

	variables.ExportEnvironmentVariables(defaults)

	return defaults, varFiles, nil
}

func shouldInclude(env, af string) bool {
	paths := strings.Split(af, string(os.PathSeparator))
	f := paths[len(paths)-1]

	if strings.Contains(f, strings.ToLower(env)) && strings.HasSuffix(f, ".tfvars") {
		return true
	}

	if f == "global.tfvars" || f == "terraform.tfvars" || f == "private.tfvars" {
		return true
	}

	return false
}

func filter(files []string, include func(string) bool) []string {
	ret := make([]string, 0)
	for _, f := range files {
		if include(f) {
			ret = append(ret, f)
		}
	}
	return ret
}
