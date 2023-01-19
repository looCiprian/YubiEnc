package input_mng

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func InputRequest(msg string, isSensitive bool) (string, error) {
	fmt.Print(msg)

	var inputPin string
	if isSensitive {
		inputPinbyte, _ := term.ReadPassword(int(syscall.Stdin))
		inputPin = string(inputPinbyte)
		fmt.Println("")
		if inputPin == "" {
			return "", errors.New("[-] Cannot read user input")
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		inputPin, _ = reader.ReadString('\n')
		inputPin = inputPin[:len(inputPin)-1]
	}
	return inputPin, nil
}
