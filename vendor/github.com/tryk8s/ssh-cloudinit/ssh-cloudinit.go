package main

import (
	"flag"
	"fmt"
	"github.com/tryk8s/ssh-cloudinit/client"
	"os"
)

var (
	remote        string
	osType        string
	port          int
	user          string
	publicKeyFile string
)

func init() {
	flag.StringVar(&remote, "remote", "", "Remote cloud-init url")
	flag.StringVar(&osType, "os", "ubuntu", "Server OS")
	flag.IntVar(&port, "port", 22, "Server SSH port")
	flag.StringVar(&user, "user", "root", "Server SSH user")
	flag.StringVar(&publicKeyFile, "key", os.Getenv("HOME")+"/.ssh/id_rsa", "SSH key path")
	flag.Usage = func() {
		fmt.Printf("Usage: ssh-cloudinit [options] <server>\n\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}
	conf := &client.Config{
		Hostname:      args[0],
		User:          user,
		PublicKeyFile: publicKeyFile,
		Server:        remote,
		Port:          port,
		Os:            osType,
		Stdout:        os.Stdout,
	}
	client.Run(conf)
}
