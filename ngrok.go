package ngrokcd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Ngrok is ngrok client
type Ngrok struct {
	host string
	port int
}

// NewNgrok returns Ngrok instance
func NewNgrok(host string, port int) *Ngrok {
	return &Ngrok{
		host,
		port,
	}
}

// Ping checks if the server is responding
func (ng Ngrok) Ping() error {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/tunnels", ng.host, ng.port))
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("unknown error")
}

// Tunnels lists all the exported tunnels
func (ng Ngrok) Tunnels() (Tunnels, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/tunnels", ng.host, ng.port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Tunnels Tunnels `json:"tunnels"`
	}{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Tunnels, nil
}
