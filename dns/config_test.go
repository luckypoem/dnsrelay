package dns

import (
	"testing"
	"fmt"
)

func TestNewConfig(t *testing.T) {

	config, err := NewConfig("../dns.toml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(config)
	fmt.Println(config.DomainRules)

}