package main

import (
	"testing"
)

type configTestCase struct {
	configString   string
	expectedResult ConfigT
}

func TestConfigReader(t *testing.T) {
	testCases := []configTestCase{
		{"localhost:25575:supersecret", ConfigT{"localhost:25575": "supersecret"}},
		{"localhost:25575:supersecret\n25576:secret", ConfigT{"localhost:25575": "supersecret", "25576": "secret"}},
		{"", nil},
		{"\n", nil},
		{"locahost\n", nil},
		{"localhost:555:humm\nlocal", nil},
	}

	for _, oneTestCase := range testCases {
		config, _ := ReadConfigFromString(oneTestCase.configString)
		identicalFlag := true
		if config == nil {
			if oneTestCase.expectedResult != nil {
				identicalFlag = false
			}
		} else {
			for key, value := range config {
				if value != oneTestCase.expectedResult[key] {
					identicalFlag = false
				}
			}
		}
		if !identicalFlag {
			t.Fail()
		}
	}
}
