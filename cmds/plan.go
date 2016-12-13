package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bitbucket.org/credomobile/terraform/crypt/encryption"
	"bitbucket.org/credomobile/terraform/state"
	"bitbucket.org/credomobile/terraform/tfvars"
)

func Plan(environment, stack, service, target string, applyplan, destroy bool) error {
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Stack: %s\n", stack)
	fmt.Printf("Service: %s\n", service)
	fmt.Printf("Target: %s\n", target)
	fmt.Printf("Destroy: %v\n", destroy)
	fmt.Printf("Apply: %v\n", applyplan)

	workingDir, _ := os.Getwd()
	workingDir += string(os.PathSeparator) + stack

	fmt.Printf("Working dir: %s\n", workingDir)
	err := encryption.DecryptFiles([]byte(os.Getenv("TERRAFORM_DECRYPT")), workingDir)
	if err != nil {
		return err
	}
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

	foundFiles, err := tfvars.Parents(workingDir, "")
	if err != nil {
		return err
	}

	varfiles := make([]string, 0)
	for _, af := range foundFiles {
		if shouldInclude(environment, af) {
			rel, ferr := filepath.Rel(workingDir, af)
			if ferr != nil {
				return ferr
			}
			args = append(args, fmt.Sprintf("-var-file=%s", rel))
			varfiles = append(varfiles, af)
		}
	}

	vars, err := tfvars.Parse(varfiles)
	if err != nil {
		return fmt.Errorf("Error parsing files: %s", err.Error())
	}

	tfvars.ExportTfvars(vars)

	fmt.Println(strings.Join(args, " "))

	bucketExists, err := state.Configure(vars["aws_region"], environment, service, stack)
	if err != nil {
		return fmt.Errorf("Error configuring remote state %s", err.Error())
	}

	wd, _ := os.Getwd()
	os.Chdir(workingDir)
	if err := beforePlan(); err != nil {
		return err
	}

	exec.Command("terraform", "get").Run()
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}

	if err := afterPlan(); err != nil {
		return err
	}
	os.Chdir(wd)

	if applyplan {
		err = Apply(outFile(namespace, destroy), environment, stack, service, target, false, !bucketExists, destroy)
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

func beforePlan() error {
	return runIfExists("before-plan.sh")
}

func afterPlan() error {
	return runIfExists("after-plan.sh")
}
