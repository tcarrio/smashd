package powerd

import (
	"os"

	"github.com/urfave/cli"
)

// New runs the powerd CLI, parsing user input and spawning a new shell
// to connect to the SuperMicro board and manage power state
func New() {
	app := cli.NewApp()
	app.Version = version()
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "version, v",
			Usage: "print the version",
		},
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "the username for authenticating against the server",
			EnvVar: "POWERD_USERNAME",
			Value:  "ADMIN",
		},
		cli.StringFlag{
			Name:   "pass, p",
			Usage:  "the password for authenticating against the server",
			EnvVar: "POWERD_PASSWORD",
		},
		cli.StringFlag{
			Name:   "address, a",
			Usage:  "the address of the server",
			EnvVar: "POWERD_ADDRESS",
		},
		cli.UintFlag{
			Name:   "port, P",
			Usage:  "the SSH port of the server",
			EnvVar: "POWERD_PORT",
			Value:  22,
		},
		cli.StringFlag{
			Name:  "state, s",
			Usage: "the desired power state",
			Value: "on",
		},
	}
	app.Action = run
	app.Run(os.Args)
}
