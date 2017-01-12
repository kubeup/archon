package client

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"path/filepath"
)

func ensureUserData(conn *ssh.Client, stdout io.Writer, data string) error {
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	dir := filepath.Dir(UserDataPath)
	filename := filepath.Base(UserDataPath)
	session.Stdout = stdout

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintln(w, "C0644", len(data), filename)
		fmt.Fprint(w, data)
		fmt.Fprint(w, "\x00")
	}()

	err = session.Run(fmt.Sprintf("/usr/bin/scp -t %s", dir))
	if err != nil {
		return err
	}
	return nil
}

func executeCmds(conn *ssh.Client, stdout io.Writer, commands []string) (err error) {
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
