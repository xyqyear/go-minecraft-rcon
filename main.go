package main

import (
	"fmt"
	"os"
)

func main() {
	address, password, command, err := ParseCommandLineArgs()
	printErrAndExit(err)

	client, err := NewClient(address)
	printErrAndExit(err)
	err = client.Login(password)
	printErrAndExit(err)
	response, err := client.SendCommandNaively(command)
	printErrAndExit(err)
	fmt.Println(response)
}

func printErrAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
