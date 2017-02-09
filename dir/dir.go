package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgellow/nanogit/log"
	"github.com/dgellow/nanogit/settings"
)

func CleanPath(path string) string {
	return strings.Replace(path, "'", "", -1)
}

func SplitPath(path string) (org string, repo string, err error) {
	sliceStr := strings.Split(path, "/")
	if len(sliceStr) < 2 {
		return "", "", fmt.Errorf("A path should be: orgname/reponame, got: %s", path)
	}
	return strings.ToLower(sliceStr[0]), strings.ToLower(sliceStr[1]), nil
}

func getDataRoot() (string, error) {
	confDataRoot := settings.ConfInfo.Conf.Server.DataRoot
	if confDataRoot == "" {
		return "", fmt.Errorf("Data root in configuration file is empty")
	}

	log.Debug("dir: AppPath: %s", settings.AppPath)
	log.Debug("dir: Server.DataRoot: %s", settings.ConfInfo.Conf.Server.DataRoot)

	if confDataRoot[0] == '/' {
		return settings.ConfInfo.Conf.Server.DataRoot, nil
	} else {
		return filepath.Join(settings.AppPath, settings.ConfInfo.Conf.Server.DataRoot), nil
	}
}

func GetOrgDir(org string) (string, error) {
	dataRoot, err := getDataRoot()
	if err != nil {
		return dataRoot, err
	}
	return filepath.Join(dataRoot, org), nil
}

func GetRepoDir(org string, repo string) (string, error) {
	dataRoot, err := getDataRoot()
	if err != nil {
		return dataRoot, err
	}
	return filepath.Join(dataRoot, org, repo), nil
}

func IsOrgExist(path string) (bool, error) {
	log.Trace("dir: IsOrgExist, path: %s", path)
	target, err := GetOrgDir(path)
	if err != nil {
		return false, err
	}

	fi, err := os.Stat(target)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return true, nil
	} else {
		return false, fmt.Errorf("A file exists with the org name, it should be a directory: %s", target)
	}
}

func IsRepoExist(orgPath string, repoPath string) (bool, error) {
	log.Trace("dir: IsRepoExist, orgPath: %s, repoPath: %s", orgPath, repoPath)
	target, err := GetRepoDir(orgPath, repoPath)
	log.Trace("dir: IsRepoExist, target: %s", target)
	if err != nil {
		return false, err
	}
	fi, err := os.Stat(target)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return true, nil
	} else {
		return false, fmt.Errorf("A file exists with the repo name, it should be a directory: %s", target)
	}
}

func IsPathExist(org string, repo string) (bool, error) {
	log.Trace("dir: IsPathExist")
	orgExists, err := IsOrgExist(org)
	if err != nil {
		return orgExists, err
	}
	repoExists, err := IsRepoExist(org, repo)
	if err != nil {
		return repoExists, err
	}
	return true, nil
}
