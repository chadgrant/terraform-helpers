package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"bitbucket.org/credomobile/terraform/state"
	"bitbucket.org/credomobile/terraform/tfvars"
)

func Apply(file, environment, stack, service, target string, pullState, pushState, destroy bool) error {
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Stack: %s\n", stack)
	fmt.Printf("Service: %s\n", service)
	fmt.Printf("Target: %s\n", target)
	fmt.Printf("Destroy: %v\n", destroy)

	dir := []string{"/terraform"}

	args := []string{"apply", "-input=true"}

	if len(target) > 0 {
		args = append(args, fmt.Sprintf("-target=%s", target))
	}
	args = append(args, file)

	if len(service) > 0 {
		dir = append(dir, service)
	}
	dir = append(dir, stack)

	workingDir := path.Join(dir...)

	fmt.Printf("Working dir: %s\n", workingDir)

	foundFiles, err := tfvars.Parents(workingDir, "")
	if err != nil {
		return fmt.Errorf("Error getting var files %s", err.Error())
	}

	varfiles := make([]string, 0)
	for _, af := range foundFiles {
		if shouldInclude(environment, af) {
			varfiles = append(varfiles, af)
		}
	}

	vars, err := tfvars.Parse(varfiles)
	if err != nil {
		return fmt.Errorf("Error parsing varfiles %s", err.Error())
	}

	tfvars.ExportTfvars(vars)

	bucketExists := true
	if pullState {
		bucketExists, err = state.Configure(vars["aws_region"], environment, service, stack)
		if err != nil {
			return fmt.Errorf("Error configuring remote state %s", err.Error())
		}
	}

	wd, _ := os.Getwd()
	os.Chdir(workingDir)
	if err = beforeApply(); err != nil {
		return err
	}

	if pullState {
		exec.Command("terraform", "get").Run()
	}
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}

	// we just created the bucket
	if !bucketExists && pushState {
		if _, serr := state.Configure(vars["aws_region"], environment, service, stack); serr != nil {
			return fmt.Errorf("Error pushing remote state %s", serr.Error())
		}
	}

	if !destroy {
		if err = afterApply(); err != nil {
			return err
		}
	}

	os.Chdir(wd)
	return nil
}

func beforeApply() error {
	return runIfExists("before-apply.sh")
}

func afterApply() error {
	return runIfExists("after-apply.sh")
}

func runIfExists(name string) error {
	if _, err := os.Stat(name); err == nil {
		fmt.Printf("Executing %s\n", name)
		cmd := exec.Command("bash", name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
