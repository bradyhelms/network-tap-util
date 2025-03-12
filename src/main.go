package main

import (
//  "bufio"
  "fmt"
  "log"
  "time"
  "github.com/bradyhelms/network-tap-util/src/utils"
)

func main() {
  // Get BERTHA credentials
  berthaIP := "172.13.201.12"
  fmt.Println("Connecting to BERTHA.")

  berthaClient, err := utils.GetSshClient(berthaIP)
  if err != nil {
    log.Fatalf("Error establishing ssh connection: %s", err)
  } else {
    fmt.Println("Conenction with BERTHA established.")
  }
  defer berthaClient.Close()
  
  fmt.Println("Starting the Vivado hardware server on BERTHA.")
  err = utils.RunCommand(berthaClient, "nohup /tools/Xilinx/Vivado/2024.1/bin/hw_server > /dev/null 2>&1 & exit")
  if err != nil {
    log.Printf("Error running command: %s", err)
  }

  // Wait a bit before trying to program the board
  time.Sleep(5 * time.Second)

  // Program the FPGA board
  almaIP := "192.168.50.54"
  fmt.Println("Connecting to Alma.")

  almaClient, err := utils.GetSshClient(almaIP)
  if err != nil {
    log.Fatalf("Error establishing ssh connection: %s", err)
  } else {
    fmt.Println("Conenction with Alma established.")
  }
  defer almaClient.Close()

  fmt.Println("Progamming the FPGA board.")
  err = utils.RunCommand(almaClient, "source /opt/PetaLinux/settings.sh" +
                               "&& cd /opt/bsp_build_20250310/xilinx-sp701-2024.1/images/linux/" + 
                               "&& /opt/PetaLinux/scripts/petalinux-boot jtag --prebuilt 3 " +
                               "--hw_server-url 172.13.201.12:3121")
  if err != nil {
    log.Printf("Error running command: %s", err)
  }
}
