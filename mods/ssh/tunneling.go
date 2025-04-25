package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh"
)

func tunneling(ctx context.Context, config *Config, client *ssh.Client) error {
	if err := promptTunneling(config); err != nil {
		return err
	}

	lsnr, err := net.Listen("tcp", "127.0.0.1:"+config.LocalPort)
	if err != nil {
		return err
	}
	defer lsnr.Close()

	select {
	case <-ctx.Done():
		return nil
	default:

		// spinner()

		localConn, err := lsnr.Accept()
		if err != nil {
			return err
		}
		defer localConn.Close()

		remoteConn, err := client.Dial("tcp", "127.0.0.1:"+config.RemotePort)
		if err != nil {
			return err
		}
		defer remoteConn.Close()

		go io.Copy(remoteConn, localConn)
		go io.Copy(localConn, remoteConn)

		fmt.Printf("Local(:%s) <-------> Remote(:%s)\r", config.LocalPort, config.RemotePort)
	}
	return nil
}

func promptTunneling(config *Config) error {
	prompts := []promptui.Prompt{
		{
			Label: "Local Port",
		},
		{
			Label: "Remote Port",
		},
	}

	for i, prompt := range prompts {
		port, err := prompt.Run()
		if err != nil {
			return err
		}

		switch i {
		case 0:
			config.LocalPort = port
		case 1:
			config.RemotePort = port
		}
	}

	return nil
}

// func localListener(config *Config) error {

// }

// func RemoteToLocal() {

// }

func spinner() {
	spinnerText := []string{"|", "/", "-", "\\"}
	for i := 0; i < 30; i++ {
		fmt.Printf("%s\r", spinnerText[i%len(spinnerText)])
		time.Sleep(100 * time.Millisecond)
	}
}
