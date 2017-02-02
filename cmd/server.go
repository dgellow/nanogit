package cmd

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/qrclabs/sshooks"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"

	"github.com/qrclabs/nanogit/auth"
	"github.com/qrclabs/nanogit/dir"
	"github.com/qrclabs/nanogit/log"
	"github.com/qrclabs/nanogit/settings"
)

var CmdServer = cli.Command{
	Name:   "server",
	Usage:  "Run the nanogit server",
	Action: runServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yml",
			Usage: "Custom configuration file path",
		},
		cli.IntFlag{
			Name:  "loglevel",
			Value: 3,
			Usage: "0=Trace, 1=Debug, 2=Info, 3=Warn, 4=Error, 5=Critical, 6=Fatal",
		},
	},
}

func pubKeyHandler(conn ssh.ConnMetadata, key ssh.PublicKey) (string, error) {
	log.Trace("server: pubKeyHandler")

	keystr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
	log.Trace("server: key: %s", keystr)

	_, err := settings.ConfInfo.LookupUserByKey(keystr)
	if err != nil {
		log.Error("server: unauthorized access: %v", err)
		return "", err
	}
	return keystr, nil
}

func runServer(c *cli.Context) error {
	log.Log.LogLevel = c.Int("loglevel")
	log.Trace("server: runServer")

	log.Trace("server: read config file")
	settings.ConfInfo.ConfigFile = c.String("config")
	settings.ConfInfo.ReadFile()

	log.Trace("server: ConfigFile: %s", settings.ConfInfo.ConfigFile)

	commandsHandlers := map[string]func(string, string, string) (io.WriteCloser, io.ReadCloser, io.ReadCloser, error){
		"git-upload-pack":    handleUploadPack,
		"git-upload-archive": handleUploadArchive,
		"git-receive-pack":   handleReceivePack,
	}

	sshooksConfig := &sshooks.ServerConfig{
		Host:              "localhost",
		Port:              1337,
		PrivatekeyPath:    "key.rsa",
		KeygenConfig:      sshooks.SSHKeygenConfig{"rsa", ""},
		PublicKeyCallback: pubKeyHandler,
		CommandsCallbacks: commandsHandlers,
		Log:               log.Log,
	}
	sshooks.Listen(sshooksConfig)

	// Keep the program running
	for {
	}
}

func handleUploadPack(keyId string, cmd string, args string) (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	log.Trace("server: Handle git-upload-pack: args: %s", args)
	org, repo := dir.SplitPath(dir.CleanPath(args))
	pathExists, err := dir.IsPathExist(org, repo)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error when checking if repository path exists: %v", err)
	}
	if !pathExists {
		return nil, nil, nil, fmt.Errorf("Repository path doesn't exist, org: %s, repo: %s", org, repo)
	}
	read, write := auth.CheckAuth(keyId, args)
	log.Trace("server: Rights policy: read: %t, write: %t", read, write)
	if !read {
		return nil, nil, nil, fmt.Errorf("Unauthorized read access: %s", args)
	}
	if !write {
		return nil, nil, nil, fmt.Errorf("Unauthorized write access: %s", args)
	}

	command := exec.Command("git-upload-pack", args)
	log.Debug("server: handleUploadPack: command: %v", command)
	stdout, err := command.StdoutPipe()
	log.Trace("server: handleUploadPack: stdout: %v", stdout)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error when trying to access stdout: %v", err)
	}
	stderr, err := command.StderrPipe()
	log.Trace("server: handleUploadPack: stderr: %v", stderr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error when trying to access stderr: %v", err)
	}
	input, err := command.StdinPipe()
	log.Trace("server: handleUploadPack: input: %v", input)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error when trying to access input: %v", err)
	}

	return input, stdout, stderr, nil
}

func handleUploadArchive(keyId string, cmd string, args string) (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	log.Trace("server: Handle git-upload-archive: args: %s", args)
	read, write := auth.CheckAuth(keyId, args)
	log.Trace("server: Rights policy: read: %t, write: %t", read, write)
	if !read {
		return nil, nil, nil, fmt.Errorf("Unauthorized read access: %s", args)
	}
	if !write {
		return nil, nil, nil, fmt.Errorf("Unauthorized write access: %s", args)
	}
	return nil, nil, nil, nil
}

func handleReceivePack(keyId string, cmd string, args string) (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	log.Trace("server: Handle git-receive-pack: args: %s", args)
	read, write := auth.CheckAuth(keyId, args)
	log.Trace("server: Rights policy: read: %t, write: %t", read, write)
	if !read {
		return nil, nil, nil, fmt.Errorf("Unauthorized read access: %s", args)
	}
	if !write {
		return nil, nil, nil, fmt.Errorf("Unauthorized write access: %s", args)
	}
	return nil, nil, nil, nil
}
