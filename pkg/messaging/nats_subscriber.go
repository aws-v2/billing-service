package messaging

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NATSSubscriber struct {
	nc     *nats.Conn
	prefix string
}

func NewNATSSubscriber(url, user, password, prefix string) (*NATSSubscriber, error) {
	opts := []nats.Option{
		nats.Name("Billing-Service-Subscriber"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * nats.DefaultReconnectWait),
	}

	if user != "" && password != "" {
		opts = append(opts, nats.UserInfo(user, password))
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS for subscriber: %w", err)
	}

	return &NATSSubscriber{
		nc:     nc,
		prefix: prefix,
	}, nil
}

// Subscribe listens for metrics on dev.v1.billing.metric.*
// It uses a Queue Group ("billing_ingestion") to ensure only one instance receives a particular message.
func (s *NATSSubscriber) Subscribe(handler func(subject string, data []byte) error) (*nats.Subscription, error) {
	subject := fmt.Sprintf("%s.billing.metric.*", s.prefix)
	queueGroup := "billing_ingestion"

	sub, err := s.nc.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		if err := handler(msg.Subject, msg.Data); err != nil {
			log.Printf("Error handling metric from subject %s: %v", msg.Subject, err)
			return
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	log.Printf("Subscribed to %s in queue group %s", subject, queueGroup)
	return sub, nil
}

func (s *NATSSubscriber) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}
