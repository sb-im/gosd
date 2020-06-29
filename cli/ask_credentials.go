package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func askCredentials() (string, string) {
	fd := int(os.Stdin.Fd())

	if !terminal.IsTerminal(fd) {
		fmt.Fprintf(os.Stderr, "This is not a terminal, exiting.\n")
		os.Exit(1)
	}

	fmt.Print("Enter Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")

	state, _ := terminal.GetState(fd)
	defer terminal.Restore(fd, state)
	bytePassword, _ := terminal.ReadPassword(fd)

	fmt.Printf("\n")
	return strings.TrimSpace(username), strings.TrimSpace(string(bytePassword))
}
