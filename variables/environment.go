package variables

import (
	"os"
	"strings"
)

func ExportEnvironmentVariables(vars map[string]string) {
	for k, v := range vars {

		if k == "aws_access_key" {
			os.Setenv("AWS_ACCESS_KEY", v)
		}

		if k == "aws_secret_key" {
			os.Setenv("AWS_SECRET_ACCESS_KEY", v)
		}

		if k == "aws_region" {
			os.Setenv("AWS_DEFAULT_REGION", v)
		}

		os.Setenv(k, v)
	}
}

func ImportEnvironmentVariables() map[string]string {
	vars := make(map[string]string, 0)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) > 0 && strings.HasPrefix(pair[0], "TF_VAR_") {
			vars[strings.Replace(pair[0], "TF_VAR_", "", -1)] = pair[1]
		}
	}

	return vars
}
