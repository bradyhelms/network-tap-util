package capture

import (
	"bufio"
	"fmt"
	"github.com/bradyhelms/network-tap-util/src/utils"
	"log"
	"os"
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

	if yesNo("Start packet capture?") {
		err = utils.RunCommand(tapClient, "/home/petalinux/scripts/tap")
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
