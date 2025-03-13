package main
//package capture

import (
	"fmt"
	"github.com/bradyhelms/network-tap-util/src/utils"
	"log"
	"time"
)

func main() {
    // Set up BERTHA for proxy
    berthaIP := "172.13.201.12"
    fmt.Println("Connecting to BERTHA.")

    berthaClient, err := utils.GetSshClient(berthaIP)
    if err != nil {
        log.Fatalf("Error establishing ssh connection: %s", err)
    } else {
        fmt.Println("Connection with BERTHA established.")
    }
    defer berthaClient.Close()

    // Connect to the FPGA board through the proxy
	fmt.Println("Establishing ssh connection the the tap.")
	tapIP := "192.168.1.96"
    tapClient, err := utils.GetSshClientWithProxy(berthaClient, tapIP)
    if err != nil {
        log.Fatalf("Error establishing proxied ssh connection: %s", err)
    } else {
        fmt.Println("Connection with the FPGA established.")
    }
	defer tapClient.Close()

	err = utils.RunCommand(tapClient, "nohup /home/petalinux/tap 20")
	if err != nil {
		log.Printf("Error running command: %s", err)
	}

	fmt.Println("Running packet capture for 20 seconds.")
	fmt.Println("Please wait")
	for i := 1; i < 13; i++ {
		time.Sleep(5 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\nDone!\n")

	fmt.Println("Starting .pcap file transfer.")
    // TEMP FIX THIS SHIT
	err = utils.RunCommand(tapClient, "scp -i /home/petalinux/.ssh/dropbear_key.db /home/root/capture.pcap petalinux@192.168.1.209:/home/petalinux/captures")
	if err != nil {
		log.Printf("Error receiving file: %s", err)
	}
}

