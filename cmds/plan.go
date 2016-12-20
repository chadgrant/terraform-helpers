package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chadgrant/terraform-helpers/state"
	"github.com/chadgrant/terraform-helpers/variables"
)

const terraformRoot = "/terraform"

func Plan(bucket, bucketPrefix, environment, stack, service, target string, applyplan, destroy bool) error {
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Stack: %s\n", stack)
	fmt.Printf("Service: %s\n", service)
	fmt.Printf("Target: %s\n", target)
	fmt.Printf("Destroy: %v\n", destroy)
	fmt.Printf("Apply: %v\n", applyplan)

	workingDir, _ := os.Getwd()
	workingDir += string(os.PathSeparator) + stack

	fmt.Printf("Working dir: %s\n", workingDir)

	namespace := []string{environment}
	args := []string{"plan", "-module-depth=1"}

	if len(target) > 0 {
		args = append(args, fmt.Sprintf("-target=%s", target))
	}

	if len(service) > 0 {
		namespace = append(namespace, service)
		os.Setenv("TF_VAR_service_name", service)
	}
	namespace = append(namespace, stack)

	if destroy {
		args = append(args, "--destroy")
	}
	args = append(args, fmt.Sprintf("-out=%s", outFile(namespace, destroy)))

	vars, varfiles, err := TFVars(workingDir, environment)
	if err != nil {
		return err
	}

	if val, ok := vars["tag_prefix"]; ok {
		variables.Replace(terraformRoot, "*.tf", "\\${\\s*var\\.tag_prefix\\s*}", val)
	}

	for _, af := range varfiles {
		rel, ferr := filepath.Rel(workingDir, af)
		if ferr != nil {
			return fmt.Errorf("Could not make file relative %s", err.Error())
		}
		args = append(args, fmt.Sprintf("-var-file=%s", rel))
	}

	fmt.Println(strings.Join(args, " "))

	bucketExists, err := state.Configure(bucket, bucketPrefix, vars["aws_region"], environment, service, stack)
	if err != nil {
		return fmt.Errorf("Error configuring remote state %s", err.Error())
	}

	wd, _ := os.Getwd()
	os.Chdir(workingDir)
	if berr := beforePlan(); berr != nil {
		return fmt.Errorf("Error before plan %s", berr.Error())
	}

	exec.Command("terraform", "get").Run()
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}

	if aerr := afterPlan(); aerr != nil {
		return fmt.Errorf("Error after plan %s", aerr.Error())
	}
	os.Chdir(wd)

	if applyplan {
		err = Apply(bucket, bucketPrefix, outFile(namespace, destroy), environment, stack, service, target, false, !bucketExists, destroy)
		if err != nil {
			return fmt.Errorf("Error applying plan: %s", err.Error())
		}
	}

	return nil
}

func outFile(ns []string, destroy bool) string {
	c := ns
	if destroy {
		c = append([]string{"destroy"}, ns...)
	}
	return fmt.Sprintf("/terraform/plans/%s.plan", strings.Join(c, "_"))
}

func beforePlan() error {
	return runIfExists("before-plan.sh")
}

func afterPlan() error {
	return runIfExists("after-plan.sh")
}
