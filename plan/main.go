package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/credomobile/terraform/cmds"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var environment string
	var stack string
	var service string
	var target string
	var destroy bool
	var apply bool

	flags := flag.NewFlagSet("plan", flag.ExitOnError)
	flags.Usage = printUsage
	flags.StringVar(&environment, "environment", os.Getenv("ENVIRONMENT"), "development|staging|production")
	//flags.StringVar(&stack, "stack", os.Getenv("STACK"), "name of stack")
	flags.StringVar(&target, "target", os.Getenv("TARGET"), "name of target provider.resource.id")
	flags.StringVar(&service, "service", os.Getenv("SERVICE"), "name of service")
	flags.BoolVar(&destroy, "destroy", false, "create a destroy plan")
	flags.BoolVar(&apply, "apply", false, "apply plan immediately")

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		return 1
	}

	if err := validateEnvironment(environment); err != nil {
		fmt.Println(err.Error())
		return 1
	}

	if len(flags.Args()) <= 0 {
		fmt.Println("stack is required")
		return 1
	}
	stack = strings.TrimRight(flags.Args()[0], "/")

	destroy = validateBoolFlag("destroy", destroy)
	apply = validateBoolFlag("apply", apply)

	if err := cmds.Plan(environment, stack, service, target, apply, destroy); err != nil {
		fmt.Println(err.Error())
		return 1
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

const helpText = `Usage: plan [options] [stack]
  plan searches recursively for tfvar files stored under /terraform root directory
  to pass into terraform as variables using a convention:

  "global.tfvars" : will be applied accross all environments
	"environment.tfvars" : will be applied to environment
	"development-private.tfvars" : will be applied to environment after decryption.
	"terraform.tfvars" : will be applied after decryption (as per the default in terraform)

	IMPORTANT! The environment-private.tfvars and terraform.tfvars are SECRETS and should not
	be stored in git. They should be encrypted before check-in and will be saved as
	filename.tfenc

	stack = the terraform "stack" or folder you are working with.
					environment var = STACK

Options:
  -environment        development,staging or production (environment var ENVIRONMENT)
  -service            the service you are deploying (optional, environment var SERVICE)
	--apply             apply the plan immedeately after planning
	--destroy           plan a destroy
  -target             the terraform target, see terraform docs

Output plan:
	The output plan will be stored in /terraform/plans and can be applied with:
	apply /terraform/plans/plan-filename.plan
`

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText)
}
