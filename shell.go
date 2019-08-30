package powerd

import (
	"fmt"
	"io"
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

func newSecureShell(conf shellConfig) (*shell, error) {
	location := fmt.Sprintf("%s:%d", conf.Address, conf.Port)
	config := &ssh.ClientConfig{
		Timeout: time.Second * 10,
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

	return &shell{
		session: session,
		in:      stdin,
		out:     stdout,
		config:  conf,
		buffer:  make([]byte, 0),
		discard: make([]byte, 0),
	}, nil
}

func (s *shell) login() {

}

func (s *shell) navigate(dir string) {

}

func (s *shell) execute(command string) {

}

func (s *shell) start() error {
	return nil
}
