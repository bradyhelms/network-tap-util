package main
//package capture

import (
	"bufio"
	"fmt"
	"github.com/bradyhelms/network-tap-util/src/utils"
	"log"
	"os"
//	"time"
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


    /*
	if !yesNo("Start packet capture?") {
		os.Exit(0)
	}
    */

	err = utils.RunCommand(tapClient, "nohup /home/petalinux/tap 5")
	if err != nil {
		log.Printf("Error running command: %s", err)
	}

    /*
	fmt.Println("Running packet capture for 60 seconds.")
	fmt.Println("Please wait")
	for i := 1; i < 13; i++ {
		time.Sleep(5 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\nDone!\n")
    */

	fmt.Println("Starting .pcap file transfer.")
    // TEMP FIX THIS SHIT
	err = utils.RunCommand(tapClient, "scp -i /home/petalinux/.ssh/dropbear_key.db /home/root/capture.pcap petalinux@192.168.1.209:/home/petalinux/captures")
	if err != nil {
		log.Printf("Error receiving file: %s", err)
	}
}

func yesNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(question + " [y/n]: ")
		input, _ := reader.ReadString('\n')

		if input == "y" {
			return true
		} else if input == "n" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
			return yesNo(question)
		}
	}
}
