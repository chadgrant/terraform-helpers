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

			fmt.Printf("Replacing %s with %s in %s\n", find, repl, path)

			read, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			re := regexp.MustCompile(find)
			if err != nil {
				return err
			}

			nc := re.ReplaceAll(read, []byte(repl))

			err = ioutil.WriteFile(path, nc, 0)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
