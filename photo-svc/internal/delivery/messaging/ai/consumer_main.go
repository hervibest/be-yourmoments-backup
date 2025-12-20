package aiconsumer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func (s *AIConsumer) ConsumeAllEvents(ctx context.Context) error {
	for _, subject := range s.subjects {
		if err := s.setupConsumer(subject); err != nil && !strings.Contains(err.Error(), "consumer name already") {
			fmt.Println("error " + err.Error())
			return fmt.Errorf("failed to setup consumer for %s: %w", subject, err)
		}

		sub, err := s.js.PullSubscribe(
			subject,
			s.durableNames[subject],
			nats.BindStream("AI_SIMILAR_STREAM"),
		)
		if err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
		}

		go s.startConsumer(ctx, sub, subject)
	}

	return nil
}

func (s *AIConsumer) startConsumer(ctx context.Context, sub *nats.Subscription, subject string) {
	s.logs.Log(fmt.Sprintf("started consumer for subject : %s", subject))

	for {
		select {
		case <-ctx.Done():
			s.logs.Log(fmt.Sprintf("stopping consumer for subject : %s", subject))

			return
		default:
			msgs, err := sub.Fetch(10, nats.MaxWait(2*time.Second))
			if err != nil {
				if err == nats.ErrTimeout {
					time.Sleep(200 * time.Millisecond) // 🔥 WAJIB
					continue
				}
				s.logs.Error(fmt.Sprintf("failed to fetch messages with error %v", err))
				time.Sleep(time.Second)
				continue
			}
			for _, msg := range msgs {
				s.handleMessage(ctx, msg)
			}
		}
	}
}
