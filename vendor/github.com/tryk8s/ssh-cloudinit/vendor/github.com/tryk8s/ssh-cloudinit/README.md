# ssh-cloudinit

ssh-cloudinit is a provisioning tool that helps you initialize cloud servers through cloudinit over ssh.

## Supported Systems

* Ubuntu
* Coreos (TBD)

## Install

  `go get github.com/tryk8s/ssh-cloudinit`
  
## Usage

```bash
ssh-cloudinit [options] <server>
  
  -key string
        SSH key path (default "~/.ssh/id_rsa")
  -os string
        Server OS (default "ubuntu")
  -port int
        Server SSH port (default 22)
  -remote string
        Remote cloud-init url
  -user string
        Server SSH user (default "root")

```

`-remote` is required.

## License
MIT
