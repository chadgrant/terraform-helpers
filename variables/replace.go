package variables

import (
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
			//nc := strings.Replace(string(read), find, repl, -1)

			nc := re.ReplaceAll(read, []byte(repl))

			err = ioutil.WriteFile(path, nc, 0)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
