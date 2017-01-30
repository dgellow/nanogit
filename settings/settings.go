package settings

import (
	"fmt"
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
		log.Fatal(4, "settings: fail to get app path: %v", err)
	}
}

func NewConsoleLogger(logLevel int) {
	log.NewLogger(0, "console", fmt.Sprintf(`{"level": %d}`, logLevel))
}
