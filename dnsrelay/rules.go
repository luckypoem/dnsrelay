package dnsrelay

import (
	"strings"
)

const (
	MATCH_TYPE_DOMAIN_SUFFIX = "suffix"
	MATCH_TYPE_DOMAIN_MATCH = "match"
	MATCH_TYPE_DOMAIN_KEYWORD = "keywork"
)

type DomainRule struct {
	MatchType string `toml:"match-type"`
	Group     string `toml:"domain-group"`
	Values    []string `toml:"value"`
}

func (rule DomainRule) Match(input string) bool {

	for _, value := range rule.Values {
		switch rule.MatchType {
		case MATCH_TYPE_DOMAIN_MATCH:
			if input == value {
				return true
			}
		case MATCH_TYPE_DOMAIN_SUFFIX:
			if strings.HasSuffix(input, value) {
				return true
			}
		case MATCH_TYPE_DOMAIN_KEYWORD:
			if strings.Contains(input, value) {
				return true
			}
		default:
			continue
		}
	}

	return false
}



