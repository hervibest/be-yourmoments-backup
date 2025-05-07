package enum

type PhotoStatusEnum string

const (
	PhotoStatusAvailableEnum     PhotoStatusEnum = "AVAILABLE"
	PhotoStatusInTransactionEnum PhotoStatusEnum = "IN_TRANSACTION"
	PhotoStatusSoldEnum          PhotoStatusEnum = "SOLD"
)
