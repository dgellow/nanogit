package cmd

import (
	"strings"

	"golang.org/x/crypto/ssh"
	"github.com/urfave/cli"
	"github.com/qrclabs/sshgit"

	"github.com/qrclabs/nanogit/log"
	"github.com/qrclabs/nanogit/settings"
)

var CmdServer = cli.Command{
	Name: "server",
	Usage: "Run the nanogit server",
	Action: runServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
			Value: "config.yml",
			Usage: "Custom configuration file path",
		},
	},
}

func pubKeyHandler(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Trace("server: pubKeyHandler")

	keystr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
	log.Trace("server: key: %s", keystr)

	_, err := settings.ConfInfo.LookupUserByKey(keystr);
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

	sshConfig := sshgit.ServerConfig{
		Host: "localhost",
		Port: 1337,
		PrivatekeyPath: "key.rsa",
		KeygenConfig: sshgit.SSHKeygenConfig{"rsa", ""},
		PublicKeyCallback: pubKeyHandler,
	}
	sshgit.Listen(sshConfig)

	// Keep the program running
	for {}
}
