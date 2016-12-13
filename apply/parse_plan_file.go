package main

import (
	"errors"
	"strings"
)

type planfile struct {
	destroy     bool
	environment string
	service     string
	stack       string
}

func parsePlanFile(filename string) (*planfile, error) {
	pf := new(planfile)

	filename = strings.Replace(filename, ".plan", "", 1)

	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return pf, errors.New("invalid planfile name")
	}

	if parts[0] == "destroy" {
		pf.destroy = true
		parts = parts[1:len(parts)]
	}

	if err := validateEnvironment(parts[0]); err != nil {
		return pf, err
	}

	pf.environment = parts[0]
	parts = parts[1:len(parts)]

	if len(parts) >= 2 {
		pf.stack = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}

	pf.service = strings.Join(parts, "_")

	return pf, nil
}
