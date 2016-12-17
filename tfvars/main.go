package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chadgrant/terraform-helpers/cmds"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var environment string

	flags := flag.NewFlagSet("tfvars", flag.ExitOnError)
	flags.Usage = printUsage
	flags.StringVar(&environment, "environment", os.Getenv("ENVIRONMENT"), "development|staging|production")

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		return 1
	}

	if len(environment) <= 0 {
		fmt.Println("environment is required")
		return 1
	}

	if err := validateEnvironment(environment); err != nil {
		fmt.Println(err.Error())
		return 1
	}

	dir, _ := os.Getwd()
	if len(flags.Args()) > 0 {
		dir = strings.TrimRight(flags.Args()[0], "/")
	}

	vars, _, err := cmds.TFVars(dir, environment)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	for k, v := range vars {
		fmt.Printf("%s=%s\n", k, v)
	}

	return 0
}

func validateEnvironment(env string) error {
	if len(env) <= 0 {
		return errors.New("environment required")
	}

	if env != "development" && env != "staging" && env != "production" {
		return errors.New("unknown environment: " + env)
	}

	return nil
}

func validateBoolFlag(name string, flagval bool) bool {
	if flagval {
		return true
	}

	for _, v := range os.Args {
		if strings.Contains(v, fmt.Sprintf("--%s", name)) {
			return true
		}
	}

	return false
}

const helpText = `Usage: tfvars [options] [dir]
  tfvars searches recursively up for tfvar and tf files stored under [dir] using a convention:

  "global.tfvars" : will be found accross all environments
	"environment.tfvars" : will be found for environment
	"development-private.tfvars" : will be found for environment and decrypted
	"terraform.tfvars" : will be applied after decryption (as per the default in terraform)

	IMPORTANT! The environment-private.tfvars and terraform.tfvars are SECRETS and should not
	be stored in git. They should be encrypted before check-in and will be saved as
	filename.tfenc

	dir = dir to start search (defaults to current directory)

Options:
  -environment        development,staging or production (environment var ENVIRONMENT)

Output key=value
`

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText)
}
