package settings

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/qrclabs/nanogit/config"
	"github.com/qrclabs/nanogit/log"
)

var (
	AppPath  string
	ConfInfo config.ConfigInfo
)

func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal("settings: fail to get app path: %v", err)
	}
}
