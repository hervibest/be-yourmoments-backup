package uploadconsumer

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

func (s *UploadConsumer) setupConsumer(subject string) error {
	consumerConfig := &nats.ConsumerConfig{
		Durable:       s.durableNames[subject],
		AckPolicy:     nats.AckExplicitPolicy,
		MaxDeliver:    5,
		BackOff:       []time.Duration{1 * time.Second, 5 * time.Second, 10 * time.Second},
		DeliverPolicy: nats.DeliverAllPolicy,
		AckWait:       30 * time.Second,
		FilterSubject: subject,
	}

	_, err := s.js.AddConsumer("UPLOAD_PHOTO_STREAM", consumerConfig)
	return err
}

func (s *UploadConsumer) handleMessage(ctx context.Context, msg *nats.Msg) {

	var err error
	switch msg.Subject {
	case "upload.bulk.photo":
		s.logs.Log(fmt.Sprintf("received message on subject: %s", msg.Subject))
		event := new(event.CreateBulkPhotoEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		s.logs.Log(fmt.Sprintf("unmarshalled event: %+v", event))
		err = s.photoWorkerUC.CreateBulkPhoto(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	case "upload.single.facecam":
		event := new(event.CreateFacecamEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("failed to unmarshal message : %s", err))
			return
		}

		err = s.facecameWorkerUC.CreateFacecam(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	case "upload.single.photo":
		event := new(event.CreatePhotoEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("f%%ailed to unmarshal message : %s", err))
			return
		}

		err = s.photoWorkerUC.CreatePhoto(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}

	case "upload.update.photo":
		event := new(event.UpdatePhotoDetailEvent)
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, event); err != nil {
			_ = msg.Nak()
			s.logs.Error(fmt.Sprintf("f%%ailed to unmarshal message : %s", err))
			return
		}

		err = s.photoWorkerUC.UpdatePhotoDetail(ctx, event)
		if err != nil {
			s.handleError(msg, err)
			return
		}
	}

	if err := msg.Ack(); err != nil {
		s.logs.Error(fmt.Sprintf("failed to ACK message : %s", err))
	}
}

func (s *UploadConsumer) handleError(msg *nats.Msg, err error) {
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
