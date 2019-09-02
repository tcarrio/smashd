package smashd

import (
	"bufio"
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
	code    chan int
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

	fmt.Printf("Started session with %s\n", conf.Address)

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
		code:    make(chan int),
	}, nil
}

func (s *shell) startLogger() {
	scanner := bufio.NewScanner(s.out)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	s.code <- 0
}

func (s *shell) navigate(dir string) error {
	return s.run("cd " + dir)
}

func (s *shell) show() error {
	return s.run("show")
}

func (s *shell) run(command string) error {
	_command := command + "\n"
	_, error := s.in.Write([]byte(_command))
	return error
}

func (s *shell) start() error {
	defer s.end()
	go s.startLogger()

	err := s.show()
	if err != nil {
		s.code <- 1
		return fmt.Errorf("Failed at show")
	}

	time.Sleep(time.Second * 2)

	err = s.navigate("/system1/pwrmgtsvc1")
	if err != nil {
		s.code <- 1
		return fmt.Errorf("Failed at navigate")
	}
	action, success := powerVerbs[s.config.State]
	if !success {
		log.Fatalf("Power state %s not found\n! Exiting...", s.config.State)
	}

	time.Sleep(time.Second * 2)

	err = s.run(action)
	if err != nil {
		s.code <- 1
		return fmt.Errorf("Failed at run [%s]", action)
	}
	return nil
}

func (s *shell) end() error {
	err := s.session.Close()
	if err != nil {
		return err
	}

	select {
	case code := <-s.code:
		if code != 0 {
			return fmt.Errorf("Exit code: %d", code)
		}
	case <-time.After(2 * time.Second):
		fmt.Println("No exit code issued after 2 seconds")
	}

	return nil
}
