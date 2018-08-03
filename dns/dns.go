package dns

import (
	"fmt"
)

var (
	// ErrUnauthorized is returned when 401 is received from API
	ErrUnauthorized = fmt.Errorf("unauthorized")
	// ErrRecordNotFound is returned when a particular record was not found
	ErrRecordNotFound = fmt.Errorf("record not found")
)

const (
	// DefaultTTL is the TTL used when creating/updating CNAME record
	DefaultTTL = 600
)

type DNS interface {
	AddRecord(name, data string, ttl int) error
	FindRecord(name string) (string, error)
	UpdateRecord(name, data string, ttl int) error
	UpsertRecord(name, data string, ttl int) error
}

type Record struct {
	CNAME   string
	Service DNS
}
