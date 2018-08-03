package ngrokcd

import (
	"fmt"

	"github.com/alyyousuf7/ngrok-custom-domain/dns"
)

type tunnelRecord struct {
	cname  string
	tunnel *Tunnel
}

type Matcher struct {
	ngrok   *Ngrok
	matches map[dns.DNS][]tunnelRecord
}

func NewMatcher(ngrok *Ngrok) *Matcher {
	return &Matcher{
		ngrok,
		map[dns.DNS][]tunnelRecord{},
	}
}

func (m *Matcher) Add(tunnelName, cname string, record dns.DNS) error {
	tunnels, err := m.ngrok.Tunnels()
	if err != nil {
		return err
	}

	tunnel := tunnels.FindName(tunnelName)

	if tunnel == nil {
		return fmt.Errorf("could not find tunnel %s", tunnelName)
	}

	m.matches[record] = append(m.matches[record], tunnelRecord{cname, tunnel})
	return nil
}

func (m Matcher) Match() error {
	for domain, records := range m.matches {
		for _, record := range records {
			if err := domain.UpsertRecord(record.cname, record.tunnel.Hostname(), dns.DefaultTTL); err != nil {
				return err
			}
		}
	}

	return nil
}
