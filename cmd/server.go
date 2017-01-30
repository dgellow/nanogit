package cmd

import (
	"strings"

	"github.com/qrclabs/sshooks"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"

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
	},
}

func pubKeyHandler(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Trace("server: pubKeyHandler")

	keystr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
	log.Trace("server: key: %s", keystr)

	_, err := settings.ConfInfo.LookupUserByKey(keystr)
	if err != nil {
		log.Error(3, "server: unauthorized access: %v", err)
		return nil, err
	}
	return &ssh.Permissions{}, nil
}

func runServer(c *cli.Context) error {
	log.Trace("server: runServer")

	settings.ConfInfo.ConfigFile = c.String("config")
	settings.ConfInfo.ReadFile()

	log.Trace("server: ConfigFile: %s", settings.ConfInfo.ConfigFile)

	commandsHandlers := map[string]func(string) error{
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
	}
	sshooks.Listen(sshooksConfig)

	// Keep the program running
	for {
	}
}

func handleUploadPack(args string) error {
	log.Trace("Handle git-upload-pack: args: %s", args)
	return nil
}

func handleUploadArchive(args string) error {
	log.Trace("Handle git-upload-archive: args: %s", args)
	return nil
}

func handleReceivePack(args string) error {
	log.Trace("Handle git-receive-pack: args: %s", args)
	return nil
}
