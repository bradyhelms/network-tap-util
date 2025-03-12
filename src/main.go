package main

import (
//  "bufio"
  "fmt"
  "golang.org/x/crypto/ssh"
  "golang.org/x/term"
  "log"
  "os"
  "time"
)

type SshCreds struct {
  User string
  Pass string
}

func main() {
  // Get BERTHA credentials
  berthaIP := "172.13.201.12"
  fmt.Println("Connecting to BERTHA.")

  berthaClient, err := getSshClient(berthaIP)
  if err != nil {
    log.Fatalf("Error establishing ssh connection: %s", err)
  } else {
    fmt.Println("Conenction with BERTHA established.")
  }
  defer berthaClient.Close()
  
  fmt.Println("Starting the Vivado hardware server on BERTHA.")
  err = runCommand(berthaClient, "/tools/Xilinx/Vivado/2024.1/bin/hw_server")
  if err != nil {
    log.Printf("Error running command: %s", err)
  }

  // Wait a bit before trying to program the board
  time.Sleep(5 * time.Second)

  // Program the FPGA board
  almaIP := "192.168.50.54"
  fmt.Println("Connecting to Alma.")

  almaClient, err := getSshClient(almaIP)
  if err != nil {
    log.Fatalf("Error establishing ssh connection: %s", err)
  } else {
    fmt.Println("Conenction with Alma established.")
  }
  defer almaClient.Close()

  fmt.Println("Progamming the FPGA board.")
  err = runCommand(almaClient, "source /opt/Petalinux/settings.sh" +
                               "&& cd /opt/bsp_build_20250310/xilinx-sp701-2024.1/images/linux/" + 
                               "&& /opt/PetaLinux/scripts/petalinux-boot jtag --prebuilt 3 " +
                               "--hw_server-url 172.13.201.12:3121")
  if err != nil {
    log.Printf("Error running command: %s", err)
  }
}

func getSshCredentials(hostname string) (creds SshCreds){
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

func getSshClient(ip string) (*ssh.Client, error) {
  creds := getSshCredentials(ip)

  clientConf := &ssh.ClientConfig{
    User: creds.User,
    Auth: []ssh.AuthMethod{
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

func runCommand(client *ssh.Client, cmd string) error {
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
