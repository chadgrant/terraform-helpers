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
	var key = os.Getenv("TERRAFORM_DECRYPT")
	var environment = os.Getenv("ENVIRONMENT")
	var service = os.Getenv("SERVICE")
	var stack string
	var target string
	var bucket string
	var bucketPrefix string
	var destroy bool

	file := os.Args[1:][0]
	parsed, err := parsePlanFile(file)
	if err == nil {
		environment = parsed.environment
		destroy = parsed.destroy
		stack = parsed.stack
		service = parsed.service
	}

	flags := flag.NewFlagSet("apply", flag.ExitOnError)
	flags.Usage = printUsage
	flags.StringVar(&key, "key", os.Getenv("TERRAFORM_DECRYPT"), "encryption key")
	flags.StringVar(&environment, "environment", environment, "development|staging|production")
	flags.StringVar(&stack, "stack", stack, "name of stack")
	flags.StringVar(&bucket, "bucket", os.Getenv("BUCKET"), "name of s3 bucket")
	flags.StringVar(&bucketPrefix, "bucket-prefix", os.Getenv("BUCKET_PREFIX"), "prefix of bucket")
	flags.StringVar(&target, "target", "", "name of target provider.resource.id")
	flags.StringVar(&service, "service", service, "name of service")
	flags.BoolVar(&destroy, "destroy", destroy, "create a destroy plan")

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		return 1
	}

	if len(flags.Args()) <= 0 {
		fmt.Println("plan file is required")
		return 1
	}
	file = flags.Args()[0]

	if len(key) <= 0 {
		fmt.Println("decryption key required")
		return 1
	}

	if len(bucket) <= 0 && len(bucketPrefix) <= 0 {
		fmt.Println("bucket or bucket-prefix is required")
		return 1
	}

	if err := validateEnvironment(environment); err != nil {
		fmt.Println(err.Error())
		return 1
	}

	destroy = validateBoolFlag("destroy", destroy)

	if err := cmds.Apply(key, bucket, bucketPrefix, file, environment, stack, service, target, true, false, destroy); err != nil {
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

const helpText = `Usage: apply [options] [planfile]
	Options are derived by the plan filename, {destroy}_{environment}_{service}_{stack}.plan

  -key								decryption key
	-environment        development,staging or production (environment var ENVIRONMENT)
  -service            the service you are deploying (optional, environment var SERVICE)
	-stack              the stack you are deploying
	-bucket							the S3 bucket state is stored (environment var BUCKET)
	-bucket-prefix			if bucket is not passed, bucket is derived. {bucket-prefix}-{aws-region}-{environment} (environment var BUCKET_PREFIX)
	--destroy						are we destroying?

Options:
  -target             the terraform target, see terraform docs
`

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText)
}
