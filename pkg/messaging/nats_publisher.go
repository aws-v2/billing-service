package messaging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	maxPublishRetries = 3
)

type NATSPublisher struct {
	nc      *nats.Conn
	profile string
}

func NewNATSPublisher(url, user, password, profile string) (*NATSPublisher, error) {
	opts := []nats.Option{
		nats.Name("Billing-Service"),
		nats.Timeout(5 * time.Second),
	}

	if user != "" && password != "" {
		opts = append(opts, nats.UserInfo(user, password))
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &NATSPublisher{
		nc:      nc,
		profile: profile,
	}, nil
}

func (p *NATSPublisher) BuildSubject(domain, action string) string {
	return fmt.Sprintf("%s.billing.v1.%s.%s", p.profile, domain, action)
}

func (p *NATSPublisher) Publish(subject string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var lastErr error
	for attempt := 1; attempt <= maxPublishRetries; attempt++ {
		if err := p.nc.Publish(subject, data); err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
			continue
		}

		if err := p.nc.Flush(); err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to publish after %d attempts: %w", maxPublishRetries, lastErr)
}

func (p *NATSPublisher) Request(subject string, payload interface{}, timeout time.Duration) (*nats.Msg, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return p.nc.Request(subject, data, timeout)
}

func (p *NATSPublisher) Close() {
	if p.nc != nil {
		p.nc.Close()
	}
}
