package helper

import "be-yourmoments/user-svc/internal/enum/fcm"

// --- Helper cek error ---
func IsFCMInvalidTokenError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	case fcm.ErrInvalidRegistration, fcm.ErrNotRegistered:
		return true
	default:
		return false
	}
}

func IsFCMRetryableError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	case fcm.ErrUnavailable, fcm.ErrInternalServerError, fcm.ErrResourceExhausted:
		return true
	default:
		return false
	}
}
