package ssh

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func makeKnownHostsCallback(knownHostsPath string) (ssh.HostKeyCallback, error) {
	knownHostsCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(knownHostsPath)
			if err != nil {
				return nil, err
			}
			f.Close()
			knownHostsCallback, err = knownhosts.New(knownHostsPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := knownHostsCallback(hostname, remote, key)
		if err == nil {
			return nil
		}
		if _, ok := err.(*knownhosts.KeyError); ok {
			f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			line := knownhosts.Line([]string{hostname}, key)
			_, err = f.WriteString(line + "\n")
			if err != nil {
				return err
			}
			fmt.Printf("[INFO] 새로운 서버 키를 known_hosts에 등록했습니다: %s\n", hostname)
		}
		return nil
	}, nil
}

func CreateSshConfig(userName, keyFile string) (*ssh.ClientConfig, error) {
	knownHostsPath := sshConfigPath("known_hosts")
	hostKeyCallback, err := makeKnownHostsCallback(knownHostsPath)
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User: userName,

		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(singer),
		},
		HostKeyCallback: hostKeyCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,       // RSA
			ssh.KeyAlgoED25519,   // ED25519
			ssh.KeyAlgoECDSA256,  // ECDSA (NIST P-256)
			ssh.KeyAlgoRSASHA256, // RSA with SHA-256
			ssh.KeyAlgoRSASHA512,
		},
		Timeout: 5 * time.Second,
	}, nil
}

func sshConfigPath(filename string) string {
	return filepath.Join(os.Getenv("HOME"), ".ssh", filename)
}
