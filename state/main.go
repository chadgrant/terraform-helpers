package state

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var env = map[string]string{
	"development": "dev",
	"staging":     "stg",
	"production":  "prd",
}

func Configure(bucket, bucketPrefix, region, environment, service, stack string) error {

	if len(bucket) <= 0 {
		b, err := getBucket(bucketPrefix, region, environment)
		if err != nil {
			return err
		}
		bucket = b
	}
	bucketDir := getBucketDir(environment, service, stack)

	args := []string{
		"remote",
		"config",
		"-backend=S3",
		fmt.Sprintf("-backend-config=bucket=%s", bucket),
		fmt.Sprintf("-backend-config=key=%s", bucketDir),
	}

	fmt.Printf("Configuring remote state as S3://%s%s\n", bucket, bucketDir)

	if err := runTerraformCmd(getWorkingDir(service, stack), args); err != nil {
		return err
	}

	return nil
}

func Pull(service, stack string) error {
	return runTerraformCmd(getWorkingDir(service, stack), []string{"remote", "pull"})
}

func Push(service, stack string) error {
	return runTerraformCmd(getWorkingDir(service, stack), []string{"remote", "push"})
}

func runTerraformCmd(directory string, args []string) error {
	wd, _ := os.Getwd()
	os.Chdir(directory)

	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	os.Chdir(wd)

	return err
}

func getBucket(prefix, region, environment string) (string, error) {
	if short, ok := env[environment]; ok {
		environment = short
	}

	b := fmt.Sprintf("%s-%s-%s", prefix, region, environment)
	exists, err := bucketExists(b)
	if err != nil {
		return "", err
	}
	if exists {
		return b, nil
	}

	b = fmt.Sprintf("%s-%s-%s", prefix, environment, region)
	exists, err = bucketExists(b)
	if err != nil {
		return "", err
	}
	if exists {
		return b, nil
	}

	return "", fmt.Errorf("No bucket for env: %s region: %s", environment, region)
}

func getBucketDir(environment, service, stack string) string {
	bdir := []string{"/terraform", environment}
	if len(service) > 0 {
		bdir = append(bdir, service)
	}
	bdir = append(bdir, stack)
	return strings.Join(bdir, "/") + ".tfstate"
}

func getWorkingDir(service, stack string) string {
	dir := []string{"/terraform"}
	if len(service) > 0 {
		dir = append(dir, service)
	}
	dir = append(dir, stack)

	return path.Join(dir...)
}

func bucketExists(name string) (bool, error) {
	sess, err := session.NewSession()
	if err != nil {
		return false, fmt.Errorf("bucketExists : %s", err)
	}
	creds := credentials.NewEnvCredentials()
	svc := s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_DEFAULT_REGION")).WithCredentials(creds))
	params := &s3.GetBucketLocationInput{Bucket: aws.String(name)}
	_, err = svc.GetBucketLocation(params)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucket") {
			return false, nil
		}
		return false, fmt.Errorf("bucketExists : Location %s", err.Error())
	}

	return true, nil
}
