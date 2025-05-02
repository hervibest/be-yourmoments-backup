package model

type AuthResponse struct {
	UserId      string
	Username    string
	Email       string
	PhoneNumber string
	Similarity  uint32
	CreatorId   string
	WalletId    string
}
