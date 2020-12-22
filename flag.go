package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// ParseStringArgs parse the argument from a string, used for testing
func ParseStringArgs(args string) (string, string, string, error) {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	return parseArgs(func() {
		pflag.CommandLine.Parse(strings.Split(args, " "))
	})
}

// ParseCommandLineArgs parse arguments from the command line.
// the return strings are server address, password and command
// the priority of "address" argument is higher than hostname and port
// if address, hostname and port are all not privided, then the server address defaults to localhost:25575
// if only the port is provided, then defaults to "localhost:{port}"
func ParseCommandLineArgs() (string, string, string, error) {
	return parseArgs(func() {
		pflag.Parse()
	})
}

func parseArgs(parseFunc func()) (string, string, string, error) {
	address := pflag.StringP("address", "a", "", "the address of the server (hostname and port)")
	hostname := pflag.StringP("hostname", "n", "", "the hostname of the server. Will be override by address flag")
	port := pflag.IntP("port", "p", 25575, "the port of the server. Will be override by address")
	password := pflag.StringP("password", "s", "", "rcon password. Overrides the one in the config file")
	configPath := pflag.StringP("config-file", "f", "", "the path for config file, defaults to ~/.go-mcrcon/pass")
	configString := pflag.StringP("config", "c", "", "the content of config, will override config file.")

	parseFunc()

	// fmt.Println("address:", *address)
	// fmt.Println("hostname:", *hostname)
	// fmt.Println("port:", *port)
	// fmt.Println("password:", *password)
	// fmt.Println("configPath:", *configPath)
	// fmt.Println("configString:", *configString)
	// fmt.Println("tail:", pflag.Args())

	var config ConfigT
	if *configString != "" {
		config, _ = ReadConfigFromString(*configString)
	} else if *configPath != "" {
		config, _ = ReadConfigFromFile(*configPath)
	} else {
		config, _ = ReadConfig()
	}

	if *address != "" {
		columnCount := strings.Count(*address, ":")
		switch columnCount {
		case 0:
			*address = fmt.Sprintf("%s:25575", *address)
		case 1:
		default:
			return "", "", "", errors.New("invalid address.")
		}
	} else {
		if *hostname == "" {
			*address = fmt.Sprintf("localhost:%d", *port)
		} else {
			*address = fmt.Sprintf("%s:%d", *hostname, *port)
		}
	}

	if *password == "" {
		_port := strings.Split(*address, ":")[1]
		if passwordFromConfig, ok := config[*address]; ok {
			*password = passwordFromConfig
		} else if passwordFromConfig, ok := config[_port]; ok {
			*password = passwordFromConfig
		} else {
			return "", "", "", errors.New("no password provided.")
		}
	}

	if pflag.NArg() == 0 {
		fmt.Println("wrong usage.")
		pflag.PrintDefaults()
		return "", "", "", errors.New("wrong usage.")
	}
	command := strings.Join(pflag.Args(), " ")

	return *address, *password, command, nil
}
