package dnsrelay

import (
	"strings"
)

const (
	CN_GROUP = "CN"
	FG_GROUP = "FG"
	REJECT_GROUP = "REJECT"

	MATCH_TYPE_DOMAIN_SUFFIX = "SUFFIX"
	MATCH_TYPE_DOMAIN_MATCH = "MATCH"
	MATCH_TYPE_DOMAIN_KEYWORD = "KEYWORD"
)

type DomainRules []DomainRule

func (self DomainRules) findGroup(domain string) string {
	for _, rule := range self {
		if rule.Match(domain) {
			return rule.Group
		}
	}
	return ""
}

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



