package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	errorcode "github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model/event"
	"github.com/nats-io/nats.go"
)

func (s *PhotoConsumer) setupConsumer(subject string) error {
	consumerConfig := &nats.ConsumerConfig{
		Durable:       s.durableNames[subject],
		AckPolicy:     nats.AckExplicitPolicy,
		MaxDeliver:    5,
		BackOff:       []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second},
		DeliverPolicy: nats.DeliverAllPolicy,
		AckWait:       30 * time.Second,
		FilterSubject: subject,
	}

	_, err := s.js.AddConsumer("PHOTO_STREAM", consumerConfig)
	return err
}

func (s *PhotoConsumer) handleMessage(ctx context.Context, msg *nats.Msg) {

	var err error
	switch msg.Subject {
	case "photo.persist.facecam":
		event := new(event.PersistFacecamEvent)
		if err := sonic.ConfigStd.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		s.logs.Log(fmt.Sprintf("unmarshalled event: %+v", event))

		err = s.userUseCase.UpdateHasFacecam(ctx, event.UserID)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	default:
		err = fmt.Errorf("unknown subject: %s", msg.Subject)
	}

	if err := msg.Ack(); err != nil {
		s.logs.Error(fmt.Sprintf("failed to ACK message : %s", err))
	}
}

func (s *PhotoConsumer) handleError(msg *nats.Msg, err error) {

	appErr, ok := err.(*helper.AppError)
	if !ok {
		appErr = &helper.AppError{Code: errorcode.ErrInternal}
	}

	switch appErr.Code {
	case errorcode.ErrInvalidArgument:
		s.logs.Error(fmt.Sprintf("invalid to ACK message : %s", err))
	default:
		delay := 10 * time.Second
		if err := msg.NakWithDelay(delay); err != nil {
		}
	}
}
