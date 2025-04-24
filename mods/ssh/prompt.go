package ssh

import "github.com/manifoldco/promptui"

func Run() error {

	sshMain := func() error {
		configList, err := SSHConfigList()
		if err != nil {
			return err
		}
		prompt := promptui.Select{
			Label: "SSH Connect",
		}

	}

	return nil
}
