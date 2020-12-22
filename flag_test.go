package main

import (
	"fmt"
	"testing"
)

type flagTestCase struct {
	commandLineArgs string
	expectedResult  []string
}

func TestFlag(t *testing.T) {
	testCases := []flagTestCase{
		{"-s supersecret list", []string{"localhost:25575", "supersecret", "list"}},
		{"-s supersecret -p 25575 whitelist list", []string{"localhost:25575", "supersecret", "whitelist list"}},
		{"-s supersecret -a example.com:25578 list", []string{"example.com:25578", "supersecret", "list"}},
		{"-c 25577:supersecret -p 25577 list", []string{"localhost:25577", "supersecret", "list"}},
		{"-c 25577:supersecret -p 25577 -s secret list", []string{"localhost:25577", "secret", "list"}},
		{"-c 25577:supersecret -p 25578 list", []string{"", "", ""}},
		{"-a localhost:4555:44 -s supersecret list", []string{"", "", ""}},
		{"-c localhost:12345:supersecret\npi:34567:secret -a pi:34567 list", []string{"pi:34567", "secret", "list"}},
		{"-c localhost:12345:supersecret\npi:34567:secret -p 34567 list", []string{"", "", ""}},
		{"-c localhost:12345:supersecret\npi:34567:secret -n pi list", []string{"", "", ""}},
		{"-c localhost:12345:supersecret\npi:34567:secret\n25575:super list", []string{"localhost:25575", "super", "list"}},
	}

	for _, oneTestCase := range testCases {
		if address, password, command, err := ParseStringArgs(oneTestCase.commandLineArgs); address != oneTestCase.expectedResult[0] || password != oneTestCase.expectedResult[1] || command != oneTestCase.expectedResult[2] {
			fmt.Println(oneTestCase, "failed")
			fmt.Println("got", address, password, command, err)
			t.Fail()
		}
	}
}
