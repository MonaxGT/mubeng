package common

import (
	"os"
	"time"
)

// Options consists of the configuration required.
type Options struct {
	File    string
	Address string
	DB      string
	DBURL   string
	DBUser  string
	DBPass  string
	DBName  string
	Check   bool
	Timeout time.Duration
	Rotate  int
	Verbose bool
	Output  string
	Result  *os.File
	List    []string
	Daemon  bool
}
