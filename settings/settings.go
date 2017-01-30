package settings

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/qrclabs/nanogit/log"
	"github.com/qrclabs/nanogit/config"
)

var (
	AppPath string
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
	log.NewLogger(0, "console", `{"level": 0}`)

	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal(4, "settings: fail to get app path: %v", err)
	}
}
