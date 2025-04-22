package model

import "time"

type RequestCreateWallet struct {
	CreatorId string
}

type WalletResponse struct {
	Id        string
	CreatorId string
	Balance   int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
