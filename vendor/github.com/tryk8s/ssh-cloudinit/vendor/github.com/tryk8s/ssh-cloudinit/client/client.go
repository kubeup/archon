package client

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"syscall"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func Run(conf *Config) error {
	auths := []ssh.AuthMethod{}
	if sshPublicKeyAuth := PublicKeyFile(conf.PublicKeyFile); sshPublicKeyAuth != nil {
		auths = append(auths, sshPublicKeyAuth)
	}
	if conf.Password != "" {
		auths = append(auths, ssh.Password(conf.Password))
	}
	passwordAuth := ssh.PasswordCallback(func() (secret string, err error) {
		fmt.Printf("%s@%s's password: ", conf.User, conf.Hostname)
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		if err != nil {
			return
		}
		return string(password), nil
	})
	auths = append(auths, passwordAuth)
	config := &ssh.ClientConfig{
		User: conf.User,
		Auth: auths,
	}
	server := fmt.Sprintf("%s:%d", conf.Hostname, conf.Port)
	err := executeCmds(server, config, conf.Stdout, conf.GetCmds())
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return nil
}
