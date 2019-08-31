package smashd

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type shell struct {
	session *ssh.Session
	in      io.WriteCloser
	out     io.Reader
	config  shellConfig
	buffer  []byte
	discard []byte
}

var powerVerbs = map[string]string{
	"on":    "start",
	"off":   "stop",
	"reset": "reset",
}

func newSecureShell(conf shellConfig) (*shell, error) {
	location := fmt.Sprintf("%s:%d", conf.Address, 22)
	config := &ssh.ClientConfig{
		Timeout: time.Second * 15,
		User:    conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", location, config)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = session.Shell()
	if err != nil {
		return nil, err
	}

	return &shell{
		session: session,
		in:      stdin,
		out:     stdout,
		config:  conf,
		buffer:  make([]byte, 0),
		discard: make([]byte, 0),
	}, nil
}

func (s *shell) navigate(dir string) {
	s.in.Write([]byte("cd " + dir + "\n"))
}

func (s *shell) run(command string) {
	s.in.Write([]byte(command + "\n"))
}

func (s *shell) start() error {
	defer s.end()
	s.navigate("/system1/pwrmgtsvc1")
	action, success := powerVerbs[s.config.State]
	if !success {
		log.Fatalf("Power state %s not found\n! Exiting...", s.config.State)
	}
	s.run(action)
	return nil
}

func (s *shell) end() error {
	return s.session.Close()
}
