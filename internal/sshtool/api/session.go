package sshtool

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

type SessionArgs struct {
	Network   string
	Host      string
	Port      int
	ClientCfg *ssh.ClientConfig
}

func NewSSHClient(args SessionArgs) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", args.Host, args.Port)
	client, err := ssh.Dial(args.Network, addr, args.ClientCfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Session(ctx context.Context, client *ssh.Client, in io.Reader, out io.Writer) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		session.Close()
	}()
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("linux", 200, 60, modes); err != nil {
		return err
	}

	session.Stdin = in
	session.Stdout = out
	session.Stderr = out

	if err := session.Shell(); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}
