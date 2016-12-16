package variables

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Action func(string) error
type Finder func(path string, re *regexp.Regexp, a Action) error

func Parents(path, re string) ([]string, error) {
	return findFiles(path, re, walkUpDirectories)
}

func Descendents(path, re string) ([]string, error) {
	return findFiles(path, re, walkDownDirectories)
}

func findFiles(path string, re string, find Finder) ([]string, error) {
	files := make([]string, 0)
	err := find(path, regexp.MustCompile(re), func(f string) error {
		files = append(files, f)
		return nil
	})
	return files, err
}

func walkDownDirectories(root string, re *regexp.Regexp, a Action) error {

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && re.MatchString(path) {
			if aerr := a(path); aerr != nil {
				return aerr
			}
		}
		return nil
	})
}

func walkUpDirectories(path string, re *regexp.Regexp, a Action) error {

	paths := strings.Split(path, string(os.PathSeparator))

	for len(paths) > 0 {
		p := strings.Join(paths, string(os.PathSeparator)) + string(os.PathSeparator)
		dir, err := ioutil.ReadDir(p)
		if err != nil {
			return err
		}

		for _, info := range dir {
			f := info.Name()
			if !info.IsDir() && re.MatchString(f) {
				if aerr := a(filepath.Join(p, f)); aerr != nil {
					return aerr
				}
			}
		}

		paths = paths[:len(paths)-1]
	}

	return nil
}
