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

const SshCode = "ssh"

type sshReader struct {
	list []Config
}
type Config struct {
	Name    string
	Addr    string
	Port    string
	KeyPath string

	LocalPort  string
	RemotePort string

	isExist bool
}

func (conf *Config) Address() string {
	return conf.Addr + ":" + conf.Port
}

func (conf *Config) String() string {
	if conf.isExist {
		return fmt.Sprintf("[ %s:%s ] %s %s", conf.Addr, conf.Port, conf.Name, conf.KeyPath)
	} else {
		return conf.Name
	}
}

func (conf *Config) Line() string {
	var line string
	if conf.Name == "" || conf.KeyPath == "" {
		line = "===="
	} else {
		line = fmt.Sprintf("%s,%s,%s,%s,%s", SshCode, conf.Addr, conf.Port, conf.Name, conf.KeyPath)
	}
	return line
}

func (sr *sshReader) Read(reader io.Reader) error {
	sr.list = make([]Config, 0)
	sr.list = append(sr.list, Config{
		Name:    "+ Config Create",
		KeyPath: "",
		isExist: false,
	})

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
			Addr:    split[1],
			Port:    split[2],
			Name:    split[3],
			KeyPath: split[4],
			isExist: true,
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
	Config *Config
}

func (sr *sshWriter) Write(writer io.Writer) error {
	if !sr.Config.isExist {
		return nil
	}
	line := sr.Config.Line()

	line = line + "\n"
	_, err := writer.Write([]byte(line))
	if err != nil {
		return err
	}

	return nil
}

func SSHConfigCreate(config *Config) error {
	sw := &sshWriter{
		Config: config,
	}

	if err := storage.Write(sw); err != nil {
		return err
	}

	return nil
}

func createSshConfig(userName, keyFile string) *ssh.ClientConfig {
	var hostKeyCallback ssh.HostKeyCallback
	knownHostsCallback, err := knownhosts.New(sshConfigPath("known_hosts"))
	if err != nil {
		log.Println(err)
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		hostKeyCallback = ssh.HostKeyCallback(knownHostsCallback)
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
		HostKeyCallback: hostKeyCallback,
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
