package ssh

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh"
)

func Run(ctx context.Context) error {
	selectConfig, err := promptConfigSelect()
	if err != nil {
		return err
	}

	if err = promptConfigCreate(selectConfig); err != nil {
		return err
	}

	idx, err := promptSSHFuncs()
	if err != nil {
		return err
	}

	config, err := createSshConfig(selectConfig.Name, selectConfig.KeyPath)
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", selectConfig.Address(), config)
	if err != nil {
		return err
	}
	defer client.Close()

	switch idx {
	case 0:
		return session(ctx, selectConfig, client)
	case 1:
		return tunneling(ctx, selectConfig, client)
	}

	return nil
}

func promptSSHFuncs() (int, error) {
	prompt := promptui.Select{
		Label: "SSH Tools",
		Items: []string{
			"Session",
			"Tunneling",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return idx, nil
}

func validateSshAddr(s string) error {
	if s == "" {
		return errors.New("is empty")
	}
	if ip := net.ParseIP(s); ip == nil {
		return errors.New("IP address is not valid")
	}
	return nil
}
func validateSshPort(s string) error {
	if s == "" {
		return errors.New("is empty")
	}

	port, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	if port > 65536 {
		return errors.New("port is too big ( < 65536 )")
	}
	return nil
}

func validateSshKeyPath(s string) error {
	if s == "" {
		return errors.New("is empty")
	}
	if _, err := os.Stat(s); os.IsNotExist(err) {
		return err
	}
	return nil
}
func validateSshName(s string) error {
	if s == "" {
		return errors.New("is empty")
	}
	if len(s) > 10 {
		return errors.New("name length is too long ( < 10 )")
	}
	return nil
}

func promptConfigCreate(config *Config) error {
	if config.isExist {
		return nil
	}

	prompts := []promptui.Prompt{
		{
			Label:    "Address",
			Validate: validateSshAddr,
		},
		{
			Label:    "Port",
			Validate: validateSshPort,
		},
		{
			Label:    "Name",
			Validate: validateSshName,
		},
		{
			Label:    "KeyPath",
			Validate: validateSshKeyPath,
		},
	}
	for i, prompt := range prompts {
		result, err := prompt.Run()
		if err != nil {
			return err
		}

		switch i {
		case 0:
			config.Addr = result
		case 1:
			config.Port = result
		case 2:
			config.Name = result
		case 3:
			config.KeyPath = result
		}
	}
	config.isExist = true

	return SSHConfigCreate(config)
}

func promptConfigSelect() (*Config, error) {
	configList, err := SSHConfigList()
	if err != nil {
		return nil, err
	}

	items := []string{}
	for _, conf := range configList {
		items = append(items, conf.String())
	}

	prompt := promptui.Select{
		Label: "SSH Connect",
		Items: items,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	if idx > 0 {
		return &configList[idx], nil
	} else {
		return &configList[0], nil
	}
}
