package settings

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dgellow/nanogit/config"
	"github.com/dgellow/nanogit/log"
)

var (
	AppPath  string
	ExecPath string
	ConfInfo config.ConfigInfo
)

func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func appPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func init() {
	var err error
	if AppPath, err = appPath(); err != nil {
		log.Fatal("settings: fail to get app path: %v", err)
	}
	if ExecPath, err = execPath(); err != nil {
		log.Fatal("settings: fail to get exec path: %v", err)
	}
}
