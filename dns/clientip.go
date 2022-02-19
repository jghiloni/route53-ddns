package dns

import (
	"errors"
	"sort"

	resolver "github.com/Focinfi/go-dns-resolver"
)

const (
	MYIP_TARGET    = "myip.opendns.com"
	OPENDNS_SERVER = "resolver1.opendns.com:53"
)

type ClientIPResolver interface {
	ResolveClientIP() (string, error)
}

type dnsResolver struct{}

func NewClientIPResolver() ClientIPResolver {
	return &dnsResolver{}
}

func (d *dnsResolver) ResolveClientIP() (string, error) {
	resolver.Config.SetTimeout(uint(2))
	resolver.Config.RetryTimes = uint(4)

	results, err := resolver.Exchange(MYIP_TARGET, OPENDNS_SERVER, resolver.TypeA)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", errors.New("no results found")
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Priority < results[j].Priority
	})

	return results[0].Content, nil
}
