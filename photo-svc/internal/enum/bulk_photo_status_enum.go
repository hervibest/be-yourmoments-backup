package enum

type BulkPhotoStatus string

var (
	BulkPhotoStatusProcessed BulkPhotoStatus = "PROCESSED"
	BulkPhotoStatusFailed    BulkPhotoStatus = "FAILED"
	BulkPhotoStatusCanceled  BulkPhotoStatus = "CANCELED"
	BulkPhotoStatusSuccess   BulkPhotoStatus = "SUCCESS"
)
