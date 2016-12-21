package variables

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func Replace(dir, ext, find, repl string) error {
	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		m, err := filepath.Match(ext, f.Name())
		if err != nil {
			return err
		}

		if m {

			read, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			re := regexp.MustCompile(find)
			if err != nil {
				return err
			}

			nc := re.ReplaceAll(read, []byte(repl))

			if !areEqual(read, nc) {
				fmt.Printf("Replacing %s with %s in %s\n", find, repl, path)

				err = ioutil.WriteFile(path, nc, 0)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func areEqual(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
