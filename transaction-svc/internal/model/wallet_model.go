package model

import "time"

type CreateWalletRequest struct {
	CreatorId string
}

type GetWalletRequest struct {
	CreatorId string
}

type WalletResponse struct {
	Id        string
	CreatorId string
	Balance   int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
