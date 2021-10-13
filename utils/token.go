package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func sshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func GetToken() (string, error) {
	sshConfig := &ssh.ClientConfig{
		User: "lagoon",
		Auth: []ssh.AuthMethod{
			sshAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", "ssh.lagoon.amazeeio.cloud:32222", sshConfig)
	if err != nil {
		return "", fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return "", err
	}

	out, err := session.CombinedOutput("token")

	if err != nil {
		return "", errors.New("Could not run Token command on remote Lagoon")
	}

	connection.Close()

	return strings.TrimSpace(string(out)), nil
}
