package client

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"path"
	"strings"
)

func runCmd(conn *ssh.Client, conf *Config, cmd string) error {
	session, err := newSession(conn, conf)
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(cmd)
}

func ensureCloudData(conn *ssh.Client, conf *Config) error {
	err := runCmd(conn, conf, fmt.Sprintf("sudo sh -c 'rm -rf %s; mkdir -p %s && chown %s %s'", CloudDataPath, CloudDataPath, conf.User, CloudDataPath))
	if err != nil {
		return err
	}

	err = runCmd(conn, conf, fmt.Sprintf("echo 'dsmode: local' > %s", path.Join(CloudDataPath, "meta-data")))
	if err != nil {
		return err
	}

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = conf.Stdout
	fmt.Fprint(conf.Stdout, "> Uploading data\n")

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintln(w, "C0644", len(conf.UserData), "user-data")
		fmt.Fprint(w, conf.UserData)
		fmt.Fprint(w, "\x00")
	}()

	err = session.Run(fmt.Sprintf("scp -t %s", CloudDataPath))
	return err
}

func newSession(conn *ssh.Client, conf *Config) (session *ssh.Session, err error) {
	session, err = conn.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 25, modes)
	if err != nil {
		return nil, fmt.Errorf("Unable to request pty: %v", err)
	}

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return
	}

	stdout := &SudoPipe{
		Stdout:    conf.Stdout,
		StdinPipe: stdinPipe,
		Password:  conf.Password,
	}
	session.Stdout = stdout

	return session, nil
}

func executeCmds(conn *ssh.Client, conf *Config) (err error) {
	commands := conf.GetCmds()
	for _, command := range commands {
		session, err := newSession(conn, conf)
		if err != nil {
			return err
		}
		defer session.Close()

		ignoreError := false
		if strings.HasPrefix(command, "! ") {
			ignoreError = true
			command = command[2:]
		}
		fmt.Printf(">>> %s\n", command)
		fmt.Fprintf(session.Stdout, ">>> %s\n", command)

		err = session.Run(command)
		if err != nil && !ignoreError {
			return err
		}
	}

	return
}

type SudoPipe struct {
	Stdout    io.Writer
	StdinPipe io.Writer
	Password  string

	line []byte
}

func (sp *SudoPipe) Write(p []byte) (count int, err error) {
	if sp.Stdout != nil {
		count, err = sp.Stdout.Write(p)
	}

	if sp.Password == "" || sp.StdinPipe == nil {
		return
	}

	for _, b := range p {
		sp.line = append(sp.line, b)
		line := string(sp.line)
		if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
			_, err = sp.StdinPipe.Write([]byte(sp.Password + "\n"))
			if err != nil {
				break
			}
		}

		if b == '\n' {
			sp.line = []byte{}
		}
	}

	return
}
