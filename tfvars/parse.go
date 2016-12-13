package tfvars

import (
	"io/ioutil"
	"strings"
)

func Parse(files []string) (map[string]string, error) {
	vars := make(map[string]string, 0)

	for _, f := range files {
		all, err := ioutil.ReadFile(f)
		if err != nil {
			return vars, err
		}
		lines := strings.Split(string(all), "\n")

		for _, line := range lines {
			if !strings.Contains(line, "=") || strings.HasPrefix(line, "#") {
				continue
			}
			nameval := strings.SplitN(line, "=", 2)
			if len(nameval) != 2 {
				continue
			}
			vars[strings.Trim(nameval[0], " ")] = cleanVal(strings.Trim(nameval[1], " "))
		}
	}

	return vars, nil
}

func cleanVal(val string) string {
	if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		return val[1 : len(val)-1]
	}
	return val
}
