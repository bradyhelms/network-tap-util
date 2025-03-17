package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"log"
	"os"
    "path/filepath"
)

type SshCreds struct {
	User string
	Pass string
}

func GetSshCredentials(hostname string) (creds SshCreds) {
	fmt.Printf("Enter ssh credentials for %s.\n", hostname)

	fmt.Printf("Username: ")
	fmt.Scanln(&creds.User)

    fmt.Printf("NOTE: You can leave this field blank if using key based auth.")
	fmt.Printf("Password: ")
	tempPass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to read password: %s", err)
	}

	creds.Pass = string(tempPass)
	fmt.Println()
	return creds
}

func GetSshClient(ip string) (*ssh.Client, error) {
	creds := GetSshCredentials(ip)

	clientConf := &ssh.ClientConfig{
		User: creds.User,
		Auth: []ssh.AuthMethod{
			PublicKeyAuth(),
			ssh.Password(creds.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", ip+":22", clientConf)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %w", err)
	}

	return client, nil
}

func RunCommand(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("Failed to run command: %w.\n\nOutput:\n%s", err, output)
	}
	fmt.Printf("Output of command '%s':\n%s\n", cmd, output)
	return nil
}

func PublicKeyAuth() ssh.AuthMethod {
    sshDir := os.ExpandEnv("$HOME/.ssh")

    files, err := os.ReadDir(sshDir)
    if err != nil {
        log.Fatalf("Failed to read SSH directory: %s", err)
    }

    var signers []ssh.Signer

    for _, file := range files {
        if file.IsDir() || filepath.Ext(file.Name()) == ".pub" {
            continue
        }

        keyPath := filepath.Join(sshDir, file.Name())
        key, err := os.ReadFile(keyPath)
        if err != nil {
            log.Printf("Skipping key %s: failed to read (%s)", keyPath, err)
            continue
        }

        signer, err := ssh.ParsePrivateKey(key)
        if err != nil {
            log.Printf("Skipping key %s: failed to read (%s)", keyPath, err)
            continue
        }

        signers = append(signers, signer)
    }

    if len(signers) == 0 {
        log.Fatal("No valid SSH private keys found")
    }

    return ssh.PublicKeys(signers...)
}

func GetSshClientWithProxy(proxyClient *ssh.Client, targetIP string) (*ssh.Client, error) {
    conn, err := proxyClient.Dial("tcp", targetIP+":22")
    if err != nil {
        return nil, fmt.Errorf("Failed to create tunnel to %s: %w", targetIP, err)
    }

    creds := GetSshCredentials(targetIP)

    clientConf := &ssh.ClientConfig{
        User: creds.User,
        Auth: []ssh.AuthMethod{
            ssh.Password(creds.Pass),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    
    clientConn, chans, reqs, err := ssh.NewClientConn(conn, targetIP+":22", clientConf)
    if err != nil {
        return nil, fmt.Errorf("Failed to establish SSH client connection: %w", err)
    }

    return ssh.NewClient(clientConn, chans, reqs), nil
}
