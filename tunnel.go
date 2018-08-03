package ngrokcd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Tunnel is ngrok tunnel struct
type Tunnel struct {
	Name string
	URL  url.URL
}

// Hostname returns hostname from URL
func (t Tunnel) Hostname() string {
	return t.URL.Hostname()
}

// UnmarshalJSON is Unmarshaler implementation
func (t *Tunnel) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	if err := json.Unmarshal(j, &rawStrings); err != nil {
		fmt.Println("Hello", err)
		return err
	}

	for k, v := range rawStrings {
		switch strings.ToLower(k) {
		case "name":
			t.Name = v.(string)
		case "public_url":
			u, err := url.Parse(v.(string))
			if err != nil {
				return err
			}
			t.URL = *u
		}
	}

	return nil
}

// Tunnels is []Tunnel
type Tunnels []Tunnel

// First return the first tunnel in list
func (l Tunnels) First() *Tunnel {
	if len(l) == 0 {
		return nil
	}

	return &l[0]
}

// FindName finds one tunnel by name
func (l Tunnels) FindName(name string) *Tunnel {
	for _, v := range l {
		if v.Name == name {
			return &v
		}
	}

	return nil
}
