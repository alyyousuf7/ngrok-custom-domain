package main

import (
	"fmt"
	"io/ioutil"

	"github.com/alyyousuf7/ngrok-custom-domain/dns"
	"gopkg.in/yaml.v2"
)

type tunnelConf map[string]struct {
	Records []string
}

type recordConf map[string]struct {
	DNS   string
	CNAME string
}

type dnsConf map[string]struct {
	Service string
	Domain  string
	Key     string
	Secret  string
}

type configuration struct {
	Tunnels tunnelConf
	Records recordConf
	DNS     dnsConf
}

func newConfiguration(tunnelName, cname, domain, service, key, secret string) *configuration {
	return &configuration{
		Tunnels: tunnelConf{
			tunnelName: {
				Records: []string{"record1"},
			},
		},
		Records: recordConf{
			"record1": {
				DNS:   "dns1",
				CNAME: cname,
			},
		},
		DNS: dnsConf{
			"dns1": {
				Domain:  domain,
				Service: service,
				Key:     key,
				Secret:  secret,
			},
		},
	}
}

func loadConfiguration(path string) (*configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &configuration{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (config configuration) validate() error {
	// Check if the required records exists
	for _, tunnel := range config.Tunnels {
		for _, record := range tunnel.Records {
			matched := false
			for rec := range config.Records {
				if record == rec {
					matched = true
					break
				}
			}

			if !matched {
				return fmt.Errorf("%s record not found", record)
			}
		}
	}

	// Check if the required dns exists
	for _, record := range config.Records {
		matched := false
		for dns := range config.DNS {
			if record.DNS == dns {
				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf("%s dns not found", record.DNS)
		}
	}

	return nil
}

// Parse returns a map of tunnels to records
func (config configuration) Parse() (map[string][]dns.Record, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	var (
		dnsMap    = make(map[string]dns.DNS)
		recordMap = make(map[string]dns.Record)
		tunnelMap = make(map[string][]dns.Record)
	)

	for k, v := range config.DNS {
		var service dns.DNS
		switch v.Service {
		case "godaddy":
			service = dns.NewGodaddy(v.Domain, v.Key, v.Secret)
		default:
			return nil, fmt.Errorf("unknown dns %s", v.Service)
		}

		dnsMap[k] = service
	}

	for k, v := range config.Records {
		service, ok := dnsMap[v.DNS]
		if !ok {
			return nil, fmt.Errorf("could not find service %s", v.DNS)
		}
		recordMap[k] = dns.Record{
			CNAME:   v.CNAME,
			Service: service,
		}
	}

	for k, tunnel := range config.Tunnels {
		for _, record := range tunnel.Records {
			rec, ok := recordMap[record]
			if !ok {
				return nil, fmt.Errorf("could not find record %s", record)
			}

			tunnelMap[k] = append(tunnelMap[k], rec)
		}
	}

	return tunnelMap, nil
}
