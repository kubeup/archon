package client

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
)

func executeCmds(server string, conf *ssh.ClientConfig, stdout io.Writer, commands []string) (err error) {
	conn, err := ssh.Dial("tcp", server, conf)
	if err != nil {
		return
	}
	for _, command := range commands {
		session, err := conn.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()

		session.Stdout = stdout

		fmt.Fprintf(stdout, ">>> %s\n", command)
		err = session.Run(command)
		if err != nil {
			return err
		}
	}

	return
}
