package powerd

import (
	"log"

	"github.com/urfave/cli"
)

type shellConfig struct {
	Username string
	Password string
	Address  string
	Port     uint
	State    string
}

func newConfig(c *cli.Context) shellConfig {
	username := c.String("username")
	password := c.String("password")
	address := c.String("Address")
	port := c.Uint("port")
	state := c.String("state")

	if len(username) == 0 {
		log.Fatalf("username not provided")
	}
	if len(password) == 0 {
		log.Fatalf("password not provided")
	}
	if len(address) == 0 {
		log.Fatalf("address not provided")
	}
	if port == 0 {
		log.Fatalf("invalid port")
	}
	if len(state) == 0 {
		log.Fatalf("state not provided")
	}

	return shellConfig{
		Username: username,
		Password: password,
		Address:  address,
		Port:     port,
		State:    state,
	}
}
