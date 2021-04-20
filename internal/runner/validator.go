package runner

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/internal/runner/storage"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// validate user-supplied option values before Runner.
func validate(opt *common.Options) error {
	var err error

	if opt.File != "" {
		opt.File, err = filepath.Abs(opt.File)
		if err != nil {
			return err
		}

		opt.List, err = readFile(opt.File)
		if err != nil {
			return err
		}

		if opt.Output != "" {
			opt.Output, err = filepath.Abs(opt.Output)
			if err != nil {
				return err
			}

			opt.Result, err = os.OpenFile(opt.Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			return nil
		}
	} else if opt.DB != "" && opt.DBURL != "" {
		var db storage.Database
		switch opt.DB {
		case "postgresql":
			c := &storage.Postgresql{
				DBURL:  opt.DBURL,
				DBUser: opt.DBUser,
				DBPass: opt.DBPass,
				DBName: opt.DBName,
			}
			err := c.Open()
			if err != nil {
				return err
			}
			db = c
		}
		l, err := db.Load()
		if err != nil {
			return err
		}
		cl := checkProxy(l)
		opt.List = cl
		return nil
	} else {
		return errors.New("no proxy provided")
	}
	return nil
}

// readFile which is returned as a unique slice proxies.
func readFile(path string) ([]string, error) {
	keys := make(map[string]bool)
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := scanner.Text()
		if _, value := keys[proxy]; !value {
			_, err := mubeng.Transport(proxy)
			if err == nil {
				keys[proxy] = true
				lines = append(lines, proxy)
			}
		}
	}

	if len(lines) < 1 {
		return lines, fmt.Errorf("open %s: has no valid proxy URLs", path)
	}

	return lines, scanner.Err()
}

func checkProxy(proxies []string) []string {
	keys := make(map[string]bool)
	var lines []string
	for _, proxy := range proxies {
		if _, value := keys[proxy]; !value {
			_, err := mubeng.Transport(proxy)
			if err != nil {
				keys[proxy] = true
				lines = append(lines, proxy)
			}
		}
	}
	return lines
}
