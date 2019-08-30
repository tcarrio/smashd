package powerd

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

func run(c *cli.Context) {
	if c.Bool("version") {
		fmt.Println(version())
		return
	}

	conf := newConfig(c)
	sh, err := newSecureShell(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = sh.start()
	if err != nil {
		log.Fatal(err)
	}
}
