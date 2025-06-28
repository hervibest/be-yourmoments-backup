package fcm

const (
	// Token errors (invalid or expired)
	ErrInvalidArgument                = "invalid-argument"
	ErrRegistrationTokenNotRegistered = "registration-token-not-registered"
	ErrUnregistered                   = "unregistered"
	ErrNotRegistered                  = "not-registered"
	ErrSenderIDMismatch               = "sender-id-mismatch"
	ErrMismatchedCredential           = "mismatched-credential"

	// Retryable errors
	ErrUnavailable         = "unavailable"
	ErrInternal            = "internal"
	ErrResourceExhausted   = "resource-exhausted"
	ErrMessageRateExceeded = "message-rate-exceeded"
	ErrQuotaExceeded       = "quota-exceeded"

	// Auth errors
	ErrPermissionDenied    = "permission-denied"
	ErrThirdPartyAuthError = "third-party-auth-error"
)
