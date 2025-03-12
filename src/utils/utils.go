package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"log"
	"os"
)

type SshCreds struct {
	User string
	Pass string
}

func GetSshCredentials(hostname string) (creds SshCreds) {
	fmt.Printf("Enter ssh credentials for %s.\n", hostname)

	fmt.Printf("Username: ")
	fmt.Scanln(&creds.User)

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
//			publicKeyAuth(),
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
	key, err := os.ReadFile("/path/to/id_rsa")
	if err != nil {
		log.Fatal("Failed to read private key: ", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	return ssh.PublicKeys(signer)
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
