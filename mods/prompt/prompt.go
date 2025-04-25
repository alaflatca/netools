package prompt

import (
	"context"
	"errors"
	"sshtn/mods/reverse"
	"sshtn/mods/ssh"
	"sshtn/mods/vpn"

	"github.com/manifoldco/promptui"
)

const (
	START   = "START"
	SSH     = "SSH"
	VPN     = "VPN"
	REVERSE = "REVERSE"
	PING    = "PING"

	EXECUTE = "EXECUTE"
	PARENT  = "PARENT"
	CHILD   = "CHILD"
	ERROR   = "ERROR"

	promptStart   = "START"
	promptSSH     = "SSH Connect"
	promptVPN     = "VPN"
	promptReverse = "Reverse Proxy"
	promptPing    = "Ping"

	promptSSHConfigAdd  = "+ Add Config"
	promptSSHConfigList = "Config List"

	executeIdx = -1
	parentIdx  = -2
	errorIdx   = -3
)

func Run(ctx context.Context) error {
	prompt := promptui.Select{
		Label: "Network Tools",
		Items: []string{
			"SSH Connect",
			"VPN",
			"Reverse Proxy",
			"Ping",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return err
	}

	switch idx {
	case 0:
		return ssh.Run(ctx)
	case 1:
		return vpn.Run()
	case 2:
		return reverse.Run()
	default:
		return errors.New("invalid index")
	}
}
