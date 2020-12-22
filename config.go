package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ConfigError is the error type returned by ReadConfig
type ConfigError uint8

func (c ConfigError) Error() string {
	switch c {
	case ConfigErrorFileNotFount:
		return "file not found"
	case ConfigErrorOSNotSupported:
		return "your system is not supported (home directory not found)"
	case ConfigErrorPermissionDenied:
		return "permission denied when reading password file"
	case ConfigErrorUnknown:
		return "unknown error happened"
	case ConfigErrorNoConfig:
		return "no config file found"
	case ConfigErrorParseFailed:
		return "failed to parse config file"
	}
	return "unknown error happened"
}

const (
	// ConfigErrorFileNotFount is a file not found error
	ConfigErrorFileNotFount ConfigError = iota
	// ConfigErrorOSNotSupported happens when os.UserHomeDir() cannot be used
	ConfigErrorOSNotSupported
	// ConfigErrorPermissionDenied happens when not having enough permission while reading config file
	ConfigErrorPermissionDenied
	// ConfigErrorUnknown happens when a unknown error occurs
	ConfigErrorUnknown
	// ConfigErrorNoConfig happens when a config file exists but does not really has any config.
	ConfigErrorNoConfig
	// ConfigErrorParseFailed happens when a config file exists but failed to parse.
	ConfigErrorParseFailed
)

// ConfigT is map[string]string
type ConfigT map[string]string

// ReadConfig needs config file lies in ~/.go-mcrcon/pass
func ReadConfig() (ConfigT, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, ConfigErrorOSNotSupported
	}
	filePath := homeDir + "/.go-mcrcon/pass"
	return ReadConfigFromFile(filePath)
}

// ReadConfigFromFile reads config file with it's path
func ReadConfigFromFile(filePath string) (ConfigT, error) {
	fileContentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ConfigErrorFileNotFount
		} else if os.IsPermission(err) {
			return nil, ConfigErrorPermissionDenied
		} else {
			return nil, ConfigErrorUnknown
		}
	}

	return ReadConfigFromString(string(fileContentBytes))
}

// ReadConfigFromString is separated from ReadConfigFromFile because of test
func ReadConfigFromString(configString string) (ConfigT, error) {
	fileLines := strings.Split(configString, "\n")
	config := make(ConfigT)
	if len(fileLines) > 0 {
		for _, fileLine := range fileLines {
			fields := strings.Split(fileLine, ":")
			switch len(fields) {
			case 2:
				config[fields[0]] = fields[1]
			case 3:
				config[fmt.Sprintf("%s:%s", fields[0], fields[1])] = fields[2]
			default:
				return nil, ConfigErrorParseFailed
			}
		}
	} else {
		return nil, ConfigErrorNoConfig
	}

	if len(config) > 0 {
		return config, nil
	}
	return nil, ConfigErrorNoConfig
}
