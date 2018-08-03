package main

import (
	"log"
	"time"

	ngrokcd "github.com/alyyousuf7/ngrok-custom-domain"
)

func getNgrok(hostname string, port int) *ngrokcd.Ngrok {
	ngrok := ngrokcd.NewNgrok(hostname, port)

	log.Print("waiting for ngrok")
	for {
		if tunnels, err := ngrok.Tunnels(); err == nil && len(tunnels) > 0 {
			break
		}
	}
	log.Print("ngrok is up")

	return ngrok
}

func performMatchingWithRetries(ngrok *ngrokcd.Ngrok, config *configuration, retries, retryDelay int) error {
	var err error
	for i := 0; i < retries; i++ {
		if i > 0 {
			log.Printf("retry #%d", i+1)
		}

		err = performMatching(ngrok, config)
		if err == nil {
			return nil
		}

		if i+1 < retries {
			time.Sleep(time.Duration(retryDelay) * time.Millisecond)
		}
	}

	return err
}

func performMatching(ngrok *ngrokcd.Ngrok, config *configuration) error {
	matcher := ngrokcd.NewMatcher(ngrok)

	tunnelMap, err := config.Parse()
	if err != nil {
		return err
	}

	for tunnel, records := range tunnelMap {
		for _, record := range records {
			if err := matcher.Add(tunnel, record.CNAME, record.Service); err != nil {
				return err
			}
		}
	}

	if err := matcher.Match(); err != nil {
		return err
	}

	return nil
}
