package ssh

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sshtn/mods/storage"
	"strings"
	"time"

	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
)

const SshCode = "1"

type sshReader struct {
	list []Config
}
type Config struct {
	Name    string
	KeyPath string
}

func (cfg *Config) String() string {
	var line string
	if cfg.Name == "" || cfg.KeyPath == "" {
		line = "===="
	} else {
		line = fmt.Sprintf("%s,%s,%s", SshCode, cfg.Name, cfg.KeyPath)
	}
	return line
}

func (sr *sshReader) Read(reader io.Reader) error {
	sr.list = make([]Config, 0)
	sr.list = append(sr.list, Config{Name: "+ Add Config"})

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 { // ssh
			continue
		}

		split := strings.Split(text, ",")
		if split[0] != SshCode {
			continue
		}

		cfg := Config{
			Name:    split[1],
			KeyPath: split[2],
		}
		sr.list = append(sr.list, cfg)
	}

	if scanner.Err() != nil {
		return fmt.Errorf("scanner error: %s", scanner.Err())
	}

	return nil
}

func SSHConfigList() ([]Config, error) {
	sr := &sshReader{}

	err := storage.Read(sr)
	if err != nil {
		return nil, err
	}

	return sr.list, nil
}

type sshWriter struct {
	Config Config
}

func (sr *sshWriter) Write(writer io.Writer) error {
	line := sr.Config.String()

	line = line + "\n"
	_, err := writer.Write([]byte(line))
	if err != nil {
		return err
	}

	return nil
}

func SSHConfigCreate(config Config) error {
	sw := &sshWriter{
		Config: config,
	}

	if err := storage.Write(sw); err != nil {
		return err
	}

	return nil
}

func createSshConfig(userName, keyFile string) *ssh.ClientConfig {
	knownHostsCallback, err := knownhosts.New(sshConfigPath("known_hosts"))
	if err != nil {
		log.Fatal(err)
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}

	return &ssh.ClientConfig{
		User: userName,

		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(singer),
		},
		HostKeyCallback: ssh.HostKeyCallback(knownHostsCallback),
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,       // RSA
			ssh.KeyAlgoED25519,   // ED25519
			ssh.KeyAlgoECDSA256,  // ECDSA (NIST P-256)
			ssh.KeyAlgoRSASHA256, // RSA with SHA-256
			ssh.KeyAlgoRSASHA512,
		},
		Timeout: 5 * time.Second,
	}
}

func sshConfigPath(filename string) string {
	return filepath.Join(os.Getenv("HOME"), ".ssh", filename)
}
