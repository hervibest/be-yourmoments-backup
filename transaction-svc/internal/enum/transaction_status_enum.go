package enum

type TransactionStatus string

var (
	TransactionStatusPendingTokenInit TransactionStatus = "PENDING_TOKEN_INIT"
	TransactionStatusPending          TransactionStatus = "PENDING"
	TransactionStatusSuccess          TransactionStatus = "SUCCESS"
	TransactionStatusFailed           TransactionStatus = "FAILED"
	TransactionStatusCancelled        TransactionStatus = "CANCELED"
	TransactionStatusExpired          TransactionStatus = "EXPIRED"
	TransactionStatusRefunded         TransactionStatus = "REFUNDED"
	TransactionStatusRefunding        TransactionStatus = "REFUNDING"
)
