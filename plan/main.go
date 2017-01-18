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
	var key string
	var stack string
	var service string
	var target string
	var bucket string
	var bucketPrefix string
	var destroy bool
	var apply bool

	flags := flag.NewFlagSet("plan", flag.ExitOnError)
	flags.Usage = printUsage
	flags.StringVar(&key, "key", os.Getenv("TERRAFORM_DECRYPT"), "encryption key")
	flags.StringVar(&environment, "environment", os.Getenv("ENVIRONMENT"), "development|staging|production")
	flags.StringVar(&bucket, "bucket", os.Getenv("BUCKET"), "name of s3 bucket")
	flags.StringVar(&bucketPrefix, "bucket-prefix", os.Getenv("BUCKET_PREFIX"), "prefix of bucket")
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

	if len(bucket) <= 0 && len(bucketPrefix) <= 0 {
		fmt.Println("bucket or bucket-prefix is required")
		return 1
	}

	destroy = validateBoolFlag("destroy", destroy)
	apply = validateBoolFlag("apply", apply)

	if err := cmds.Plan(key, bucket, bucketPrefix, environment, stack, service, target, apply, destroy); err != nil {
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

func getBucket(prefix, region, environment string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, region, environment)
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
  -key								decryption key
  -environment        development,staging or production (environment var ENVIRONMENT)
  -service            the service you are deploying (optional, environment var SERVICE)
	-bucket							the S3 bucket state is stored (environment var BUCKET)
	-bucket-prefix			if bucket is not passed, bucket is derived. {bucket-prefix}-{aws-region}-{environment} (environment var BUCKET_PREFIX)
	--apply             apply the plan immediately after planning
	--destroy           plan a destroy
  -target             the terraform target, see terraform docs

Output plan:
	The output plan will be stored in /terraform/plans and can be applied with:
	apply /terraform/plans/plan-filename.plan
`

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText)
}
