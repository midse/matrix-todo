package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	docopt "github.com/docopt/docopt-go"
	"golang.org/x/crypto/ssh/terminal"
)

func checkMultipleBoolArgs(arguments docopt.Opts, argNames []string) bool {
	checked := false

	for _, item := range argNames {
		res, err := arguments.Bool(item)

		if err != nil {
			continue
		}

		if res {
			checked = true
			break
		}

	}

	return checked
}

func getPassword() string {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to read input")
		os.Exit(1)
	}
	password := string(bytePassword)

	return strings.TrimSpace(password)
}
