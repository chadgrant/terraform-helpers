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

func Configure(bucket, bucketPrefix, region, environment, service, stack string) (bool, error) {

	if len(bucket) <= 0 {
		bucket = getBucket(bucketPrefix, region, environment)
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

	//let the first push go through (create bucket)
	exists, err := bucketExists(bucket)
	if err != nil {
		return false, err
	}

	if stack == "network" && !exists {
		return false, nil
	}

	if err := runTerraformCmd(getWorkingDir(service, stack), args); err != nil {
		return true, err
	}

	return true, nil
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

func getBucket(prefix, region, environment string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, region, environment)
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
