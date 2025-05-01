package helper

import (
	"regexp"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/enum/fcm"
)

type ErrorFCM struct {
	StatusCode string
	Reason     string
	Code       string
	Details    string
	Token      string
	Raw        string
}

var fcmErrorRegex = regexp.MustCompile(`status: (\d+); reason: (.*?); code: (.*?); details: (.*)`)

func ParseFCMError(err error) *ErrorFCM {
	if err == nil {
		return nil
	}

	matches := fcmErrorRegex.FindStringSubmatch(err.Error())
	if len(matches) == 5 {
		return &ErrorFCM{
			StatusCode: matches[1],
			Reason:     matches[2],
			Code:       matches[3],
			Details:    matches[4],
			Raw:        err.Error(),
		}
	}

	return &ErrorFCM{
		Raw: err.Error(),
	}
}

func (e *ErrorFCM) IsInvalidToken() bool {
	switch e.Code {
	case
		fcm.ErrInvalidArgument,
		fcm.ErrRegistrationTokenNotRegistered,
		fcm.ErrUnregistered,
		fcm.ErrNotRegistered,
		fcm.ErrSenderIDMismatch,
		fcm.ErrMismatchedCredential:
		return true
	}
	return false
}

func (e *ErrorFCM) IsRetryable() bool {
	switch e.Code {
	case
		fcm.ErrUnavailable,
		fcm.ErrInternal,
		fcm.ErrResourceExhausted,
		fcm.ErrMessageRateExceeded,
		fcm.ErrQuotaExceeded:
		return true
	}
	return false
}

func (e *ErrorFCM) IsAuthError() bool {
	return e.Code == fcm.ErrPermissionDenied || e.Code == fcm.ErrThirdPartyAuthError
}

func (e *ErrorFCM) IsRateLimitError() bool {
	return e.Code == fcm.ErrMessageRateExceeded || e.Code == fcm.ErrQuotaExceeded
}
