package prompt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sshtn/mods/ssh"

	"github.com/manifoldco/promptui"
)

const (
	START   = "ROOT"
	SSH     = "SSH"
	VPN     = "VPN"
	REVERSE = "REVERSE"
	PING    = "PING"

	executeIdx = -1
	parentIdx  = -2
	errorIdx   = -3
)

type promptFn func() (int, string, error)

var promptMap = map[string]promptFn{
	startPrompt: StartPrompt(),
}

type treePrompt struct {
	ui func() (int, string, error)

	current *treePrompt
	parent  *treePrompt
	child   []*treePrompt
}

func (tp *treePrompt) Run() error {
	idx, _, err := tp.ui()
	if err != nil {
		return err
	}

	switch idx {
	case executeIdx:
		err = tp.Run()
	case parentIdx:
		err = tp.parent.Run()
	default:
		err = tp.child[idx].Run()
	}

	if err != nil {
		return err
	}

	return nil
}

func StartPrompt() promptFn {
	prompt := promptui.Select{
		Label: "Network Tools",
		Items: []string{
			"SSH Connect",
			"VPN",
			"Reverse Proxy",
			"Ping",
		},
	}
	tree.child = []*treePrompt{
		promptSSHConnect(tree),
		promptVPN(tree),
		promptReverseProxy(tree),
		promptPing(tree),
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return 0, "", err
	}

	return idx, "", nil
}

func promptSSHConnect(parent *treePrompt) *treePrompt {
	tree := &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			configList, err := ssh.SSHConfigList()
			if err != nil {
				return 0, "", err
			}

			prompt := promptui.Select{
				Label: "SSH Connect",
				Items: configList,
			}

			idx, _, err := prompt.Run()
			if err != nil {
				return 0, "", err
			}
			return idx, "", nil
		},
	}

	tree.child = append(tree.child, promptSSHConfigAdd(tree))
	if len(tree.configList) > 1 {
		for _, config := range configList[1:] {
			tree.child = append(tree.child, promptSSHConfigLoad(tree, &config))
		}
	}

	return tree
}

func inputValidateFunc(s string) error {
	if s == "" {
		return errors.New("input is empty")
	} else {
		return nil
	}
}

func keyPathValidateFunc(keypath string) error {
	if keypath == "" {
		return errors.New("input is empty")
	}

	if _, err := os.Stat(keypath); os.IsNotExist(err) {
		return fmt.Errorf("'%s' is not exist", keypath)
	}
	return nil
}

func promptSSHConfigAdd(parent *treePrompt) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			prompts := []promptui.Prompt{
				{
					Label:    "username",
					Validate: inputValidateFunc,
				},
				{
					Label:    "keypath",
					Validate: keyPathValidateFunc,
				},
			}

			config := ssh.Config{}
			for i, prompt := range prompts {
				result, err := prompt.Run()
				if err != nil {
					return 0, "", err
				}
				switch i {
				case 0:
					config.Name = result
				case 1:
					config.KeyPath = result
				}
			}
			if err := ssh.SSHConfigCreate(config); err != nil {
				return errorIdx, "", err
			}
			return parentIdx, "", nil
		},
	}
}
func promptSSHConfigLoad(parent *treePrompt, config *ssh.Config) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			log.Printf("config: %+v\n", config)
			prompt := promptui.Select{
				Label: "SSH Tools",
				Items: []string{
					"Session Shell",
					"Tunneling",
				},
			}
			idx, _, err := prompt.Run()
			if err != nil {
				return 0, "", err
			}
			return idx, "", nil
		},
		child: []*treePrompt{},
	}
}

func SSHSession(parent *treePrompt) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			return 0, "", nil
		},
	}
}

func SSHTunneling(parent *treePrompt) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			return 0, "", nil
		},
	}
}

func promptVPN(parent *treePrompt) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			prompt := promptui.Select{
				Label: "VPN",
				Items: []string{
					"+ Add Config",
					"[twelve, /home/twelve/twelve_rsa]",
				},
			}
			idx, _, err := prompt.Run()
			if err != nil {
				return 0, "", err
			}

			return idx, "", nil
		},
	}
}
func promptReverseProxy(parent *treePrompt) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			prompt := promptui.Select{
				Label: "Reverse Proxy",
				Items: []string{
					"+ Add Config",
					"[twelve, /home/twelve/twelve_rsa]",
				},
			}
			idx, _, err := prompt.Run()
			if err != nil {
				return 0, "", err
			}

			return idx, "", nil
		},
	}
}

func promptPing(parent *treePrompt) *treePrompt {
	return &treePrompt{
		parent: parent,
		ui: func() (int, string, error) {
			prompt := promptui.Select{
				Label: "Ping",
				Items: []string{
					"+ Add Config",
					"[twelve, /home/twelve/twelve_rsa]",
				},
			}
			idx, _, err := prompt.Run()
			if err != nil {
				return 0, "", err
			}

			return idx, "", nil
		},
	}
}
