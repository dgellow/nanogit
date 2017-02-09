package main

import (
	"os"
	"sort"

	"github.com/urfave/cli"

	"github.com/dgellow/nanogit/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "nanogit"
	app.Usage = "simple git server"
	app.Action = func(c *cli.Context) error {
		return cli.NewExitError("nanogit: no argument given: Run `nanogit --help` for more information", 1)
	}

	app.Commands = []cli.Command{
		cmd.CmdServer,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}
