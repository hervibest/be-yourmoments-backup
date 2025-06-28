package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model/event"
	"github.com/nats-io/nats.go"
)

func (s *AIConsumer) setupConsumer(subject string) error {
	consumerConfig := &nats.ConsumerConfig{
		Durable:       s.durableNames[subject],
		AckPolicy:     nats.AckExplicitPolicy,
		MaxDeliver:    5,
		BackOff:       []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second},
		DeliverPolicy: nats.DeliverAllPolicy,
		AckWait:       30 * time.Second,
		FilterSubject: subject,
	}

	_, err := s.js.AddConsumer("AI_SIMILAR_STREAM", consumerConfig)
	return err
}

func (s *AIConsumer) handleMessage(ctx context.Context, msg *nats.Msg) {

	var err error
	switch msg.Subject {
	case "ai.bulk.photo":
		s.logs.Log(fmt.Sprintf("received message on subject: %s", msg.Subject))
		event := new(event.BulkUserSimilarPhotoEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		s.logs.Log(fmt.Sprintf("unmarshalled event: %+v", event))
		err = s.userSimilarWorkerUC.CreateBulkUserSimilarPhotos(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	case "ai.single.facecam":
		event := new(event.UserSimiliarFacecamEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		err = s.userSimilarWorkerUC.CreateUserFacecam(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	case "ai.single.photo":
		event := new(event.UserSimilarEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("f%%ailed to unmarshal message : %s", err))
			return
		}

		err = s.userSimilarWorkerUC.CreateUserSimilar(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}
	}

	if err := msg.Ack(); err != nil {
		s.logs.Error(fmt.Sprintf("failed to ACK message : %s", err))
	}
}

func (s *AIConsumer) handleError(msg *nats.Msg, err error) {
	s.logs.Error(fmt.Sprintf("failed to process user similar with error : %v ", err))

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
