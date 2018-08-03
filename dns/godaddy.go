package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var _ DNS = &GoDaddy{}

type GoDaddy struct {
	domain, key, secret string
}

func NewGodaddy(domain, key, secret string) *GoDaddy {
	return &GoDaddy{
		domain,
		key,
		secret,
	}
}

func (g GoDaddy) AddRecord(name, data string, ttl int) error {
	body, err := json.Marshal([]struct {
		Data       string `json:"data"`
		Name       string `json:"name"`
		TTL        int    `json:"ttl"`
		RecordType string `json:"type"`
	}{{data, name, ttl, "CNAME"}})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records", g.domain),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", g.key, g.secret))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	if res.StatusCode == 401 {
		return ErrUnauthorized
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := struct {
		Code    string
		Message string
	}{}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	return fmt.Errorf(response.Message)
}

func (g GoDaddy) FindRecord(name string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/CNAME/%s", g.domain, name),
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", g.key, g.secret))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 200 {
		response := []struct {
			Data string
		}{}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", err
		}

		if len(response) == 0 {
			return "", ErrRecordNotFound
		}

		return response[0].Data, nil
	}

	if res.StatusCode == 401 {
		return "", ErrUnauthorized
	}

	response := struct {
		Code    string
		Message string
	}{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return "", fmt.Errorf("%d: %s", res.StatusCode, response.Message)
}

func (g GoDaddy) UpdateRecord(name, data string, ttl int) error {
	body, err := json.Marshal([]struct {
		Data       string `json:"data"`
		Name       string `json:"name"`
		TTL        int    `json:"ttl"`
		RecordType string `json:"type"`
	}{{data, name, ttl, "CNAME"}})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/CNAME/%s", g.domain, name),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", g.key, g.secret))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	if res.StatusCode == 401 {
		return ErrUnauthorized
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := struct {
		Code    string
		Message string
	}{}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	return fmt.Errorf("%d: %s", res.StatusCode, response.Message)
}

func (g GoDaddy) UpsertRecord(name, data string, ttl int) error {
	oldData, err := g.FindRecord(name)
	if err != nil && err != ErrRecordNotFound {
		return err
	}

	if err == ErrRecordNotFound {
		log.Printf("adding %s: %s", name, data)
		return g.AddRecord(name, data, ttl)
	}

	if oldData == data {
		return nil
	}

	log.Printf("updating %s: %s", name, data)
	return g.UpdateRecord(name, data, ttl)
}
