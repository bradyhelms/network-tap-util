package capture

import (
	"bufio"
	"fmt"
	"github.com/bradyhelms/network-tap-util/src/utils"
	"log"
	"os"
	"time"
)

// First connect to board over ssh
func main() {
	tapIP := "192.168.0.0"
	fmt.Println("Establishing ssh connection the the tap.")

	tapClient, err := utils.GetSshClient(tapIP)
	if err != nil {
		log.Fatalf("Error establishing ssh connection: %s", err)
	} else {
		fmt.Println("Connection with tap established.")
	}
	defer tapClient.Close()

	if !yesNo("Start packet capture?") {
		os.Exit(0)
	}

	err = utils.RunCommand(tapClient, "/home/petalinux/scripts/tap 60")
	if err != nil {
		log.Printf("Error running command: %s", err)
	}

	fmt.Println("Running packet capture for 60 seconds.")
	fmt.Println("Please wait")
	for i := 1; i < 13; i++ {
		time.Sleep(5 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\nDone!\n")

	fmt.Println("Starting .pcap file transfer.")
	err = utils.RunCommand(tapClient, "scp /home/petalinux/scripts/capture.pcap petalinux@192.168.1.209:/home/petalinux/captures")
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
