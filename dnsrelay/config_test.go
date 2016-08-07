package dnsrelay

import (
	"testing"
	"fmt"
)

func TestNewConfig(t *testing.T) {

	config, err := NewConfig("../dnsrelay.toml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(config)

}